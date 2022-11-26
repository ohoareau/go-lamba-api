package services

import (
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/common"
	"github.com/ohoareau/gola/features"
)

func CreateHttpRouter(options *common.Options, local bool) *common.HttpRouter {
	r := chi.NewRouter()

	if nil != options.HttpRouterConfigurator {
		options.HttpRouterConfigurator(r)
	}

	if local && nil != options.LocalHttpRouterConfigurator {
		options.LocalHttpRouterConfigurator(r)
	}

	applyHttpFeatures(r, options.Features)

	if nil != options.Apigw2Configurator {
		options.Apigw2Configurator(r)
	}
	if nil != options.Apigw1Configurator {
		options.Apigw1Configurator(r)
	}

	return r
}

func applyHttpFeatures(r *common.HttpRouter, f common.Features) {
	if f["logger"] {
		features.Logger(r)
	}
	if f["recoverer"] {
		features.Recoverer(r)
	}
	if f["cors"] {
		features.Cors(r)
	}
	if f["ping"] {
		features.Ping(r)
	}
	if f["jwtauth"] {
		features.JwtAuth(r)
	}
	if f["profiler"] {
		features.Profiler(r)
	}

	// these features add routes and must be executed after all the features that add middlewares

	if f["robots"] {
		features.Robots(r)
	}
	if f["sitemap"] {
		features.Sitemap(r)
	}
	if f["root"] {
		features.Root(r)
	}
	if f["notfound"] {
		features.Notfound(r)
	}
}
