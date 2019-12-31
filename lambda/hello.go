package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Music struct {
	Artist    string `dynamo:"Artist"`    //パーティションキー
	SongTitle string `dynamo:"SongTitle"` //ソートキー
}

func hello() {
	// 大きな流れ：認証情報→セッション→dynamoDB
	creds := credentials.NewStaticCredentials("ACCESS_KEY", "SECRET_ACCESS_KEY", "") //第３引数はtoken
	sess, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1")},
	)

	db := dynamo.New(sess)
	table := db.Table("Music")
	// データ入れる
	u := Music{Artist: "イケメン風の人", SongTitle: "バレない程度にパクったメロディ"}
	fmt.Println(u)
	if err := table.Put(u).Run(); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("成功！")
	}
	// データ取得
	var music []Music
	err := table.Get("Artist", "イケメン風の人").All(&music)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(music)
}

func main() {
	lambda.Start(hello)
	// hello()
}
