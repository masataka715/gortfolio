package utils

import (
	"io"
	"log"
	"os"
)

func LoggingSettings(logFile string) {
	// ファイル開き方の設定（O_RDWR:読み書き両方可能、O_CREATE:ファイルがない時に作成、O_APPEND:書き込み方は追記）
	logfile, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// 標準出力に表示＆logfileに書き込みを同時実行(MultiWriter)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	// ログの日時設定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// ログを出力するio.Writerを設定する
	log.SetOutput(multiLogFile)
}
