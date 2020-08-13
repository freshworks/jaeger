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
	assert.Equal(t, f.config.EsAllTagsAsFields, config.DEFAULT_ES_ALL_TAGS_AS_FIELDS)
	assert.Equal(t, f.config.EsTagsAsFieldsConfigFile, config.DEFAULT_ES_TAGS_AS_FIELDS_CONFIG_FILE)
	assert.Equal(t, f.config.EsTagsAsFieldsDotReplacement, config.DEFAULT_ES_TAGS_AS_FIELDS_DOT_REPLACEMENT)
	assert.Equal(t, f.config.HttpMaxIdleConnsPerHost, config.DEFAULT_HTTP_MAX_IDLE_CONNECTIONS_PER_HOST)
	assert.Equal(t, f.config.HttpMaxIdleConns, config.DEFAULT_HTTP_MAX_IDLE_CONNECTIONS)
	assert.Equal(t, f.config.HttpRequestTimeout, config.DEFAULT_HTTP_REQUEST_TIMEOUT)
	assert.Equal(t, f.config.ProxyURL, "")
	assert.Equal(t, f.config.AuthToken, "")
	assert.Equal(t, f.config.BulkSize, config.DEFAULT_BULK_SIZE)
	assert.Equal(t, f.config.BulkActions, config.DEFAULT_BULK_ACTIONS)
	assert.Equal(t, f.config.SpanBatchFlushInterval, config.DEFAULT_BULK_FLUSH_INTERVAL)
	assert.Equal(t, f.config.WorkersCount, config.DEFAULT_BULK_WORKERS)
	assert.Equal(t, f.config.EnableJsonMsgFormat, config.DEFAULT_ENABLE_JSON_MSG_FORMAT)
	assert.Equal(t, f.config.SpanServiceCacheTTL, config.DEFAULT_SERVICE_CACHE_TTL)
	assert.Equal(t, f.config.SpanServiceCacheSize, config.DEFAULT_SERVICE_CACHE_SIZE)
}
