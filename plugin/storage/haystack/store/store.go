package store

import (
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type Store struct {
	reader *HaystackSpanReader
	writer *HaystackSpanWriter
}

func NewHaystackStore(config config.HaystackConfig, logger *zap.Logger, metricsFactory metrics.Factory) *Store {
	logger.Info("Initialized haystack storage...")
	reader := NewHaystackSpanReader(config, logger)
	writer, err := NewHaystackSpanWriter(config, logger, metricsFactory)
	if err != nil {
		logger.Fatal("Failed to create haystack span writer ", zap.String("error", err.Error()))
	}
	return &Store{
		reader: reader,
		writer: writer,
	}
}

func (store *Store) Close() {
	store.writer.Close()
}

func (store *Store) SpanReader() spanstore.Reader {
	return store.reader
}

func (store *Store) SpanWriter() spanstore.Writer {
	return store.writer
}

func (store *Store) DependencyReader() dependencystore.Reader {
	return store.reader
}
