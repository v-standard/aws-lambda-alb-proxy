package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/v-standard/aws-lambda-alb-proxy/proxy"
)

func APIGatewayHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	params := req.MultiValueQueryStringParameters
	res := events.APIGatewayProxyResponse{
		StatusCode:        200,
		MultiValueHeaders: map[string][]string{"Content-Type": {"text/html"}},
		Body:              fmt.Sprintf("Hello %s", params["name"][0]),
		IsBase64Encoded:   false,
	}
	return res, nil
}

func main() {
	lambda.Start(proxy.ALBProxyWithContext(APIGatewayHandler))
}
