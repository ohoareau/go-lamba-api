package adapters

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/go-chi/chi/v5"
)

// ChiAdapter makes it easy to send API Gateway proxy events to a Chi
// Mux. The library transforms the proxy event into an HTTP request and then
// creates a proxy response object from the http.ResponseWriter
type ChiAdapter struct {
	core.RequestAccessorV2
	router *chi.Mux
}

// CreateChiAdapter creates a new instance of the ChiAdapter object.
// Receives an initialized *chi.Mux object - normally created with chi.NewRouter().
// It returns the initialized instance of the ChiAdapter object.
func CreateChiAdapter(router *chi.Mux) *ChiAdapter {
	return &ChiAdapter{router: router}
}

// Proxy receives an API Gateway proxy event, transforms it into an http.Request
// object, and sends it to the chi.Mux for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (g *ChiAdapter) Proxy(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	chiReq, err := g.ProxyEventToHTTPRequest(req)
	return g.proxyInternal(chiReq, err)
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the chi.Mux for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (g *ChiAdapter) ProxyWithContext(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	chiReq, err := g.EventToRequestWithContext(ctx, req)
	return g.proxyInternal(chiReq, err)
}

func (g *ChiAdapter) proxyInternal(chiReq *http.Request, err error) (events.APIGatewayV2HTTPResponse, error) {

	if err != nil {
		return core.GatewayTimeoutV2(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriterV2()
	g.router.ServeHTTP(http.ResponseWriter(respWriter), chiReq)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeoutV2(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}
