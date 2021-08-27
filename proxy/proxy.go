package proxy

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaAdapterWithContext struct {
	f func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func (a *lambdaAdapterWithContext) ProxyWithContext(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	apigwReq := events.APIGatewayProxyRequest{
		Path:                            req.Path,
		HTTPMethod:                      req.HTTPMethod,
		Headers:                         req.Headers,
		MultiValueHeaders:               req.MultiValueHeaders,
		QueryStringParameters:           req.QueryStringParameters,
		MultiValueQueryStringParameters: req.MultiValueQueryStringParameters,
		Body:                            req.Body,
		IsBase64Encoded:                 req.IsBase64Encoded,
	}
	apigwRes, err := a.f(ctx, apigwReq)
	return events.ALBTargetGroupResponse{
		StatusCode:        apigwRes.StatusCode,
		Headers:           apigwRes.Headers,
		MultiValueHeaders: apigwRes.MultiValueHeaders,
		Body:              apigwRes.Body,
		IsBase64Encoded:   apigwRes.IsBase64Encoded,
	}, err
}

func ALBProxyWithContext(f func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	adapter := lambdaAdapterWithContext{f: f}
	return adapter.ProxyWithContext
}
