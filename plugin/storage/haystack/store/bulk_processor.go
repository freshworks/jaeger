package store

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/client"
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/objects"
	storageMetrics "github.com/jaegertracing/jaeger/storage/spanstore/metrics"
	"go.uber.org/zap"
)

// bulkProcessor is responsible for starting concurrent worker threads
// worker threads are responsible for batching and sending data to storage endpoint
type bulkProcessor struct {
	logger             *zap.Logger
	httpClient         *client.HttpClient
	workerCount        int
	bulkSize           int
	bulkActions        int
	batchFlushInterval int
	request            <-chan objects.HaystackSpan
	writeMetrics       *storageMetrics.WriteMetrics
	waitGroup          sync.WaitGroup
}

type BulkProcessor interface {
	Start()
	Stop()
}

func NewBulkProcessor(logger *zap.Logger, cf config.HaystackConfig, request <-chan objects.HaystackSpan, writeMetrics *storageMetrics.WriteMetrics) BulkProcessor {
	return &bulkProcessor{
		logger:             logger,
		httpClient:         client.NewHttpClient(cf, logger),
		workerCount:        cf.WorkersCount,
		bulkActions:        cf.BulkActions,
		bulkSize:           cf.BulkSize,
		batchFlushInterval: cf.SpanBatchFlushInterval,
		request:            request,
		writeMetrics:       writeMetrics,
		waitGroup:          sync.WaitGroup{},
	}
}

func (b *bulkProcessor) Start() {
	b.logger.Info("Starting bulk processor workers", zap.Int("workersCount", b.workerCount))
	for i := 0; i < b.workerCount; i++ {
		b.waitGroup.Add(1)
		b.logger.Info("Started haystack storage writer bulk worker", zap.Int("workerId", i))
		go func(workerId int) {
			defer b.waitGroup.Done()
			var (
				batch          []objects.HaystackSpan
				ticker         = time.NewTicker(time.Second * time.Duration(b.batchFlushInterval))
				batchSizeSoFar int // bytes
			)
			for {
				select {
				case span, ok := <-b.request:
					if !ok { // on close of channel write last batch
						if len(batch) > 0 {
							b.commit(&batch, batchSizeSoFar, workerId)
							batch = []objects.HaystackSpan{} // reset buffer
							batchSizeSoFar = 0
						}
						b.logger.Info("Exiting worker", zap.Int("workerId", workerId))
						return
					}

					batch = append(batch, span)
					if ok := b.isCommitRequired(&batch, &batchSizeSoFar); ok {
						b.commit(&batch, batchSizeSoFar, workerId)
						batch = []objects.HaystackSpan{} // reset buffer
						batchSizeSoFar = 0
					}

				case <-ticker.C:
					if len(batch) > 0 {
						b.commit(&batch, batchSizeSoFar, workerId)
						batch = []objects.HaystackSpan{} // reset buffer
						batchSizeSoFar = 0
					}
				}
			}
		}(i)
	}
}

func (b *bulkProcessor) Stop() {
	b.waitGroup.Wait() // Wait un till all worker goroutines return.
}

func (b *bulkProcessor) isCommitRequired(batch *[]objects.HaystackSpan, batchSizeSoFar *int) bool {
	length := len(*batch)
	if length > 0 {
		span := &(*batch)[length-1]
		*batchSizeSoFar = *batchSizeSoFar + span.Size()
		if length >= b.bulkActions || *batchSizeSoFar >= b.bulkSize {
			return true
		}
	}
	return false
}

func (b *bulkProcessor) commit(batch *[]objects.HaystackSpan, batchSize int, workerId int) {
	var (
		err            error
		spanBatchBytes []byte
		start          = time.Now()
	)
	defer func() {
		diff := time.Since(start)
		b.logger.Debug("Time elapsed to perform commit operation", zap.String("duration", diff.String()), zap.Int("batchSizeInBytes", batchSize), zap.Int("batchLength", len(*batch)), zap.Int("workerId", workerId))
		b.writeMetrics.Emit(err, diff)
	}()

	spanBatch := objects.HaystackSpanBatchEvent{
		Events: *batch,
		Size:   len(*batch),
	}
	spanBatchBytes, err = json.Marshal(spanBatch)
	if err != nil {
		b.logger.Error("failed to marshal span batch", zap.Int("batchSize", len(*batch)), zap.String("error", err.Error()))
		return
	}
	err = b.httpClient.Post(spanBatchBytes)
	if err != nil {
		b.logger.Warn("failed to write batch to proxy service", zap.Int("batchSize", len(*batch)), zap.String("error", err.Error()))
		return
	}
	b.logger.Debug("successfully written batch of spans to proxy service ", zap.Int("batchSize", len(*batch)))
}
