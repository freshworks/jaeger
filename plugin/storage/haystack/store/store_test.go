package store

import (
	"github.com/spf13/viper"
	"testing"

	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

func TestNewHaystackStore(t *testing.T) {
	cfg := config.HaystackConfig{}
	v := viper.New()
	cfg.InitFromViper(v)
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	metricsFactory := metrics.NullFactory
	store := NewHaystackStore(cfg, logger, metricsFactory)
	assert.NotEmpty(t, store)
}
