package config

import (
	"encoding/json"
	"flag"
	"go.uber.org/zap"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	DefaultBulkSize                       = 2 << 20 // bytes
	DefaultBulkActions                    = 1000
	DefaultBulkFlushInterval              = 30 // seconds
	DefaultBulkWorkers                    = 1
	DefaultServiceCacheSize               = 100000
	DefaultServiceCacheTTL                = 86400 // seconds
	DefaultEnableJSONMsgFormat            = false
	DefaultHTTPRequestTimeout             = 1 // seconds
	DefaultHTTPMaxIdleConnections         = 100
	DefaultHTTPMaxIdleConnectionsPerHosts = 100
	DefaultEsAllTagsAsFields              = false
	DefaultEsTagsAsFieldsConfigFile       = ""
	DefaultEsTagsAsFieldsDotReplacement   = "@"
)

var (
	authToken               = "haystack.proxy.auth-token"
	proxyURL                = "haystack.proxy.url"
	httpRequestTimeout      = "haystack.proxy.client.request-timeout"
	httpMaxIdleConns        = "haystack.proxy.client.max-idle-conns"
	httpMaxIdleConnsPerHost = "haystack.proxy.client.max-idle-conns-per-host"
	bulkActions             = "haystack.bulk.actions"
	bulkSize                = "haystack.bulk.size"
	spanBatchFlushInterval  = "haystack.bulk.flush-interval"
	workersCount            = "haystack.bulk.workers"
	spanServiceCacheSize    = "haystack.service.cache.size"
	spanServiceCacheTTL     = "haystack.service.cache.ttl"
	enableJSONMsgFormat     = "haystack.enable-json-msg-format"

	esAllTagsAsFields            = "es.tags-as-fields.all"
	esTagsAsFieldsConfigFile     = "es.tags-as-fields.config-file"
	esTagsAsFieldsDotReplacement = "es.tags-as-fields.dot-replacement"
)

// HaystackConfig contains haystack storage config
type HaystackConfig struct {
	AuthToken               string `yaml:"haystack.proxy.auth-token"`
	ProxyURL                string `yaml:"haystack.proxy.url"`
	BulkSize                int    `yaml:"haystack.bulk.size"`
	BulkActions             int    `yaml:"haystack.bulk.actions"`
	SpanBatchFlushInterval  int    `yaml:"haystack.bulk.flush-interval"`
	WorkersCount            int    `yaml:"haystack.bulk.workers"`
	SpanServiceCacheSize    int    `yaml:"haystack.service.cache.size"`
	SpanServiceCacheTTL     int    `yaml:"haystack.service.cache.ttl"`
	EnableJSONMsgFormat     bool   `yaml:"haystack.enable-json-msg-format"`
	HTTPRequestTimeout      int    `yaml:"haystack.proxy.client.request-timeout"`
	HTTPMaxIdleConns        int    `yaml:"haystack.proxy.client.max-idle-conns"`
	HTTPMaxIdleConnsPerHost int    `yaml:"haystack.proxy.client.max-idle-conns-per-host"`

	EsAllTagsAsFields            bool   `yaml:"es.tags-as-fields.all"`
	EsTagsAsFieldsConfigFile     string `yaml:"es.tags-as-fields.config-file"`
	EsTagsAsFieldsDotReplacement string `yaml:"es.tags-as-fields.dot-replacement"`
}

// Validate configuration for mandatory fields
func (config *HaystackConfig) Validate(logger *zap.Logger) error {
	var (
		err error
	)
	if config.ProxyURL == "" {
		logger.Error("Proxy URL is not specified")
		err = errors.New("No proxy url provided in the config")
	}
	if config.AuthToken == "" {
		logger.Error("Auth token is not specified")
		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		err = errors.New("No auth token provided in the config. " + errStr)
	}
	return err
}

// used for logging
func (config *HaystackConfig) String() string {
	hc, _ := json.Marshal(config)
	return string(hc)
}

