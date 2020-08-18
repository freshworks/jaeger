package store

import (
	"context"
	"go.uber.org/zap"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

// HaystackSpanReader is a struct which holds haystack span reader properties
type HaystackSpanReader struct {
}

// NewHaystackSpanReader creates a new haystack span reader
func NewHaystackSpanReader(config config.HaystackConfig, logger *zap.Logger) *HaystackSpanReader {
	reader := &HaystackSpanReader{}
	return reader
}

// GetTrace returns a Jaeger trace by traceID
func (reader *HaystackSpanReader) GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error) {
	return nil, nil
}

// GetServices returns an array of all the service names that are being monitored
func (reader *HaystackSpanReader) GetServices(ctx context.Context) ([]string, error) {
	return nil, nil
}

// GetOperations returns an array of all the operations a specific service performed
func (reader *HaystackSpanReader) GetOperations(ctx context.Context, service spanstore.OperationQueryParameters) ([]spanstore.Operation, error) {
	return nil, nil
}

// FindTraces return an array of Jaeger traces by a search query
func (reader *HaystackSpanReader) FindTraces(ctx context.Context, query *spanstore.TraceQueryParameters) ([]*model.Trace, error) {
	return nil, nil
}

// FindTraceIDs retrieve traceIDs that match the traceQuery
func (reader *HaystackSpanReader) FindTraceIDs(ctx context.Context, query *spanstore.TraceQueryParameters) ([]model.TraceID, error) {
	return nil, nil
}

// GetDependencies returns an array of all the dependencies in a specific time range
func (*HaystackSpanReader) GetDependencies(endTs time.Time, lookback time.Duration) ([]model.DependencyLink, error) {
	return nil, nil
}
