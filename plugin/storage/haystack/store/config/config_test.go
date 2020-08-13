package config

import (
	"flag"
	"go.uber.org/zap"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	v := viper.New()
	cfg := HaystackConfig{}

	cfg.InitFromViper(v)
	assert.NotEmpty(t, cfg)
	assert.Equal(t, cfg.EsAllTagsAsFields, DEFAULT_ES_ALL_TAGS_AS_FIELDS)
	assert.Equal(t, cfg.EsTagsAsFieldsConfigFile, DEFAULT_ES_TAGS_AS_FIELDS_CONFIG_FILE)
	assert.Equal(t, cfg.EsTagsAsFieldsDotReplacement, DEFAULT_ES_TAGS_AS_FIELDS_DOT_REPLACEMENT)
	assert.Equal(t, cfg.HttpMaxIdleConnsPerHost, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS_PER_HOST)
	assert.Equal(t, cfg.HttpMaxIdleConns, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS)
	assert.Equal(t, cfg.HttpRequestTimeout, DEFAULT_HTTP_REQUEST_TIMEOUT)
	assert.Equal(t, cfg.ProxyURL, "")
	assert.Equal(t, cfg.AuthToken, "")
	assert.Equal(t, cfg.BulkSize, DEFAULT_BULK_SIZE)
	assert.Equal(t, cfg.BulkActions, DEFAULT_BULK_ACTIONS)
	assert.Equal(t, cfg.SpanBatchFlushInterval, DEFAULT_BULK_FLUSH_INTERVAL)
	assert.Equal(t, cfg.WorkersCount, DEFAULT_BULK_WORKERS)
	assert.Equal(t, cfg.EnableJsonMsgFormat, DEFAULT_ENABLE_JSON_MSG_FORMAT)
	assert.Equal(t, cfg.SpanServiceCacheTTL, DEFAULT_SERVICE_CACHE_TTL)
	assert.Equal(t, cfg.SpanServiceCacheSize, DEFAULT_SERVICE_CACHE_SIZE)

	logger, e := zap.NewDevelopment()
	assert.Nil(t, e)
	err := cfg.Validate(logger)
	assert.NotNil(t, err)

	var eh flag.ErrorHandling
	flagset := flag.NewFlagSet("test", eh)
	cfg.AddFlags(flagset)
	flags := []string{
		authToken,
		proxyURL,
		httpRequestTimeout,
		httpMaxIdleConns,
		httpMaxIdleConnsPerHost,
		bulkActions,
		bulkSize,
		spanBatchFlushInterval,
		workersCount,
		spanServiceCacheSize,
		spanServiceCacheTTL,
		enableJsonMsgFormat,

		esAllTagsAsFields,
		esTagsAsFieldsConfigFile,
		esTagsAsFieldsDotReplacement,
	}
	for _, flagName := range flags {
		f := flagset.Lookup(flagName)
		assert.NotNil(t, f)
	}
}
