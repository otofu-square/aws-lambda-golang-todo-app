package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	utils "github.com/otofu-square/aws-lambda-golang-todo-app/utils"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		errorMsg := "{\"error\": \"ID is empty.\"}"
		return events.APIGatewayProxyResponse{Body: errorMsg, StatusCode: 400}, nil
	}

	var todo utils.Todo
	if err := utils.DynamoDB().Get("ID", id).One(&todo); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	responseJSON, _ := json.Marshal(todo)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
