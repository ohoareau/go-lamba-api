package common

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type KinesisRouter struct {
	Routes        map[string]KinesisRouteHandler
	BeforeRecords KinesisRouterBeforeRecordsFunc
	BeforeRecord  KinesisRouterBeforeRecordFunc
	AfterRecord   KinesisRouterAfterRecordFunc
	AfterRecords  KinesisRouterAfterRecordsFunc
}

func (r *KinesisRouter) AddRoute(selector string, handler KinesisRouteHandler) {
	r.Routes[selector] = handler
}

func (r *KinesisRouter) RegisterBeforeRecordsFunction(f KinesisRouterBeforeRecordsFunc) {
	r.BeforeRecords = f
}

func (r *KinesisRouter) RegisterBeforeRecordFunction(f KinesisRouterBeforeRecordFunc) {
	r.BeforeRecord = f
}

func (r *KinesisRouter) RegisterAfterRecordsFunction(f KinesisRouterAfterRecordsFunc) {
	r.AfterRecords = f
}

func (r *KinesisRouter) RegisterAfterRecordFunction(f KinesisRouterAfterRecordFunc) {
	r.AfterRecord = f
}

type KinesisRecordInfo struct {
	RecordIndex int
	Record      events.KinesisEventRecord
	Event       events.KinesisEvent
	Context     context.Context
}

type KinesisRouteHandler func(data []byte, info KinesisRecordInfo) (interface{}, error)
type Options struct {
	Apigw2Configurator   Apigw2Configurator
	Apigw1Configurator   Apigw1Configurator
	SnsConfigurator      SnsConfigurator
	SqsConfigurator      SqsConfigurator
	S3Configurator       S3Configurator
	KinesisConfigurator  KinesisConfigurator
	DynamodbConfigurator DynamodbConfigurator
	Features             Features
}

type KinesisRouterBeforeRecordsFunc func(event events.KinesisEvent, ctx context.Context) error
type KinesisRouterBeforeRecordFunc func(info KinesisRecordInfo) (KinesisRecordInfo, error)
type KinesisRouterAfterRecordFunc func(result KinesisRouterHandlerRecordResult, info KinesisRecordInfo) (KinesisRouterHandlerRecordResult, error)
type KinesisRouterAfterRecordsFunc func(results []KinesisRouterHandlerRecordResult, event events.KinesisEvent, ctx context.Context) (interface{}, error)

type KinesisRouterHandlerRecordResult struct {
	Result interface{}
	Error  error
}

func (r *KinesisRouter) ConvertKinesisEventRecordToData(info KinesisRecordInfo) ([]byte, error) {
	return info.Record.Kinesis.Data, nil
}

func (r *KinesisRouter) HandleRecord(info KinesisRecordInfo) (interface{}, error) {
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

func (r *KinesisRouter) SelectRouteHandler(info KinesisRecordInfo) (KinesisRouteHandler, error) {
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

func (r *KinesisRouter) Handle(event events.KinesisEvent, ctx context.Context) (interface{}, error) {
	var err error

	if nil != r.BeforeRecords {
		err = r.BeforeRecords(event, ctx)
		if nil != err {
			return nil, err
		}
	}
	var rrs []KinesisRouterHandlerRecordResult
	var info KinesisRecordInfo
	var result KinesisRouterHandlerRecordResult
	for i, record := range event.Records {
		info = KinesisRecordInfo{
			RecordIndex: i,
			Record:      record,
			Context:     ctx,
			Event:       event,
		}
		if nil != r.BeforeRecord {
			info, err = r.BeforeRecord(info)
			if nil != err {
				rrs = append(rrs, KinesisRouterHandlerRecordResult{
					Result: nil,
					Error:  err,
				})
				continue
			}
		}
		rr, err := r.HandleRecord(info)
		result = KinesisRouterHandlerRecordResult{
			Result: rr,
			Error:  err,
		}
		if nil != r.AfterRecord {
			result, err = r.AfterRecord(result, info)
			if nil != err {
				rrs = append(rrs, KinesisRouterHandlerRecordResult{
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
