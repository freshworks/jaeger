package objects

import (
	"encoding/json"

	"github.com/jaegertracing/jaeger/plugin/storage/es/spanstore/dbmodel"
)

// SPAN TYPE
const (
	TypeSpan    = "span"
	TypeService = "service"
)

// HaystackSpan struct specifies the haystack storage span model
type HaystackSpan struct {
	Meta        MetaData    `json:"meta"`
	Message     interface{} `json:"msg"`
	messageSize int
}

// MetaData struct contains meta info
type MetaData struct {
	Type        string `json:"type"`
	ServiceName string `json:"serviceName"`
}

// HaystackSpanBatchEvent struct specifies batch request
type HaystackSpanBatchEvent struct {
	Events []HaystackSpan `json:"events"`
	Size   int            `json:"size"`
}

// Used for testing and logging
func (hs *HaystackSpan) String() string {
	haystackSpanBytes, _ := json.Marshal(hs)
	return string(haystackSpanBytes)
}

// Size returns the span message size in bytes
func (hs *HaystackSpan) Size() int {
	return hs.messageSize
}

// TransformToHaystackSpan converts dbmode.Span into HaystackSpan model
func TransformToHaystackSpan(span *dbmodel.Span, jsonMsgFormat bool) (HaystackSpan, error) {
	var (
		haystackSpan = HaystackSpan{}
		message      interface{}
		messageSize  int
	)
	if jsonMsgFormat {
		message = span
	} else {
		spanData, err := json.Marshal(span)
		if err != nil {
			return haystackSpan, err
		}
		message = string(spanData)
		messageSize = len(string(spanData))
	}
	haystackSpan = HaystackSpan{
		Meta: MetaData{
			Type:        TypeSpan,
			ServiceName: span.Process.ServiceName,
		},
		Message:     message,
		messageSize: messageSize,
	}
	return haystackSpan, nil
}
