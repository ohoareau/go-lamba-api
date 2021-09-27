package gola

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/adapters"
	"net/http"
	"os"
)

func start(handler ApiGwV2Handler) {
	runtime.Start(handler)
}

func startLocal(r Router) {
	port, hasPort := os.LookupEnv("PORT")
	if !hasPort {
		port = "6000"
	}
	err := http.ListenAndServe(":" + port, r)

	if nil != err {
		fmt.Println(err.Error())
	}
}

func createRouter(configure RouterConfigurator, features Features) Router {
	r := chi.NewRouter()

	applyFeatures(r, features)

	configure(r)

	return r
}

func createHandler(configure RouterConfigurator, features Features) ApiGwV2Handler {
	return func(event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		r := createRouter(configure, features)

		adapter := adapters.CreateChiAdapter(r)

		response, err := adapter.ProxyWithContext(context.Background(), event)

		return response, err
	}
}

//goland:noinspection GoUnusedExportedFunction
func Gola(configure RouterConfigurator, featurize FeaturesCreator) {
	_, isLocal := os.LookupEnv("GOLA_LOCAL")
	if isLocal {
		startLocal(createRouter(configure, featurize()))
	} else {
		start(createHandler(configure, featurize()))
	}
}
