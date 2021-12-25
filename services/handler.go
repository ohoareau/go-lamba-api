package services

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
)

func CreateHandler(options common.Options) interface{} {
	return func(ctx context.Context, event interface{}) (interface{}, error) {
		mode := detectModeFromEvent(event)
		switch mode {
		case "apigw2":
			return HandleApiGatewayV2Event(event.(events.APIGatewayV2HTTPRequest), ctx, options)
		case "kinesis":
			return HandleKinesisEvent(event.(events.KinesisEvent), ctx, options)
		case "apigw1":
			return HandleApiGatewayV1Event(event.(events.APIGatewayProxyRequest), ctx, options)
		case "sqs":
			return HandleSqsEvent(event.(events.SQSEvent), ctx, options)
		case "s3":
			return HandleS3Event(event.(events.S3Event), ctx, options)
		case "dynamodb":
			return HandleDynamoDBEvent(event.(events.DynamoDBEvent), ctx, options)
		case "sns":
			return HandleSnsEvent(event.(events.SNSEvent), ctx, options)
		default:
			return HandleApiGatewayV2Event(event.(events.APIGatewayV2HTTPRequest), ctx, options)
		}
	}
}

func detectModeFromEvent(event interface{}) string {
	if IsApiGatewayV2Event(event) {
		return "apigw2"
	}
	if IsKinesisEvent(event) {
		return "kinesis"
	}
	if IsSqsEvent(event) {
		return "sqs"
	}
	if IsDynamoDBEvent(event) {
		return "dynamodb"
	}
	if IsS3Event(event) {
		return "s3"
	}
	if IsApiGatewayV1Event(event) {
		return "apigw1"
	}
	if IsSnsEvent(event) {
		return "sns"
	}
	return "unknown"
}
