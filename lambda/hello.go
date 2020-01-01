package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/guregu/dynamo"
)

type Music struct {
	Artist    string `dynamo:"Artist"`    //パーティションキー
	SongTitle string `dynamo:"SongTitle"` //ソートキー
}

func getMessage(sess *session.Session) *string {
	svc := sqs.New(sess)
	qURL := "https://sqs.ap-northeast-1.amazonaws.com/607012455302/dead.fifo"

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return nil
	}

	mesBody := result.Messages[0].Body
	return mesBody
}

func dynamoInsert(sess *session.Session, mesBody *string) {
	db := dynamo.New(sess)
	table := db.Table("Music")
	// データ入れる
	music := Music{Artist: "イケメン風の人", SongTitle: *mesBody}
	fmt.Println(music)
	if err := table.Put(music).Run(); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Success!")
	}
}

func insertMessage() {
	// 大きな流れ：認証情報→セッション→SQS→dynamoDB
	creds := credentials.NewStaticCredentials("ACCESS_KEY", "SECRET_ACCESS_KEY", "")

	sess, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1")},
	)

	// SQS
	mesBody := getMessage(sess)
	// dynamoDB
	if mesBody != nil {
		dynamoInsert(sess, mesBody)
	}
}

func main() {
	// ラムダ実行
	// lambda.Start(insertMessage)
	// ローカルでテストする用
	insertMessage()
}
