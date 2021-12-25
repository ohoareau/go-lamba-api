package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
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
	r := KinesisRouter{
		Routes: map[string]common.KinesisRouteHandler{},
	}
	if nil != options.KinesisConfigurator {
		options.KinesisConfigurator(r)
	}
	return r
}

type BeforeRecordsFunc func(event events.KinesisEvent, ctx context.Context) error
type BeforeRecordFunc func(info common.KinesisRecordInfo) (common.KinesisRecordInfo, error)
type AfterRecordFunc func(result HandlerRecordResult, info common.KinesisRecordInfo) (HandlerRecordResult, error)
type AfterRecordsFunc func(results []HandlerRecordResult, event events.KinesisEvent, ctx context.Context) (interface{}, error)

type KinesisRouter struct {
	Routes        map[string]common.KinesisRouteHandler
	BeforeRecords BeforeRecordsFunc
	BeforeRecord  BeforeRecordFunc
	AfterRecord   AfterRecordFunc
	AfterRecords  AfterRecordsFunc
}

func (r KinesisRouter) AddRoute(selector string, handler common.KinesisRouteHandler) {
	r.Routes[selector] = handler
}

func (r KinesisRouter) SelectRouteHandler(info common.KinesisRecordInfo) (common.KinesisRouteHandler, error) {
	streamName := info.Record.EventSourceArn[:strings.LastIndex(info.Record.EventSourceArn, "/")]
	v, found := r.Routes[streamName]
	if found {
		return v, nil
	}
	v, found = r.Routes["*"]
	if found {
		return v, nil
	}
	return nil, errors.New("no kinesis route handler found for stream '" + streamName + "'")
}
func (r KinesisRouter) ConvertKinesisEventRecordToData(info common.KinesisRecordInfo) ([]byte, error) {
	return info.Record.Kinesis.Data, nil
}

func (r KinesisRouter) HandleRecord(info common.KinesisRecordInfo) (interface{}, error) {
	handler, err := r.SelectRouteHandler(info)
	if nil != err {
		return nil, err
	}
	data, err := r.ConvertKinesisEventRecordToData(info)
	if nil != err {
		return nil, err
	}
	return handler(data, info)

}

type HandlerRecordResult struct {
	Result interface{}
	Error  error
}

func (r KinesisRouter) Handle(event events.KinesisEvent, ctx context.Context) (interface{}, error) {
	var err error

	if nil != r.BeforeRecords {
		err = r.BeforeRecords(event, ctx)
		if nil != err {
			return nil, err
		}
	}
	var rrs []HandlerRecordResult
	var info common.KinesisRecordInfo
	var result HandlerRecordResult
	for i, record := range event.Records {
		info = common.KinesisRecordInfo{
			RecordIndex: i,
			Record:      record,
			Context:     ctx,
			Event:       event,
		}
		if nil != r.BeforeRecord {
			info, err = r.BeforeRecord(info)
			if nil != err {
				rrs = append(rrs, HandlerRecordResult{
					Result: nil,
					Error:  err,
				})
				continue
			}
		}
		rr, err := r.HandleRecord(info)
		result = HandlerRecordResult{
			Result: rr,
			Error:  err,
		}
		if nil != r.AfterRecord {
			result, err = r.AfterRecord(result, info)
			if nil != err {
				rrs = append(rrs, HandlerRecordResult{
					Result: result,
					Error:  err,
				})
				continue
			}
		}
		rrs = append(rrs, result)
	}
	if nil != r.AfterRecords {
		return r.AfterRecords(rrs, event, ctx)
	}
	return rrs, err
}
