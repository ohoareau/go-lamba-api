package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ohoareau/gola/common"
	"log"
)

type lambdaHandler func(ctx context.Context, payload []byte) ([]byte, error)

// Invoke calls the handler, and serializes the response.
// If the underlying handler returned an error, or an error occurs during serialization, error is returned.
func (handler lambdaHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	response, err := handler(ctx, payload)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

func CreateHandler(options common.Options) lambda.Handler {
	return lambdaHandler(func(ctx context.Context, payload []byte) ([]byte, error) {
		mode := detectModeFromPayload(payload)
		log.Println(mode)
		var r interface{}
		var err error
		switch mode {
		case "apigw2":
			r, err = HandleApiGatewayV2Event(ConvertPayloadToApiGatewayV2Event(payload), ctx, options)
		case "kinesis":
			r, err = HandleKinesisEvent(ConvertPayloadToKinesisEvent(payload), ctx, options)
		case "apigw1":
			r, err = HandleApiGatewayV1Event(ConvertPayloadToApiGatewayV1Event(payload), ctx, options)
		case "sqs":
			r, err = HandleSqsEvent(ConvertPayloadToSqsEvent(payload), ctx, options)
		case "s3":
			r, err = HandleS3Event(ConvertPayloadToS3Event(payload), ctx, options)
		case "dynamodb":
			r, err = HandleDynamoDBEvent(ConvertPayloadToDynamoDBEvent(payload), ctx, options)
		case "sns":
			r, err = HandleSnsEvent(ConvertPayloadToSnsEvent(payload), ctx, options)
		default:
			r, err = HandleApiGatewayV2Event(ConvertPayloadToApiGatewayV2Event(payload), ctx, options)
		}
		var output []byte
		if nil != err {
			fmt.Println(err)
		} else {
			output, err = json.Marshal(r)
			if nil != err {
				fmt.Println(err)
			}
		}
		return output, err
	})
}

type BasicBulkEventRecord struct {
	EventSource1 string `json:"EventSource"`
	EventSource2 string `json:"eventSource"`
}

type BulkEvent struct {
	Records []BasicBulkEventRecord `json:"records"`
}

type SingleEvent struct {
	Version string `json:"version"`
}

func detectModeFromPayload(payload []byte) string {
	var bulkEvent BulkEvent
	isBulkEvent := false
	err := json.Unmarshal(payload, &bulkEvent)
	if nil == err && (0 < len(bulkEvent.Records)) {
		isBulkEvent = true
	}
	if isBulkEvent {
		if IsBulkEventFromSource(bulkEvent, "aws:kinesis") {
			return "kinesis"
		}
		if IsBulkEventFromSource(bulkEvent, "aws:sqs") {
			return "sqs"
		}
		if IsBulkEventFromSource(bulkEvent, "aws:dynamodb") {
			return "dynamodb"
		}
		if IsBulkEventFromSource(bulkEvent, "aws:s3") {
			return "s3"
		}
		if IsBulkEventFromSource(bulkEvent, "aws:sns") {
			return "sns"
		}
		return "unknown"
	}
	var singleEvent SingleEvent
	err = json.Unmarshal(payload, &singleEvent)
	if nil == err {
		if "2.0" == singleEvent.Version {
			return "apigw2"
		}
		if "1.0" == singleEvent.Version {
			return "apigw1"
		}
	}
	return "unknown"
}

func IsBulkEventFromSource(bulkEvent BulkEvent, source string) bool {
	return (source == bulkEvent.Records[0].EventSource1) || (source == bulkEvent.Records[0].EventSource2)
}
