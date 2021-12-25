package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
)

func ConvertPayloadToKinesisEvent(payload []byte) events.KinesisEvent {
	var event events.KinesisEvent
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err.Error())
	}
	return event
}

//goland:noinspection GoUnusedParameter
func HandleKinesisEvent(event events.KinesisEvent, ctx context.Context, options *common.Options) (interface{}, error) {
	return CreateKinesisRouter(options).Handle(event, ctx)
}

func CreateKinesisRouter(options *common.Options) *common.KinesisRouter {
	r := &common.KinesisRouter{
		Routes: map[string]common.KinesisRouteHandler{},
	}
	if nil != options.KinesisConfigurator {
		options.KinesisConfigurator(r)
	}
	return r
}
