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
	assert.Equal(t, cfg.EsAllTagsAsFields, DefaultEsAllTagsAsFields)
	assert.Equal(t, cfg.EsTagsAsFieldsConfigFile, DefaultEsTagsAsFieldsConfigFile)
	assert.Equal(t, cfg.EsTagsAsFieldsDotReplacement, DefaultEsTagsAsFieldsDotReplacement)
	assert.Equal(t, cfg.HTTPMaxIdleConnsPerHost, DefaultHTTPMaxIdleConnectionsPerHosts)
	assert.Equal(t, cfg.HTTPMaxIdleConns, DefaultHTTPMaxIdleConnections)
	assert.Equal(t, cfg.HTTPRequestTimeout, DefaultHTTPRequestTimeout)
	assert.Equal(t, cfg.ProxyURL, "")
	assert.Equal(t, cfg.AuthToken, "")
	assert.Equal(t, cfg.BulkSize, DefaultBulkSize)
	assert.Equal(t, cfg.BulkActions, DefaultBulkActions)
	assert.Equal(t, cfg.SpanBatchFlushInterval, DefaultBulkFlushInterval)
	assert.Equal(t, cfg.WorkersCount, DefaultBulkWorkers)
	assert.Equal(t, cfg.EnableJSONMsgFormat, DefaultEnableJSONMsgFormat)
	assert.Equal(t, cfg.SpanServiceCacheTTL, DefaultServiceCacheTTL)
	assert.Equal(t, cfg.SpanServiceCacheSize, DefaultServiceCacheSize)

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
		enableJSONMsgFormat,

		esAllTagsAsFields,
		esTagsAsFieldsConfigFile,
		esTagsAsFieldsDotReplacement,
	}
	for _, flagName := range flags {
		f := flagset.Lookup(flagName)
		assert.NotNil(t, f)
	}
}
