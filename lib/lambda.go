package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type response struct {
	Msg string `json:"email"`
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

func (x Lambda) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.QueryStringParameters["email"]

	if email != "" {
		err := x.writeToS3(email)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
	} else {
		return events.APIGatewayProxyResponse{}, errors.New("email blank")
	}

	message := fmt.Sprintf("%v succesfully unsubscribed", email)
	resp := &response{
		Msg: message,
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
