package objects

import (
	"encoding/json"
	"testing"

	"github.com/jaegertracing/jaeger/plugin/storage/es/spanstore/dbmodel"
	"github.com/stretchr/testify/assert"
)

func TestTransformToHaystackSpan(t *testing.T) {
	var (
		serviceName = "frontend"
	)
	span := &dbmodel.Span{
		TraceID:         "53c3135acaec68f5",
		SpanID:          "53c3135acaec68t3",
		Flags:           0,
		OperationName:   "cart",
		References:      nil,
		StartTime:       0,
		StartTimeMillis: 0,
		Duration:        0,
		Tags:            nil,
		Tag:             nil,
		Logs:            nil,
		Process: dbmodel.Process{
			ServiceName: serviceName,
		},
	}
	expectedMessage, err := json.Marshal(span)
	assert.Nil(t, err)
	haystackSpan, err := TransformToHaystackSpan(span, false)
	assert.Nil(t, err)
	actualMessage, err := json.Marshal(span)
	assert.Nil(t, err)
	assert.Equal(t, string(expectedMessage), string(actualMessage))
	assert.Equal(t, serviceName, haystackSpan.Meta.ServiceName)
	assert.Equal(t, TypeSpan, haystackSpan.Meta.Type)
	assert.GreaterOrEqual(t, len(expectedMessage), haystackSpan.Size())
}
