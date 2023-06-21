package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Request struct {
	PlantName string `json:"plantName"`
	Value     int    `json:"value"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	decoded, _ := base64.StdEncoding.DecodeString(request.Body)
	formData, err := url.ParseQuery(string(decoded))

	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	plantName := formData.Get("plantName")
	value, err := strconv.Atoi(formData.Get("value"))

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, nil
	}

	req := Request{
		PlantName: plantName,
		Value:     value,
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})
	if err != nil {
		log.Fatal(err)
	}
	snsClient := sns.New(sess)

	topicARN := "arn:aws:sns:us-east-2:730384195053:UpdateTemp"
	payload, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	result, err := snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(string(payload)),
	})
	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(*result.MessageId),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