// InitFromViper initializes config from file or environment variables
func (config *HaystackConfig) InitFromViper(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetDefault(httpRequestTimeout, DefaultHTTPRequestTimeout)
	v.SetDefault(bulkActions, DefaultBulkActions)
	v.SetDefault(bulkSize, DefaultBulkSize) // bytes
	v.SetDefault(spanBatchFlushInterval, DefaultBulkFlushInterval)
	v.SetDefault(spanServiceCacheSize, DefaultServiceCacheSize)
	v.SetDefault(spanServiceCacheTTL, DefaultServiceCacheTTL)
	v.SetDefault(workersCount, DefaultBulkWorkers)
	v.SetDefault(enableJSONMsgFormat, DefaultEnableJSONMsgFormat)
	v.SetDefault(httpMaxIdleConns, DefaultHTTPMaxIdleConnections)
	v.SetDefault(httpMaxIdleConnsPerHost, DefaultHTTPMaxIdleConnectionsPerHosts)
	v.SetDefault(esAllTagsAsFields, DefaultEsAllTagsAsFields)
	v.SetDefault(esTagsAsFieldsConfigFile, DefaultEsTagsAsFieldsConfigFile)
	v.SetDefault(esTagsAsFieldsDotReplacement, DefaultEsTagsAsFieldsDotReplacement)
	v.AutomaticEnv()

	*config = HaystackConfig{
		AuthToken:                    v.GetString(authToken),
		ProxyURL:                     v.GetString(proxyURL),
		BulkActions:                  v.GetInt(bulkActions),
		BulkSize:                     v.GetInt(bulkSize),
		SpanBatchFlushInterval:       v.GetInt(spanBatchFlushInterval),
		SpanServiceCacheSize:         v.GetInt(spanServiceCacheSize),
		SpanServiceCacheTTL:          v.GetInt(spanServiceCacheTTL),
		WorkersCount:                 v.GetInt(workersCount),
		EnableJSONMsgFormat:          v.GetBool(enableJSONMsgFormat),
		HTTPMaxIdleConns:             v.GetInt(httpMaxIdleConns),
		HTTPRequestTimeout:           v.GetInt(httpRequestTimeout),
		HTTPMaxIdleConnsPerHost:      v.GetInt(httpMaxIdleConnsPerHost),
		EsAllTagsAsFields:            v.GetBool(esAllTagsAsFields),
		EsTagsAsFieldsConfigFile:     v.GetString(esTagsAsFieldsConfigFile),
		EsTagsAsFieldsDotReplacement: v.GetString(esTagsAsFieldsDotReplacement),
	}
}

// AddFlags adds flags for config
func (config *HaystackConfig) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(authToken, "", "Auth token required by haystack proxy api")
	flagSet.String(proxyURL, "", "Proxy URL required by haystack")
	flagSet.Int(httpRequestTimeout, DefaultHTTPRequestTimeout, "Request timeout specified in seconds")
	flagSet.Int(bulkActions, DefaultBulkActions, "Number of spans to be sent in each request")
	flagSet.Int(bulkSize, DefaultBulkSize, "Maximum size of request payload")
	flagSet.Int(spanBatchFlushInterval, DefaultBulkFlushInterval, "Periodic time interval to flush the batch of spans. Specified in seconds")
	flagSet.Int(spanServiceCacheSize, DefaultServiceCacheSize, "Service cache size")
	flagSet.Int(spanServiceCacheTTL, DefaultServiceCacheTTL, "Service cache ttl")
	flagSet.Int(workersCount, DefaultBulkWorkers, "Number of workers concurrently read from write channel")
	flagSet.Bool(enableJSONMsgFormat, DefaultEnableJSONMsgFormat, "enableJsonMsgFormat false send message as json string else as a object")
	flagSet.Int(httpMaxIdleConns, DefaultHTTPMaxIdleConnections, "max idle http client connections")
	flagSet.Int(httpMaxIdleConnsPerHost, DefaultHTTPMaxIdleConnectionsPerHosts, "max idle http client connections per host")
	flagSet.Bool(esAllTagsAsFields, DefaultEsAllTagsAsFields, "es all tags as fields")
	flagSet.String(esTagsAsFieldsConfigFile, DefaultEsTagsAsFieldsConfigFile, "es tags as fields config file")
	flagSet.String(esTagsAsFieldsDotReplacement, DefaultEsTagsAsFieldsDotReplacement, "es fields dot replacement")
}
