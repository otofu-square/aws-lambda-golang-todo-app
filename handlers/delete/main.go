package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	utils "github.com/otofu-square/aws-lambda-golang-todo-app/utils"
)

type RequestJSON struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		errorMsg := "{\"error\": \"ID is empty.\"}"
		return events.APIGatewayProxyResponse{Body: errorMsg, StatusCode: 400}, nil
	}

	var todo utils.Todo
	if err := utils.DynamoDB().Delete("ID", id).OldValue(&todo); err != nil {
		return events.APIGatewayProxyResponse{Body: "Couldn't delete.", StatusCode: 502}, nil
	}

	responseJSON, _ := json.Marshal(todo)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
