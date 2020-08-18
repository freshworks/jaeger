package store

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/objects"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

func TestNewHaystackSpanWriter(t *testing.T) {
	hs := testNewHaystackSpanWriter(t)
	defer hs.Close()
	assert.NotEmpty(t, hs)
}

func TestHaystackSpanWriter_WriteSpan(t *testing.T) {
	hs := testNewHaystackSpanWriter(t)
	defer hs.Close()
	span := &model.Span{
		TraceID:              model.TraceID{},
		SpanID:               0,
		OperationName:        "",
		References:           nil,
		Flags:                0,
		StartTime:            time.Time{},
		Duration:             0,
		Tags:                 nil,
		Logs:                 nil,
		Process:              nil,
		ProcessID:            "",
		Warnings:             nil,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	err := hs.WriteSpan(span)
	assert.Nil(t, err)
}

func TestHaystackSpanWriter_Close(t *testing.T) {
	hs := testNewHaystackSpanWriter(t)
	hs.Close()
	var pf assert.PanicTestFunc = func() {
		hs.writeCh <- objects.HaystackSpan{}
	}
	assert.Panics(t, pf)
}

func testNewHaystackSpanWriter(t *testing.T) *HaystackSpanWriter {
	v := viper.New()
	cfg := config.HaystackConfig{}
	cfg.InitFromViper(v)
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	var metricsFactory = metrics.NullFactory
	cfg.ProxyURL = "http://localhost:1111"
	cfg.AuthToken = "dummy"
	hs, err := NewHaystackSpanWriter(cfg, logger, metricsFactory)
	assert.Nil(t, err)
	return hs
}
