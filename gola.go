package gola

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/adapters"
	"net/http"
	"os"
)

func start(handler interface{}) {
	runtime.Start(handler)
}

func startLocal(r Router) {
	port, hasPort := os.LookupEnv("PORT")
	if !hasPort {
		port = "5000"
	}

	address := ":" + port
	url := "http://localhost" + address
	fmt.Println("🚀 Server ready at " + url)
	err := http.ListenAndServe(address, r)

	if nil != err {
		fmt.Println(err.Error())
	}
}

func createRouter(options Options) Router {
	r := chi.NewRouter()

	applyFeatures(r, options.Features)

	options.Configure(r)

	return r
}

func createHandler(options Options) interface{} {
	return func(ctx interface{}, event interface{}) (interface{}, error) {
		mode := detectModeFromEvent(event)
		switch mode {
		case "apigw2":
			return handleApiGatewayV2Event(event.(events.APIGatewayV2HTTPRequest), ctx, options)
		case "kinesis":
			return handleKinesisEvent(event.(events.KinesisEvent), ctx, options)
		case "apigw1":
			return handleApiGatewayV1Event(event.(events.APIGatewayProxyRequest), ctx, options)
		case "sqs":
			return handleSqsEvent(event.(events.SQSEvent), ctx, options)
		case "s3":
			return handleS3Event(event.(events.S3Event), ctx, options)
		case "dynamodb":
			return handleDynamoDBEvent(event.(events.DynamoDBEvent), ctx, options)
		case "sns":
			return handleSnsEvent(event.(events.SNSEvent), ctx, options)
		default:
			return handleApiGatewayV2Event(event.(events.APIGatewayV2HTTPRequest), ctx, options)
		}
	}
}

func detectModeFromEvent(event interface{}) string {
	switch event.(type) {
	case events.APIGatewayV2HTTPRequest:
		return "apigw2"
	case events.KinesisEvent:
		return "kinesis"
	case events.SQSEvent:
		return "sqs"
	case events.DynamoDBEvent:
		return "dynamodb"
	case events.S3Event:
		return "s3"
	case events.APIGatewayProxyRequest:
		return "apigw1"
	case events.SNSEvent:
		return "sns"
	default:
		return "unknown"
	}
}

//goland:noinspection GoUnusedParameter
func handleApiGatewayV2Event(event events.APIGatewayV2HTTPRequest, ctx interface{}, options Options) (events.APIGatewayV2HTTPResponse, error) {
	r := createRouter(options)

	adapter := adapters.CreateChiAdapter(r)

	response, err := adapter.ProxyWithContext(context.Background(), event)

	return response, err
}

//goland:noinspection GoUnusedParameter
func handleApiGatewayV1Event(event events.APIGatewayProxyRequest, ctx interface{}, options Options) (events.APIGatewayProxyResponse, error) {
	err := errors.New("gola not yet implemented: ApiGatewayV1 Event handler")
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, err
}

//goland:noinspection GoUnusedParameter
func handleSqsEvent(event events.SQSEvent, ctx interface{}, options Options) (interface{}, error) {
	err := errors.New("gola not yet implemented: SQS Event handler")
	return nil, err
}

//goland:noinspection GoUnusedParameter
func handleSnsEvent(event events.SNSEvent, ctx interface{}, options Options) (interface{}, error) {
	err := errors.New("gola not yet implemented: SNS Event handler")
	return nil, err
}

//goland:noinspection GoUnusedParameter
func handleS3Event(event events.S3Event, ctx interface{}, options Options) (interface{}, error) {
	err := errors.New("gola not yet implemented: S3 Event handler")
	return nil, err
}

//goland:noinspection GoUnusedParameter
func handleDynamoDBEvent(event events.DynamoDBEvent, ctx interface{}, options Options) (interface{}, error) {
	err := errors.New("gola not yet implemented: DynamoDB Event handler")
	return nil, err
}

//goland:noinspection GoUnusedParameter
func handleKinesisEvent(event events.KinesisEvent, ctx interface{}, options Options) (interface{}, error) {
	err := errors.New("gola not yet implemented: Kinesis Event handler")
	return nil, err
}

//goland:noinspection GoUnusedExportedFunction
func Gola(configure RouterConfigurator, featurize FeaturesCreator) {
	_, isAwsDeployed := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	options := Options{
		Configure: configure,
		Features:  featurize(),
	}
	if !isAwsDeployed {
		startLocal(createRouter(options))
	} else {
		start(createHandler(options))
	}
}
