package awsutil

import (
	"encoding/json"
	"github.com/RadioCheckerApp/api/model"
	"github.com/aws/aws-lambda-go/events"
)

func CreateResponse(statusCode int, message model.APIResponseMessage) events.
	APIGatewayProxyResponse {
	encodedMessage, _ := json.Marshal(message)
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
			// required for CORS support,
			// see https://github.com/serverless/serverless/issues/1955#issuecomment-266235353
			"Access-Control-Allow-Origin": "*",
		},
		Body:       string(encodedMessage),
		StatusCode: statusCode,
	}
}
