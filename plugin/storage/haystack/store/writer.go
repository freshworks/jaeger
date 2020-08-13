package store

import (
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/cache"
	esconfig "github.com/jaegertracing/jaeger/pkg/es/config"
	"github.com/jaegertracing/jaeger/plugin/storage/es/spanstore/dbmodel"
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/objects"
	storageMetrics "github.com/jaegertracing/jaeger/storage/spanstore/metrics"
	"github.com/uber/jaeger-lib/metrics"
)

// HaystackSpanWriter is a struct which holds haystack span writer properties
type HaystackSpanWriter struct {
	config          config.HaystackConfig
	logger          *zap.Logger
	serviceCache    cache.Cache
	writeCh         chan objects.HaystackSpan
	bulkProcessor   BulkProcessor
	writeMetrics    *storageMetrics.WriteMetrics
	esSpanConverter dbmodel.FromDomain
}

// NewHaystackSpanWriter creates a new haystack span writer for jaeger
func NewHaystackSpanWriter(config config.HaystackConfig, logger *zap.Logger, metricsFactory metrics.Factory) (*HaystackSpanWriter, error) {
	logger.Info("Haystack configuration", zap.String("config", config.String()))
	if err := config.Validate(logger); err != nil {
		return nil, err
	}
	spanWriter := &HaystackSpanWriter{
		config: config,
		logger: logger,
		serviceCache: cache.NewLRUWithOptions(
			config.SpanServiceCacheSize,
			&cache.Options{
				TTL: time.Duration(config.SpanServiceCacheTTL) * time.Second,
			},
		),
		writeCh:      make(chan objects.HaystackSpan),
		writeMetrics: storageMetrics.NewWriteMetrics(metricsFactory, "bulk_index"),
	}
	var tags []string
	if config.EsTagsAsFieldsConfigFile != "" {
		var err error
		if tags, err = esconfig.LoadTagsFromFile(config.EsTagsAsFieldsConfigFile); err != nil {
			logger.Error("Could not open file with tags", zap.Error(err))
			return nil, err
		}
	}
	spanWriter.esSpanConverter = dbmodel.NewFromDomain(config.EsAllTagsAsFields, tags, config.EsTagsAsFieldsDotReplacement)
	spanWriter.bulkProcessor = NewBulkProcessor(logger, config, spanWriter.writeCh, spanWriter.writeMetrics)
	// start bulkProcessor to read haystack spans and write to proxy service
	spanWriter.bulkProcessor.Start()

	return spanWriter, nil
}

// WriteSpan receives a Jaeger span, converts it to Haystack span and sends it to proxy service
func (sw *HaystackSpanWriter) WriteSpan(span *model.Span) error {
	defer func() {
		if r := recover(); r != nil {
			sw.logger.Warn("recovered from panic", zap.Any("error", r))
		}
	}()

	// Transform span into es span model
	esSpan := sw.esSpanConverter.FromDomainEmbedProcess(span)

	// Transform es span into haystack span model
	haystackSpan, err := objects.TransformToHaystackSpan(esSpan, sw.config.EnableJsonMsgFormat)
	if err != nil {
		sw.logger.Error("Failed to transform jaeger span to haystack span model", zap.String("error", err.Error()))
		return err
	}

	sw.writeCh <- haystackSpan

	service := objects.NewHaystackService(span)
	serviceHash, err := service.HashCode()
	if sw.serviceCache.Get(serviceHash) == nil || err != nil {
		if err == nil {
			sw.serviceCache.Put(serviceHash, serviceHash)
		}

		// Transform span into haystack span model
		haystackServiceSpan, err := objects.TransformToHaystackServiceSpan(service, sw.config.EnableJsonMsgFormat)
		if err != nil {
			sw.logger.Error("Failed to transform service span to haystack span model", zap.String("error", err.Error()))
			return err
		}
		sw.writeCh <- haystackServiceSpan
	}
	return nil
}

// Close stops and drains worker buffers
func (sw *HaystackSpanWriter) Close() error {
	sw.logger.Info("Stopping haystack span writer...")
	sw.logger.Info("Closing haystack storage request pipe")
	close(sw.writeCh)
	sw.bulkProcessor.Stop() // Stops Consuming spans from writeCh
	sw.logger.Info("Stopped haystack span writer...")
	return nil
}

// Not used
func (sw *HaystackSpanWriter) dropEmptyTags(tags []model.KeyValue) []model.KeyValue {
	for i, tag := range tags {
		if tag.Key == "" {
			tags[i] = tags[len(tags)-1] // Copy last element to index i.
			tags = tags[:len(tags)-1]   // Truncate slice.
			sw.logger.Warn(fmt.Sprintf("Found tag empty key: %s, dropping tag..", tag.String()))
		}
	}
	return tags
}
