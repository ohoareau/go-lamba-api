package gola

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-chi/chi/v5"
)

type ApiGwV2Handler func(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type RouterConfigurator func(r Router)

type FeatureKey string
type FeatureValue bool

type Features map[FeatureKey]FeatureValue

type Router = *chi.Mux
