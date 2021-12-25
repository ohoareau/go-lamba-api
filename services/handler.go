package services

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ohoareau/gola/common"
)

type lambdaHandler func(ctx context.Context, payload []byte) (interface{}, error)

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

func CreateHandler(options *common.Options) lambda.Handler {
	return lambdaHandler(func(ctx context.Context, payload []byte) (interface{}, error) {
		mode := detectModeFromPayload(payload)
		switch mode {
		case "apigw2":
			return HandleApiGatewayV2Event(ConvertPayloadToApiGatewayV2Event(payload), ctx, options)
		case "kinesis":
			return HandleKinesisEvent(ConvertPayloadToKinesisEvent(payload), ctx, options)
		case "apigw1":
			return HandleApiGatewayV1Event(ConvertPayloadToApiGatewayV1Event(payload), ctx, options)
		case "sqs":
			return HandleSqsEvent(ConvertPayloadToSqsEvent(payload), ctx, options)
		case "s3":
			return HandleS3Event(ConvertPayloadToS3Event(payload), ctx, options)
		case "dynamodb":
			return HandleDynamoDBEvent(ConvertPayloadToDynamoDBEvent(payload), ctx, options)
		case "sns":
			return HandleSnsEvent(ConvertPayloadToSnsEvent(payload), ctx, options)
		default:
			return HandleApiGatewayV2Event(ConvertPayloadToApiGatewayV2Event(payload), ctx, options)
		}
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
