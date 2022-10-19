package common

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type SqsRouter struct {
	Routes        map[string]SqsRouteHandler
	BeforeRecords SqsRouterBeforeRecordsFunc
	BeforeRecord  SqsRouterBeforeRecordFunc
	AfterRecord   SqsRouterAfterRecordFunc
	AfterRecords  SqsRouterAfterRecordsFunc
}

func (r *SqsRouter) AddRoute(selector string, handler SqsRouteHandler) {
	r.Routes[selector] = handler
}

func (r *SqsRouter) RegisterBeforeRecordsFunction(f SqsRouterBeforeRecordsFunc) {
	r.BeforeRecords = f
}

func (r *SqsRouter) RegisterBeforeRecordFunction(f SqsRouterBeforeRecordFunc) {
	r.BeforeRecord = f
}

func (r *SqsRouter) RegisterAfterRecordsFunction(f SqsRouterAfterRecordsFunc) {
	r.AfterRecords = f
}

func (r *SqsRouter) RegisterAfterRecordFunction(f SqsRouterAfterRecordFunc) {
	r.AfterRecord = f
}

type SqsRecordInfo struct {
	RecordIndex int
	Record      events.SQSMessage
	Event       events.SQSEvent
	Context     context.Context
}

type SqsRouteHandler func(data []byte, info SqsRecordInfo, queueName string) (interface{}, error)

type SqsRouterBeforeRecordsFunc func(event events.SQSEvent, ctx context.Context) error
type SqsRouterBeforeRecordFunc func(info SqsRecordInfo) (SqsRecordInfo, error)
type SqsRouterAfterRecordFunc func(result SqsRouterHandlerRecordResult, info SqsRecordInfo) (SqsRouterHandlerRecordResult, error)
type SqsRouterAfterRecordsFunc func(results []SqsRouterHandlerRecordResult, event events.SQSEvent, ctx context.Context) (interface{}, error)

type SqsRouterHandlerRecordResult struct {
	Result interface{}
	Error  error
}

func (r *SqsRouter) ConvertSqsEventRecordToData(info SqsRecordInfo) ([]byte, error) {
	return []byte(info.Record.Body), nil
}

func (r *SqsRouter) HandleRecord(info SqsRecordInfo) (interface{}, error) {
	handler, queueName, err := r.SelectRouteHandler(info)
	if nil != err {
		return nil, err
	}
	data, err := r.ConvertSqsEventRecordToData(info)
	if nil != err {
		return nil, err
	}
	return handler(data, info, queueName)

}

func (r *SqsRouter) SelectRouteHandler(info SqsRecordInfo) (SqsRouteHandler, string, error) {
	queueName := info.Record.EventSourceARN[strings.LastIndex(info.Record.EventSourceARN, ":")+1:]
	v, found := r.Routes[queueName]
	if found {
		return v, queueName, nil
	}
	v, found = r.Routes["*"]
	if found {
		return v, queueName, nil
	}
	return nil, queueName, errors.New("no sqs route handler found for queue '" + queueName + "'")
}

func (r *SqsRouter) Handle(event events.SQSEvent, ctx context.Context) (interface{}, error) {
	var err error

	if nil != r.BeforeRecords {
		err = r.BeforeRecords(event, ctx)
		if nil != err {
			return nil, err
		}
	}
	var rrs []SqsRouterHandlerRecordResult
	var info SqsRecordInfo
	var result SqsRouterHandlerRecordResult
	for i, record := range event.Records {
		info = SqsRecordInfo{
			RecordIndex: i,
			Record:      record,
			Context:     ctx,
			Event:       event,
		}
		if nil != r.BeforeRecord {
			info, err = r.BeforeRecord(info)
			if nil != err {
				rrs = append(rrs, SqsRouterHandlerRecordResult{
					Result: nil,
					Error:  err,
				})
				continue
			}
		}
		rr, err := r.HandleRecord(info)
		result = SqsRouterHandlerRecordResult{
			Result: rr,
			Error:  err,
		}
		if nil != r.AfterRecord {
			result, err = r.AfterRecord(result, info)
			if nil != err {
				rrs = append(rrs, SqsRouterHandlerRecordResult{
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
