package main

import (
	"encoding/json"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"github.com/RadioCheckerApp/api/shared"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// AWS config implicitly defined by serverless.yml
	dbSession, _ := session.NewSession(&aws.Config{})

	db := dynamodb.New(dbSession)
	trackRecordsDAO := datalayer.NewDDBTrackRecordDAO(db, os.Getenv("TRACKRECORDS_TABLE"))

	tracks, err := shared.Tracks(trackRecordsDAO, request.PathParameters,
		request.QueryStringParameters)
	responseMessage := model.NewAPIResponseMessage(tracks, err)
	return createResponse(200, responseMessage), nil
}

func main() {
	lambda.Start(Handler)
}

func createResponse(statusCode int, message model.APIResponseMessage) events.
	APIGatewayProxyResponse {
	encodedMessage, _ := json.Marshal(message)
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(encodedMessage),
		StatusCode: statusCode,
	}
}