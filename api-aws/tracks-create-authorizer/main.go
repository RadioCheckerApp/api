package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"strings"
)

const crawlerPrincipalID = "crawler-service"

func Handler(authRequest events.APIGatewayCustomAuthorizerRequest) (events.
	APIGatewayCustomAuthorizerResponse, error) {
	if authRequest.AuthorizationToken == "" {
		return statusUnauthorized(&authRequest)
	}

	split := strings.Split(authRequest.AuthorizationToken, "Bearer")
	if len(split) != 2 {
		return statusUnauthorized(&authRequest)
	}

	token := strings.TrimSpace(split[1])

	if strings.ToLower(token) == os.Getenv("TRACKS_CREATE_AUTH_TOKEN") {
		return statusAuthorized(&authRequest)
	}

	return statusUnauthorized(&authRequest)
}

func main() {
	lambda.Start(Handler)
}

func statusUnauthorized(authRequest *events.APIGatewayCustomAuthorizerRequest) (events.
	APIGatewayCustomAuthorizerResponse, error) {
	log.Printf("UNAUTHORIZE REQUEST: Type: `%s`, Token: `%s`, ARN: `%s`",
		authRequest.Type, authRequest.AuthorizationToken, authRequest.MethodArn)

	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
}

func statusAuthorized(authRequest *events.APIGatewayCustomAuthorizerRequest) (events.
	APIGatewayCustomAuthorizerResponse, error) {
	wildcardArn := buildWildcardResourceArn(authRequest.MethodArn)
	log.Printf("AUTHORIZE REQUEST: Type: `%s`, Token: `%s`, ARN: `%s`",
		authRequest.Type, authRequest.AuthorizationToken, wildcardArn)

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:    crawlerPrincipalID,
		PolicyDocument: generatePolicy("Allow", wildcardArn),
	}, nil
}

func generatePolicy(effect, resourceArn string) events.APIGatewayCustomAuthorizerPolicy {
	statement := events.IAMPolicyStatement{
		Action:   []string{"execute-api:Invoke"}, // default action
		Effect:   effect,
		Resource: []string{resourceArn},
	}

	return events.APIGatewayCustomAuthorizerPolicy{
		Version:   "2012-10-17", // default version
		Statement: []events.IAMPolicyStatement{statement},
	}
}

func buildWildcardResourceArn(resourceArn string) string {
	// resource ARN example layout:
	// arn:aws:execute-api:eu-central-1:001975686909:pul5mro035/dev/PUT/stations/hitradio-oe3/tracks/1537701181
	split := strings.Split(resourceArn, "/")
	if len(split) != 7 {
		log.Printf("ERROR: Unable to split ARN `%s`. Not adding any wildcards.", resourceArn)
		return resourceArn
	}
	split[4] = "*"
	split[6] = "*"
	return strings.Join(split, "/")
}
