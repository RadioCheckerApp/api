package main

import (
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/shared"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// AWS config implicitly defined by serverless.yml
	dbSession, _ := session.NewSession(&aws.Config{})

	db := dynamodb.New(dbSession)
	trackRecordsDAO := datalayer.NewDDBTrackRecordDAO(db, os.Getenv("TRACKRECORDS_TABLE"))

	jsonStr, err := shared.Tracks(trackRecordsDAO, request.PathParameters, request.QueryStringParameters)
	if err != nil {
		log.Printf("Tracks: %s", err)
		return events.APIGatewayProxyResponse{}, errors.New("internal server error")
	}

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       jsonStr,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
