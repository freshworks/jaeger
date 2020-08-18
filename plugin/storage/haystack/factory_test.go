package haystack

import (
	"flag"
	//"fmt"
	"testing"

	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	f := NewFactory()
	assert.Empty(t, f)

	flagSet := &flag.FlagSet{}
	f.AddFlags(flagSet)
	assert.NotEmpty(t, flagSet)

	v := viper.New()
	f.InitFromViper(v)
	assert.NotEmpty(t, f.config)
	assert.Equal(t, f.config.EsAllTagsAsFields, config.DefaultEsAllTagsAsFields)
	assert.Equal(t, f.config.EsTagsAsFieldsConfigFile, config.DefaultEsTagsAsFieldsConfigFile)
	assert.Equal(t, f.config.EsTagsAsFieldsDotReplacement, config.DefaultEsTagsAsFieldsDotReplacement)
	assert.Equal(t, f.config.HTTPMaxIdleConnsPerHost, config.DefaultHTTPMaxIdleConnectionsPerHosts)
	assert.Equal(t, f.config.HTTPMaxIdleConns, config.DefaultHTTPMaxIdleConnections)
	assert.Equal(t, f.config.HTTPRequestTimeout, config.DefaultHTTPRequestTimeout)
	assert.Equal(t, f.config.ProxyURL, "")
	assert.Equal(t, f.config.AuthToken, "")
	assert.Equal(t, f.config.BulkSize, config.DefaultBulkSize)
	assert.Equal(t, f.config.BulkActions, config.DefaultBulkActions)
	assert.Equal(t, f.config.SpanBatchFlushInterval, config.DefaultBulkFlushInterval)
	assert.Equal(t, f.config.WorkersCount, config.DefaultBulkWorkers)
	assert.Equal(t, f.config.EnableJSONMsgFormat, config.DefaultEnableJSONMsgFormat)
	assert.Equal(t, f.config.SpanServiceCacheTTL, config.DefaultServiceCacheTTL)
	assert.Equal(t, f.config.SpanServiceCacheSize, config.DefaultServiceCacheSize)
}
