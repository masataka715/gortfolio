package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type MyEvent struct {
	BucketName string `json:"bucket_name"`
}

func hello(event MyEvent) {
	// NewStaticCredentialsの第３引数はtoken
	creds := credentials.NewStaticCredentials("ACCESS_KEY", "SECRET_ACCESS_KEY", "")
	session, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1")},
	)
	svc := s3.New(session)
	resp, _ := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(event.BucketName),
	})
	log.Println(resp) //ここの結果がテストの実行結果のログ出力やCloudWatchのロググループで確認できる
}

func main() {
	lambda.Start(hello)
}
