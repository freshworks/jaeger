package config

import (
	"encoding/json"
	"flag"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_BULK_SIZE                          = 2 << 20 // bytes
	DEFAULT_BULK_ACTIONS                       = 1000
	DEFAULT_BULK_FLUSH_INTERVAL                = 30 // seconds
	DEFAULT_BULK_WORKERS                       = 1
	DEFAULT_SERVICE_CACHE_SIZE                 = 100000
	DEFAULT_SERVICE_CACHE_TTL                  = 86400 // seconds
	DEFAULT_ENABLE_JSON_MSG_FORMAT             = false
	DEFAULT_HTTP_REQUEST_TIMEOUT               = 1 // seconds
	DEFAULT_HTTP_MAX_IDLE_CONNECTIONS          = 100
	DEFAULT_HTTP_MAX_IDLE_CONNECTIONS_PER_HOST = 100

	DEFAULT_ES_ALL_TAGS_AS_FIELDS             = false
	DEFAULT_ES_TAGS_AS_FIELDS_CONFIG_FILE     = ""
	DEFAULT_ES_TAGS_AS_FIELDS_DOT_REPLACEMENT = "@"
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
	enableJsonMsgFormat     = "haystack.enable-json-msg-format"

	esAllTagsAsFields            = "es.tags-as-fields.all"
	esTagsAsFieldsConfigFile     = "es.tags-as-fields.config-file"
	esTagsAsFieldsDotReplacement = "es.tags-as-fields.dot-replacement"
)

type HaystackConfig struct {
	AuthToken               string `yaml:"haystack.proxy.auth-token"`
	ProxyURL                string `yaml:"haystack.proxy.url"`
	BulkSize                int    `yaml:"haystack.bulk.size"`
	BulkActions             int    `yaml:"haystack.bulk.actions"`
	SpanBatchFlushInterval  int    `yaml:"haystack.bulk.flush-interval"`
	WorkersCount            int    `yaml:"haystack.bulk.workers"`
	SpanServiceCacheSize    int    `yaml:"haystack.service.cache.size"`
	SpanServiceCacheTTL     int    `yaml:"haystack.service.cache.ttl"`
	EnableJsonMsgFormat     bool   `yaml:"haystack.enable-json-msg-format"`
	HttpRequestTimeout      int    `yaml:"haystack.proxy.client.request-timeout"`
	HttpMaxIdleConns        int    `yaml:"haystack.proxy.client.max-idle-conns"`
	HttpMaxIdleConnsPerHost int    `yaml:"haystack.proxy.client.max-idle-conns-per-host"`

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

//ParseConfig receives a config file path, parse it and returns haystack span store config
func ParseConfig(filePath string, logger *zap.Logger) (*HaystackConfig, error) {
	var haystackConfig *HaystackConfig
	if filePath != "" { // read config from file
		haystackConfig = &HaystackConfig{}
		yamlFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(yamlFile, haystackConfig)
	} else { // read config from environment variables
		v := viper.New()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.SetDefault(httpRequestTimeout, DEFAULT_HTTP_REQUEST_TIMEOUT)
		v.SetDefault(bulkActions, DEFAULT_BULK_ACTIONS)
		v.SetDefault(bulkSize, DEFAULT_BULK_SIZE) // bytes
		v.SetDefault(spanBatchFlushInterval, DEFAULT_BULK_FLUSH_INTERVAL)
		v.SetDefault(spanServiceCacheSize, DEFAULT_SERVICE_CACHE_SIZE)
		v.SetDefault(spanServiceCacheTTL, DEFAULT_SERVICE_CACHE_TTL)
		v.SetDefault(workersCount, DEFAULT_BULK_WORKERS)
		v.SetDefault(enableJsonMsgFormat, DEFAULT_ENABLE_JSON_MSG_FORMAT)
		v.SetDefault(httpMaxIdleConns, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS)
		v.SetDefault(httpMaxIdleConnsPerHost, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS_PER_HOST)
		v.SetDefault(esAllTagsAsFields, DEFAULT_ES_ALL_TAGS_AS_FIELDS)
		v.SetDefault(esTagsAsFieldsConfigFile, DEFAULT_ES_TAGS_AS_FIELDS_CONFIG_FILE)
		v.SetDefault(esTagsAsFieldsDotReplacement, DEFAULT_ES_TAGS_AS_FIELDS_DOT_REPLACEMENT)
		v.AutomaticEnv()

		haystackConfig = &HaystackConfig{
			AuthToken:                    v.GetString(authToken),
			ProxyURL:                     v.GetString(proxyURL),
			HttpRequestTimeout:           v.GetInt(httpRequestTimeout),
			BulkActions:                  v.GetInt(bulkActions),
			BulkSize:                     v.GetInt(bulkSize),
			SpanBatchFlushInterval:       v.GetInt(spanBatchFlushInterval),
			SpanServiceCacheSize:         v.GetInt(spanServiceCacheSize),
			SpanServiceCacheTTL:          v.GetInt(spanServiceCacheTTL),
			WorkersCount:                 v.GetInt(workersCount),
			EnableJsonMsgFormat:          v.GetBool(enableJsonMsgFormat),
			HttpMaxIdleConns:             v.GetInt(httpMaxIdleConns),
			HttpMaxIdleConnsPerHost:      v.GetInt(httpMaxIdleConnsPerHost),
			EsAllTagsAsFields:            v.GetBool(esAllTagsAsFields),
			EsTagsAsFieldsConfigFile:     v.GetString(esTagsAsFieldsConfigFile),
			EsTagsAsFieldsDotReplacement: v.GetString(esTagsAsFieldsDotReplacement),
		}
	}

	if err := haystackConfig.Validate(logger); err != nil {
		return nil, err
	}
	return haystackConfig, nil
}

// used for logging
func (config *HaystackConfig) String() string {
	hc, _ := json.Marshal(config)
	return string(hc)
}

func (config *HaystackConfig) InitFromViper(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetDefault(httpRequestTimeout, DEFAULT_HTTP_REQUEST_TIMEOUT)
	v.SetDefault(bulkActions, DEFAULT_BULK_ACTIONS)
	v.SetDefault(bulkSize, DEFAULT_BULK_SIZE) // bytes
	v.SetDefault(spanBatchFlushInterval, DEFAULT_BULK_FLUSH_INTERVAL)
	v.SetDefault(spanServiceCacheSize, DEFAULT_SERVICE_CACHE_SIZE)
	v.SetDefault(spanServiceCacheTTL, DEFAULT_SERVICE_CACHE_TTL)
	v.SetDefault(workersCount, DEFAULT_BULK_WORKERS)
	v.SetDefault(enableJsonMsgFormat, DEFAULT_ENABLE_JSON_MSG_FORMAT)
	v.SetDefault(httpMaxIdleConns, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS)
	v.SetDefault(httpMaxIdleConnsPerHost, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS_PER_HOST)
	v.SetDefault(esAllTagsAsFields, DEFAULT_ES_ALL_TAGS_AS_FIELDS)
	v.SetDefault(esTagsAsFieldsConfigFile, DEFAULT_ES_TAGS_AS_FIELDS_CONFIG_FILE)
	v.SetDefault(esTagsAsFieldsDotReplacement, DEFAULT_ES_TAGS_AS_FIELDS_DOT_REPLACEMENT)
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
		EnableJsonMsgFormat:          v.GetBool(enableJsonMsgFormat),
		HttpMaxIdleConns:             v.GetInt(httpMaxIdleConns),
		HttpRequestTimeout:           v.GetInt(httpRequestTimeout),
		HttpMaxIdleConnsPerHost:      v.GetInt(httpMaxIdleConnsPerHost),
		EsAllTagsAsFields:            v.GetBool(esAllTagsAsFields),
		EsTagsAsFieldsConfigFile:     v.GetString(esTagsAsFieldsConfigFile),
		EsTagsAsFieldsDotReplacement: v.GetString(esTagsAsFieldsDotReplacement),
	}
}

