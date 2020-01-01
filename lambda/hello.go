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

func insertMessage() {
	// 大きな流れ：認証情報→セッション→sqs→dynamoDB
	creds := credentials.NewStaticCredentials("ACCESS_KEY", "SECRET_ACCESS_KEY", "") //第３引数はtoken

	sess, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1")},
	)

	// SQS
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
		return
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

	mesBody := result.Messages[0].Body
	fmt.Println(*mesBody)

	// dynamoDB
	db := dynamo.New(sess)
	table := db.Table("Music")
	// データ入れる
	u := Music{Artist: "イケメン風の人", SongTitle: *mesBody}
	fmt.Println(u)
	if err := table.Put(u).Run(); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("成功！")
	}
}

func main() {
	// ラムダ実行
	// lambda.Start(insertMessage)
	// ローカルでテストする用
	insertMessage()
}
