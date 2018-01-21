package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	utils "github.com/otofu-square/aws-lambda-golang-todo-app/utils"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var todos []utils.Todo
	if err := utils.DynamoDB().Scan().All(&todos); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	responseJSON, _ := json.Marshal(todos)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
