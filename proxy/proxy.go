package proxy

import (
	"context"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaAdapterWithContext struct {
	f func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func (a *lambdaAdapterWithContext) ProxyWithContext(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	queryStringParameters, err := UnescapeQueryString(req.QueryStringParameters)
	if err != nil {
		return events.ALBTargetGroupResponse{}, err
	}
	multiValueQueryStringParameters, err := UnescapeMultiValueQueryString(req.MultiValueQueryStringParameters)
	if err != nil {
		return events.ALBTargetGroupResponse{}, err
	}
	apigwReq := events.APIGatewayProxyRequest{
		Path:                            req.Path,
		HTTPMethod:                      req.HTTPMethod,
		Headers:                         req.Headers,
		MultiValueHeaders:               req.MultiValueHeaders,
		QueryStringParameters:           queryStringParameters,
		MultiValueQueryStringParameters: multiValueQueryStringParameters,
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

func UnescapeQueryString(qs map[string]string) (map[string]string, error) {
	res := map[string]string{}
	for k, v := range qs {
		key, err := url.QueryUnescape(k)
		if err != nil {
			return nil, err
		}
		value, err := url.QueryUnescape(v)
		if err != nil {
			return nil, err
		}
		res[key] = value
	}
	return res, nil
}

func UnescapeMultiValueQueryString(qs map[string][]string) (map[string][]string, error) {
	res := map[string][]string{}
	for k, v := range qs {
		key, err := url.QueryUnescape(k)
		if err != nil {
			return nil, err
		}
		for i := range v {
			value, err := url.QueryUnescape(v[i])
			if err != nil {
				return nil, err
			}
			res[key] = append(res[key], value)
		}
	}
	return res, nil
}