func (config *HaystackConfig) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(authToken, "", "Auth token required by haystack proxy api")
	flagSet.String(proxyURL, "", "Proxy URL required by haystack")
	flagSet.Int(httpRequestTimeout, DEFAULT_HTTP_REQUEST_TIMEOUT, "Request timeout specified in seconds")
	flagSet.Int(bulkActions, DEFAULT_BULK_ACTIONS, "Number of spans to be sent in each request")
	flagSet.Int(bulkSize, DEFAULT_BULK_SIZE, "Maximum size of request payload")
	flagSet.Int(spanBatchFlushInterval, DEFAULT_BULK_FLUSH_INTERVAL, "Periodic time interval to flush the batch of spans. Specified in seconds")
	flagSet.Int(spanServiceCacheSize, DEFAULT_SERVICE_CACHE_SIZE, "Service cache size")
	flagSet.Int(spanServiceCacheTTL, DEFAULT_SERVICE_CACHE_TTL, "Service cache ttl")
	flagSet.Int(workersCount, DEFAULT_BULK_WORKERS, "Number of workers concurrently read from write channel")
	flagSet.Bool(enableJsonMsgFormat, DEFAULT_ENABLE_JSON_MSG_FORMAT, "enableJsonMsgFormat false send message as json string else as a object")
	flagSet.Int(httpMaxIdleConns, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS, "max idle http client connections")
	flagSet.Int(httpMaxIdleConnsPerHost, DEFAULT_HTTP_MAX_IDLE_CONNECTIONS_PER_HOST, "max idle http client connections per host")
	flagSet.Bool(esAllTagsAsFields, DEFAULT_ES_ALL_TAGS_AS_FIELDS, "es all tags as fields")
	flagSet.String(esTagsAsFieldsConfigFile, DEFAULT_ES_TAGS_AS_FIELDS_CONFIG_FILE, "es tags as fields config file")
	flagSet.String(esTagsAsFieldsDotReplacement, DEFAULT_ES_TAGS_AS_FIELDS_DOT_REPLACEMENT, "es fields dot replacement")
}
