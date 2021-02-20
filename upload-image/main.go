package main

import (
	"hello-world/aws"
	"hello-world/multipart"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var keys = []string{
		"name",
		"image",
	}

	body, err := multipart.DecodeFromBase64String(request, keys)
	if err != nil {
		return BadRequest(err), nil
	}

	path, err := aws.NewS3().UploadFile(body["image"], "image.png")
	if err != nil {
		return BadRequest(err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       path,
	}, nil
}

//BadRequest - helper function to return bad request error
func BadRequest(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}
}
