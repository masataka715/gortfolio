// package shiritori

// import (
// 	"gortfolio/config"
// 	"gortfolio/pkg/common"
// 	"unicode/utf8"

// 	"github.com/gin-gonic/gin"
// )

// type Shiritori struct {
// 	Word       string `json:"word"`
// 	LastLetter string `json:"lastLetter"`
// 	Message    string `json:"message"`
// }

// func OnlyGo(router *gin.Engine) {
// 	router.POST("/shiritori", func(ctx *gin.Context) {
// 		shiritoriWord := ctx.PostForm("shiritoriWord")
// 		lastLetter, message := judge(ctx, shiritoriWord, "test")
// 		if lastLetter != "" {
// 			ctx.SetCookie("lastLetter", lastLetter, 30, "/", "localhost", false, true)
// 		}
// 		ctx.SetCookie("shiritoriMessage", message, 30, "/", "localhost", false, true)
// 		ctx.Redirect(302, "/")
// 	})
// }

// func judge(ctx *gin.Context, word string, correctFirstLetter string) (string, string) {
// 	var message string
// 	var lastLetter string
// 	wordRune := []rune(word)
// 	firstLetter := string(wordRune[0])

// 	// OnlyGoでは必要
// 	correctFirstLetter, _ := ctx.Cookie("lastLetter")
// 	if correctFirstLetter == "" {
// 		correctFirstLetter = "り"
// 	}
// 	// 平仮名かどうかの判断
// 	whiteList := []string{
// 		"あ", "い", "う", "え", "お",
// 		"か", "き", "く", "け", "こ",
// 		"さ", "し", "す", "せ", "そ",
// 		"た", "ち", "つ", "て", "と",
// 		"な", "に", "ぬ", "ね", "の",
// 		"は", "ひ", "ふ", "へ", "ほ",
// 		"ま", "み", "む", "め", "も",
// 		"や", "ゆ", "よ",
// 		"ら", "り", "る", "れ", "ろ",
// 		"わ", "を", "ん",
// 		"が", "ぎ", "ぐ", "げ", "ご",
// 		"ざ", "じ", "ず", "ぜ", "ぞ",
// 		"だ", "ぢ", "づ", "で", "ど",
// 		"ば", "び", "ぶ", "べ", "ぼ",
// 		"ぱ", "ぴ", "ぷ", "ぺ", "ぽ",
// 		"ぁ", "ぃ", "ぅ", "ぇ", "ぉ",
// 		"っ", "ゃ", "ゅ", "ょ",
// 	}
// 	flag := 0
// 	for _, v1 := range word {
// 		for _, v2 := range whiteList {
// 			if string(v1) == v2 {
// 				flag++
// 			}
// 		}
// 	}
// 	if flag != utf8.RuneCountInString(word) {
// 		message = "平仮名で入力して下さい"
// 		lastLetter = correctFirstLetter
// 		return lastLetter, message
// 	}

// 	// 最初の文字が合っているかの判断
// 	if firstLetter != correctFirstLetter {
// 		message = "最初の文字が違います"
// 		lastLetter = correctFirstLetter
// 		return lastLetter, message
// 	}
// 	// 最後の文字が「ん」かどうかの判断
// 	size := len(wordRune)
// 	lastLetter = string(wordRune[size-1])
// 	if lastLetter == "ん" {
// 		message = "最後が「ん」なのでゲームオーバー！　また「り」からです"
// 		return "り", message
// 	}

// 	return lastLetter, message
// }
