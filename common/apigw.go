package common

import "github.com/aws/aws-lambda-go/events"

type ApiGwV2Handler func(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
