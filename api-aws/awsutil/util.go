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
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(encodedMessage),
		StatusCode: statusCode,
	}
}
