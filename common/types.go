package common

import (
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
type SqsRouter interface{}
type S3Router interface{}
type DynamodbRouter interface{}
