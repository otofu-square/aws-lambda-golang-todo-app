package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type RequestJSON struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func createResponse(body string, status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: body, StatusCode: status}, nil
}

func createErrorResponse(errorMsg string, status int) (events.APIGatewayProxyResponse, error) {
	body := fmt.Sprintf("{\"error\": \"%s\"}", errorMsg)
	return createResponse(body, status)
}

func CreateHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestJSON RequestJSON
	if err := json.Unmarshal([]byte(request.Body), &requestJSON); err != nil {
		return createErrorResponse("Bad Request", 400)
	}

	todo := NewTodo(requestJSON.Title, requestJSON.Completed)
	if err := DynamoDB().Put(todo).Run(); err != nil {
		return createErrorResponse(err.Error(), 502)
	}

	return createResponse(EncodeTodoJSON(todo), 201)
}

func DeleteHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		return createErrorResponse("ID is empty.", 400)
	}

	var todo Todo
	if err := DynamoDB().Delete("ID", id).OldValue(&todo); err != nil {
		return createErrorResponse("Couldn't delete.", 502)
	}

	return createResponse(EncodeTodoJSON(&todo), 200)
}

func IndexHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var todos []Todo
	if err := DynamoDB().Scan().All(&todos); err != nil {
		return createErrorResponse(err.Error(), 404)
	}

	return createResponse(EncodeTodosJSON(&todos), 200)
}

func ShowHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		return createErrorResponse("ID is empty.", 400)
	}

	var todo Todo
	if err := DynamoDB().Get("ID", id).One(&todo); err != nil {
		return createErrorResponse(err.Error(), 404)
	}

	return createResponse(EncodeTodoJSON(&todo), 200)
}

func UpdateHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		return createErrorResponse("ID is empty.", 400)
	}

	var requestJSON RequestJSON
	if err := json.Unmarshal([]byte(request.Body), &requestJSON); err != nil {
		return createErrorResponse("Bad Request.", 400)
	}

	var todo Todo
	err := DynamoDB().
		Update("ID", id).
		Set("Title", requestJSON.Title).
		Set("Completed", requestJSON.Completed).
		Value(&todo)
	if err != nil {
		return createErrorResponse(err.Error(), 502)
	}

	return createResponse(EncodeTodoJSON(&todo), 200)
}
