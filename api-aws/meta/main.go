package main

import (
	"github.com/RadioCheckerApp/api/api-aws/awsutil"
	"github.com/RadioCheckerApp/api/model"
	"github.com/RadioCheckerApp/api/request"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() (events.APIGatewayProxyResponse, error) {
	worker := request.CreateMetaWorker()
	message, err := worker.HandleRequest()
	responseMessage := model.NewAPIResponseMessage(message, err)
	return awsutil.CreateResponse(200, responseMessage), nil
}

func main() {
	lambda.Start(Handler)
}
