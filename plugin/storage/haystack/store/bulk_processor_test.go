package store

import (
	"testing"

	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/objects"
	storageMetrics "github.com/jaegertracing/jaeger/storage/spanstore/metrics"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

func TestNewBulkProcessor(t *testing.T) {
	bp := testNewBulkProcessor(t)
	assert.NotEmpty(t, bp.(*bulkProcessor))
}

func TestBulkProcessor_IsCommitRequired(t *testing.T) {
	bp := testNewBulkProcessor(t)
	batch := []objects.HaystackSpan{}
	batchSizeSoFar := 0
	ok := bp.(*bulkProcessor).isCommitRequired(&batch, &batchSizeSoFar)
	assert.Equal(t, false, ok)
	for i := 0; i < 100; i++ {
		batch = append(batch, objects.HaystackSpan{})
	}

	ok = bp.(*bulkProcessor).isCommitRequired(&batch, &batchSizeSoFar)
	assert.Equal(t, false, ok)

	batchSizeSoFar = config.DEFAULT_BULK_SIZE + 100000
	ok = bp.(*bulkProcessor).isCommitRequired(&batch, &batchSizeSoFar)
	assert.Equal(t, true, ok)

	for i := 0; i < 1000; i++ {
		batch = append(batch, objects.HaystackSpan{})
	}
	batchSizeSoFar = 1000
	ok = bp.(*bulkProcessor).isCommitRequired(&batch, &batchSizeSoFar)
	assert.Equal(t, true, ok)
}

func TestBulkProcessor_Commit(t *testing.T) {

}

func testNewBulkProcessor(t *testing.T) BulkProcessor {
	v := viper.New()
	cfg := config.HaystackConfig{}
	cfg.InitFromViper(v)
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	request := make(chan objects.HaystackSpan)
	metricsFactory := metrics.NullFactory
	writeMetrics := storageMetrics.NewWriteMetrics(metricsFactory, "bulk_index")
	bp := NewBulkProcessor(logger, cfg, request, writeMetrics)
	return bp
}
