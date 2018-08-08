package main

import (
	"github.com/RadioCheckerApp/api/api-aws/awsutil"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"github.com/RadioCheckerApp/api/request"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func Handler(apiRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// AWS config implicitly defined by serverless.yml
	dbSession, _ := session.NewSession(&aws.Config{})

	db := dynamodb.New(dbSession)
	trackRecordsDAO := datalayer.NewDDBTrackRecordDAO(db, os.Getenv("TRACKRECORDS_TABLE"))

	worker, err := request.CreateTracksWorker(trackRecordsDAO, apiRequest.PathParameters,
		apiRequest.QueryStringParameters)
	if err != nil {
		responseMessage := model.NewAPIResponseMessage(nil, err)
		return awsutil.CreateResponse(200, responseMessage), nil
	}

	tracks, err := worker.HandleRequest()
	responseMessage := model.NewAPIResponseMessage(tracks, err)
	return awsutil.CreateResponse(200, responseMessage), nil
}

func main() {
	lambda.Start(Handler)
}
