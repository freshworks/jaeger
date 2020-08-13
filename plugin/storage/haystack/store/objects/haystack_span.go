package objects

import (
	"encoding/json"

	"github.com/jaegertracing/jaeger/plugin/storage/es/spanstore/dbmodel"
)

// SPAN TYPE
const (
	TYPE_SPAN    = "span"
	TYPE_SERVICE = "service"
)

type HaystackSpan struct {
	Meta        MetaData    `json:"meta"`
	Message     interface{} `json:"msg"`
	messageSize int
}

type MetaData struct {
	Type        string `json:"type"`
	ServiceName string `json:"serviceName"`
}

type HaystackSpanBatchEvent struct {
	Events []HaystackSpan `json:"events"`
	Size   int            `json:"size"`
}

func (hs *HaystackSpan) String() string {
	haystackSpanBytes, _ := json.Marshal(hs)
	return string(haystackSpanBytes)
}

func (hs *HaystackSpan) Size() int {
	return hs.messageSize
}

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
			Type:        TYPE_SPAN,
			ServiceName: span.Process.ServiceName,
		},
		Message:     message,
		messageSize: messageSize,
	}
	return haystackSpan, nil
}
