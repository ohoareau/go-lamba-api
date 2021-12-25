package services

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/adapters"
	"github.com/ohoareau/gola/common"
)

func IsApiGatewayV2Event(event interface{}) bool {
	switch event.(type) {
	case events.APIGatewayV2HTTPRequest:
		return true
	default:
		return false
	}
}

func IsApiGatewayV1Event(event interface{}) bool {
	switch event.(type) {
	case events.APIGatewayProxyRequest:
		return true
	default:
		return false
	}
}

//goland:noinspection GoUnusedParameter
func HandleApiGatewayV2Event(event events.APIGatewayV2HTTPRequest, ctx interface{}, options common.Options) (events.APIGatewayV2HTTPResponse, error) {
	return adapters.CreateChiAdapter(CreateHttpRouter(options)).ProxyWithContext(context.Background(), event)
}

//goland:noinspection GoUnusedParameter
func HandleApiGatewayV1Event(event events.APIGatewayProxyRequest, ctx interface{}, options common.Options) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, errors.New("gola not yet implemented: ApiGatewayV1 Event handler")
}
