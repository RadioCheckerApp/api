package main

import (
	"github.com/RadioCheckerApp/api/api-aws/awsutil"
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

func Handler() (events.APIGatewayProxyResponse, error) {
	// AWS config implicitly defined by serverless.yml
	dbSession, _ := session.NewSession(&aws.Config{})

	db := dynamodb.New(dbSession)
	stationDAO := datalayer.NewDDBStationDAO(db, os.Getenv("STATIONS_TABLE"))

	stations, err := shared.Stations(stationDAO)
	responseMessage := model.NewAPIResponseMessage(stations, err)
	return awsutil.CreateResponse(200, responseMessage), nil
}

func main() {
	lambda.Start(Handler)
}
