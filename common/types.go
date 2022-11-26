package common

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
)

type Apigw2Configurator func(r *HttpRouter)
type Apigw1Configurator func(r *HttpRouter)
type SnsConfigurator func(r *SnsRouter)
type SqsConfigurator func(r *SqsRouter)
type S3Configurator func(r *S3Router)
type KinesisConfigurator func(r *KinesisRouter)
type DynamodbConfigurator func(r *DynamodbRouter)

type Features map[string]bool

type HttpRouter = chi.Mux
type SnsRouter interface{}
type S3Router interface{}
type DynamodbRouter interface{}

type HandlerWrapper func(h lambda.Handler) lambda.Handler

type Options struct {
	Apigw2Configurator   Apigw2Configurator
	Apigw1Configurator   Apigw1Configurator
	SnsConfigurator      SnsConfigurator
	SqsConfigurator      SqsConfigurator
	S3Configurator       S3Configurator
	KinesisConfigurator  KinesisConfigurator
	DynamodbConfigurator DynamodbConfigurator
	Features             Features
	HandlerWrapper       HandlerWrapper
}
