package proxy

import (
	"context"
	"net/url"
	"strings"

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
	var hostname string
	if host, ok := req.Headers["host"]; ok && len(host) > 0 {
		hostname = host
	}
	if host, ok := req.MultiValueHeaders["host"]; ok && len(host) > 0 && len(host[0]) > 0 {
		hostname = host[0]
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
		RequestContext: events.APIGatewayProxyRequestContext{
			DomainName:   hostname,
			DomainPrefix: strings.Split(hostname, ".")[0],
		},
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
