package main

import (
	"github.com/RadioCheckerApp/api/shared"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       shared.APIMetadata(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
