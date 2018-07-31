package main

import (
	"github.com/RadioCheckerApp/api/api-aws/awsutil"
	"github.com/RadioCheckerApp/api/model"
	"github.com/RadioCheckerApp/api/shared"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() (events.APIGatewayProxyResponse, error) {
	responseMessage := model.NewAPIResponseMessage(shared.APIMetadata(), nil)
	return awsutil.CreateResponse(200, responseMessage), nil
}

func main() {
	lambda.Start(Handler)
}
