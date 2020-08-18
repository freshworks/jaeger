package store

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	batchSizeSoFar = config.DefaultBulkSize + 100000
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
	var (
		batch = []objects.HaystackSpan{objects.HaystackSpan{
			Meta: objects.MetaData{
				ServiceName: "test",
				Type:        "span",
			},
			Message: "test",
		}}
		batchSize = 100
	)
	bp := testNewBulkProcessor(t)

	var jobHandler = func(w http.ResponseWriter, r *http.Request) {
		var p = objects.HaystackSpanBatchEvent{}
		err := json.NewDecoder(r.Body).Decode(&p)
		assert.Nil(t, err)
		assert.Equal(t, 1, p.Size)
		assert.Equal(t, 1, len(p.Events))
	}
	server := startMockESCluster(jobHandler)
	bp.(*bulkProcessor).httpClient.SetEndpoint(server.URL)
	bp.(*bulkProcessor).commit(&batch, batchSize, 0)
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

func startMockESCluster(jobHandler func(_ http.ResponseWriter, r *http.Request)) *httptest.Server {
	server := httptest.NewUnstartedServer(http.HandlerFunc(jobHandler))
	server.Start()
	return server
}
