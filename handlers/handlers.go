package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type RequestJSON struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func CreateHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestJSON RequestJSON
	if err := json.Unmarshal([]byte(request.Body), &requestJSON); err != nil {
		return events.APIGatewayProxyResponse{Body: "Bad Request", StatusCode: 400}, nil
	}

	todo := NewTodo(requestJSON.Title, requestJSON.Completed)
	if err := DynamoDB().Put(todo).Run(); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}, nil
	}

	responseJSON, _ := json.Marshal(todo)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 201}, nil
}

func DeleteHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		errorMsg := "{\"error\": \"ID is empty.\"}"
		return events.APIGatewayProxyResponse{Body: errorMsg, StatusCode: 400}, nil
	}

	var todo Todo
	if err := DynamoDB().Delete("ID", id).OldValue(&todo); err != nil {
		return events.APIGatewayProxyResponse{Body: "Couldn't delete.", StatusCode: 502}, nil
	}

	responseJSON, _ := json.Marshal(todo)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}

func IndexHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var todos []Todo
	if err := DynamoDB().Scan().All(&todos); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	responseJSON, _ := json.Marshal(todos)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}

func ShowHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		errorMsg := "{\"error\": \"ID is empty.\"}"
		return events.APIGatewayProxyResponse{Body: errorMsg, StatusCode: 400}, nil
	}

	var todo Todo
	if err := DynamoDB().Get("ID", id).One(&todo); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	responseJSON, _ := json.Marshal(todo)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}

func UpdateHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		errorMsg := "{\"error\": \"ID is empty.\"}"
		return events.APIGatewayProxyResponse{Body: errorMsg, StatusCode: 400}, nil
	}

	var requestJSON RequestJSON
	if err := json.Unmarshal([]byte(request.Body), &requestJSON); err != nil {
		return events.APIGatewayProxyResponse{Body: "Bad Request", StatusCode: 400}, nil
	}

	var todo Todo
	err := DynamoDB().
		Update("ID", id).
		Set("Title", requestJSON.Title).
		Set("Completed", requestJSON.Completed).
		Value(&todo)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}, nil
	}

	responseJSON, _ := json.Marshal(todo)
	return events.APIGatewayProxyResponse{Body: string(responseJSON), StatusCode: 200}, nil
}
