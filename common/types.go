package common

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-chi/chi/v5"
)

type ApiGwV2Handler func(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type Apigw2Configurator func(r HttpRouter)
type Apigw1Configurator func(r HttpRouter)
type SnsConfigurator func(r SnsRouter)
type SqsConfigurator func(r SqsRouter)
type S3Configurator func(r S3Router)
type KinesisConfigurator func(r KinesisRouter)
type DynamodbConfigurator func(r DynamodbRouter)

type FeaturesCreator func() Features

type FeatureKey string
type FeatureValue bool

type Features map[FeatureKey]FeatureValue

type HttpRouter = *chi.Mux
type SnsRouter = interface{}
type SqsRouter = interface{}
type S3Router = interface{}
type KinesisRouter interface {
	Handle(event events.KinesisEvent, ctx context.Context) (interface{}, error)
	AddRoute(selector string, route KinesisRouteHandler)
}
type DynamodbRouter = interface{}

type KinesisRecordInfo struct {
	RecordIndex int
	Record      events.KinesisEventRecord
	Event       events.KinesisEvent
	Context     context.Context
}

type KinesisRouteHandler func(data interface{}, info KinesisRecordInfo) (interface{}, error)
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
