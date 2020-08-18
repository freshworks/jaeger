package objects

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/stretchr/testify/assert"
)

var (
	serviceName   = "frontend"
	operationName = "/cart"
	span          = &model.Span{
		TraceID:       model.NewTraceID(1, 1),
		SpanID:        model.NewSpanID(1),
		Flags:         0,
		OperationName: operationName,
		References:    nil,
		StartTime:     time.Now(),
		Duration:      0,
		Tags:          nil,
		Logs:          nil,
		Process: &model.Process{
			ServiceName: serviceName,
		},
		ProcessID:            "",
		Warnings:             nil,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
)

func TestTransformToHaystackServiceSpan(t *testing.T) {
	haystackService := NewHaystackService(span)
	assert.Equal(t, serviceName, haystackService.ServiceName)
	assert.Equal(t, operationName, haystackService.OperationName)

	haystackServiceSpan, err := TransformToHaystackServiceSpan(haystackService, false)
	assert.Nil(t, err)
	assert.Equal(t, TypeService, haystackServiceSpan.Meta.Type)
	assert.Equal(t, serviceName, haystackServiceSpan.Meta.ServiceName)
	expected, err := json.Marshal(haystackService)
	assert.Nil(t, err)
	assert.Equal(t, string(expected), haystackServiceSpan.Message)
	assert.GreaterOrEqual(t, len(expected), haystackServiceSpan.messageSize)
}

func TestHaystackService_HashCode(t *testing.T) {
	haystackService := NewHaystackService(span)
	hash, err := haystackService.HashCode()
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}
