package awsutil

import (
	"github.com/RadioCheckerApp/api/model"
	"github.com/aws/aws-lambda-go/events"
	"reflect"
	"testing"
)

type test struct {
	Payload string `json:"payload"`
}

func TestCreateResponse(t *testing.T) {
	var tests = []struct {
		inputStatusCode  int
		inputMessage     model.APIResponseMessage
		expectedResponse events.APIGatewayProxyResponse
	}{
		{
			200,
			model.APIResponseMessage{true, test{"hello world"}, ""},
			events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"success\":true,\"data\":{\"payload\":\"hello world\"}}",
				StatusCode: 200,
			},
		},
		{
			200,
			model.APIResponseMessage{false, nil, "errormsg"},
			events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"success\":false,\"message\":\"errormsg\"}",
				StatusCode: 200,
			},
		},
	}

	for _, test := range tests {
		response := CreateResponse(test.inputStatusCode, test.inputMessage)
		if !reflect.DeepEqual(response, test.expectedResponse) {
			t.Errorf("TestCreateResponse(%v, %v): got: %v, expected: %v",
				test.inputStatusCode, test.inputMessage, response, test.expectedResponse)
		}
	}
}
