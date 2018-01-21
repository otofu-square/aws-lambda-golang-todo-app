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
	var requestJSON RequestJSON
	if err := json.Unmarshal([]byte(request.Body), &requestJSON); err != nil {
		return events.APIGatewayProxyResponse{Body: "Bad Request", StatusCode: 400}, nil
	}

	todo := utils.NewTodo(requestJSON.Title, requestJSON.Completed)
	if err := utils.DynamoDB().Put(todo).Run(); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}, nil
	}

	responseJSON, _ := json.Marshal(todo)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 201}, nil
}

func main() {
	lambda.Start(handleRequest)
}
