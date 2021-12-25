package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
	"strings"
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
func HandleKinesisEvent(event events.KinesisEvent, ctx context.Context, options common.Options) (interface{}, error) {
	return CreateKinesisRouter(options).Handle(event, ctx)
}

func CreateKinesisRouter(options common.Options) common.KinesisRouter {
	var r common.KinesisRouter
	if nil != options.KinesisConfigurator {
		options.KinesisConfigurator(r)
	}
	return r
}

type KinesisRouter map[string]common.KinesisRouteHandler

func (r KinesisRouter) AddRoute(selector string, handler common.KinesisRouteHandler) {
	r[selector] = handler
}

func (r KinesisRouter) SelectRouteHandler(info common.KinesisRecordInfo) (common.KinesisRouteHandler, error) {
	streamName := info.Record.EventSourceArn[:strings.LastIndex(info.Record.EventSourceArn, "/")]
	v, found := r[streamName]
	if found {
		return v, nil
	}
	v, found = r["*"]
	if found {
		return v, nil
	}
	return nil, errors.New("no kinesis route handler found for stream '" + streamName + "'")
}
func (r KinesisRouter) ConvertKinesisEventRecordToData(info common.KinesisRecordInfo) (interface{}, error) {
	return base64.StdEncoding.DecodeString(string(info.Record.Kinesis.Data))
}

func (r KinesisRouter) HandleRecord(info common.KinesisRecordInfo) error {
	handler, err := r.SelectRouteHandler(info)
	if nil != err {
		return err
	}
	data, err := r.ConvertKinesisEventRecordToData(info)
	if nil != err {
		return err
	}
	_, err = handler(data, info)
	return err

}
func (r KinesisRouter) Handle(event events.KinesisEvent, ctx context.Context) (interface{}, error) {
	var err error
	for i, record := range event.Records {
		err = r.HandleRecord(common.KinesisRecordInfo{
			RecordIndex: i,
			Record:      record,
			Context:     ctx,
			Event:       event,
		})
		if nil != err {
			log.Println(i, err)
		}
	}
	return nil, nil
}
