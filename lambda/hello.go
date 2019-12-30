package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func hello() {
	// ~/.aws/credentialsからアクセスキーを読み込む
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sns.New(sess)
	result, _ := svc.ListTopics(nil)
	for _, t := range result.Topics {
		log.Println(*t.TopicArn)
		input := &sns.PublishInput{
			Message:  aws.String("Hello world!"),
			TopicArn: aws.String(*t.TopicArn),
		}
		result, _ := svc.Publish(input)
		log.Println(result)
	}
}

func main() {
	// lambda.Start(hello)
	hello()
}
