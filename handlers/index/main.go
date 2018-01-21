package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	handlers "github.com/otofu-square/aws-lambda-golang-todo-app/handlers"
)

func main() {
	lambda.Start(handlers.IndexHandler)
}
