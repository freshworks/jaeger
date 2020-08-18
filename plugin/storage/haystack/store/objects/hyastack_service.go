package objects

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/jaegertracing/jaeger/model"
)

//HaystackService type, for query purposes
type HaystackService struct {
	OperationName string `json:"operationName"`
	ServiceName   string `json:"serviceName"`
}

//NewHaystackService creates a new haystack service from a span
func NewHaystackService(span *model.Span) *HaystackService {
	service := &HaystackService{
		ServiceName:   span.Process.ServiceName,
		OperationName: span.OperationName,
	}
	return service
}

// HashCode receives a haystack service and returns a hash representation of it's service name and operation name.
func (service *HaystackService) HashCode() (string, error) {
	hash := fnv.New64a()
	_, err := hash.Write(append([]byte(service.ServiceName), []byte(service.OperationName)...))
	return fmt.Sprintf("%x", hash.Sum64()), err
}

// TransformToHaystackServiceSpan converts service span into HaystackSpan model
func TransformToHaystackServiceSpan(service *HaystackService, jsonMsgFormat bool) (HaystackSpan, error) {
	var (
		haystackServiceSpan = HaystackSpan{}
		message             interface{}
		messageSize         int
	)

	if jsonMsgFormat {
		message = service
	} else {
		serviceData, err := json.Marshal(service)
		if err != nil {
			return haystackServiceSpan, err
		}
		message = string(serviceData)
		messageSize = len(string(serviceData))
	}
	haystackServiceSpan = HaystackSpan{
		Meta: MetaData{
			Type:        TypeService,
			ServiceName: service.ServiceName,
		},
		Message:     message,
		messageSize: messageSize,
	}
	return haystackServiceSpan, nil
}
