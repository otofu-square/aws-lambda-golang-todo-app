package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestJSON struct {
	Name string `json:"name"`
}

type ResponseJSON struct {
	Message string `json:"message"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestJSON RequestJSON
	if err := json.Unmarshal([]byte(request.Body), &requestJSON); err != nil || requestJSON.Name == "" {
		return events.APIGatewayProxyResponse{Body: "Bad Request", StatusCode: 400}, nil
	}

	greeting := "Hello, " + requestJSON.Name
	responseJSON, _ := json.Marshal(ResponseJSON{Message: greeting})
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
