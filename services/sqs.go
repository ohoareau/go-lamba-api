package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
)

func ConvertPayloadToSqsEvent(payload []byte) events.SQSEvent {
	var event events.SQSEvent
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err.Error())
	}
	return event
}

//goland:noinspection GoUnusedParameter
func HandleSqsEvent(event events.SQSEvent, ctx context.Context, options *common.Options) (interface{}, error) {
	return CreateSqsRouter(options).Handle(event, ctx)
}

func CreateSqsRouter(options *common.Options) *common.SqsRouter {
	r := &common.SqsRouter{
		Routes: map[string]common.SqsRouteHandler{},
	}
	if nil != options.SqsConfigurator {
		options.SqsConfigurator(r)
	}
	return r
}
