package gola

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/adapters"
)

func start(handler ApiGwV2Handler) {
	runtime.Start(handler)
}

func createHandler(configure RouterConfigurator, features Features) ApiGwV2Handler {
	return func(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		r := chi.NewRouter()

		applyFeatures(r, features)

		configure(r)

		adapter := adapters.CreateChiAdapter(r)

		response, err := adapter.ProxyWithContext(context.Background(), event)

		return response, err
	}
}

//goland:noinspection GoUnusedExportedFunction
func Gola(configure RouterConfigurator, featurize FeaturesCreator) {
	start(createHandler(configure, featurize()))
}