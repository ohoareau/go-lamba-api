package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/adapters"
	"github.com/ohoareau/gola/common"
)

func ConvertPayloadToApiGatewayV2Event(payload []byte) events.APIGatewayV2HTTPRequest {
	var event events.APIGatewayV2HTTPRequest
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err.Error())
	}
	return event
}

func ConvertPayloadToApiGatewayV1Event(payload []byte) events.APIGatewayProxyRequest {
	var event events.APIGatewayProxyRequest
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err.Error())
	}
	return event
}

//goland:noinspection GoUnusedParameter
func HandleApiGatewayV2Event(event events.APIGatewayV2HTTPRequest, ctx context.Context, options *common.Options) (events.APIGatewayV2HTTPResponse, error) {
	return adapters.CreateChiAdapter(CreateHttpRouter(options, false)).ProxyWithContext(ctx, event)
}

//goland:noinspection GoUnusedParameter
func HandleApiGatewayV1Event(event events.APIGatewayProxyRequest, ctx interface{}, options *common.Options) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, errors.New("gola not yet implemented: ApiGatewayV1 Event handler")
}
