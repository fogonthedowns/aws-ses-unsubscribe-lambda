package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type response struct {
	UTC time.Time `json:"utc"`
}

type Lambda struct {
	Session *session.Session
}

func Init(x *Lambda) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-1")})
	if err != nil {
		fmt.Println("NewSession error")
	}
	x.Session = sess
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.QueryStringParameters['email']
	
	if email != "" {
	
	}

	resp := &response{
		UTC: now.UTC(),
	}
	body, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// writeToS3 writes email address to S3
func (x Lambda) writeToS3(emailAddress string) (err error) {
	uploader := s3manager.NewUploader(x.Session)
	currentTime := time.Now()
	id := uuid.New()

	bucketPrefix := fmt.Sprintf(
		"%s/%s/",
		"emailunsubscribe",
		currentTime.Format("2006-01-02"))

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketPrefix),
		Key:    aws.String(id.String()),
		Body:   bytes.NewReader([]byte(emailAddress)),
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
