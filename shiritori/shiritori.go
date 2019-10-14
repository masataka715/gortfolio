package shiritori

import (
	"html/template"
	"net/http"
	"unicode/utf8"
)

type Shiritori struct {
	Word       string
	LastLetter string
	Message    string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// if post
	shiritoriWord := r.FormValue("shiritoriWord")
	var cookie *http.Cookie
	cookie.Value = "test"
	lastLetter, message := judge(r, shiritoriWord, cookie)
	if lastLetter != "" {
		http.SetCookie(w, &http.Cookie{
			Name:  "lastLetter",
			Value: lastLetter,
			Path:  "/",
		})
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "shiritoriMessage",
		Value: message,
		Path:  "/",
	})

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/shiritori.html"))
	_ = templates.ExecuteTemplate(w, "layout", nil)
}

func judge(r *http.Request, word string, correctFirstLetter *http.Cookie) (string, string) {
	var message string
	var lastLetter string
	wordRune := []rune(word)
	firstLetter := string(wordRune[0])

	// OnlyGoでは必要
	correctFirstLetter, _ = r.Cookie("lastLetter")
	if correctFirstLetter.Value == "" {
		correctFirstLetter.Value = "り"
	}
	// 平仮名かどうかの判断
	whiteList := []string{
		"あ", "い", "う", "え", "お",
		"か", "き", "く", "け", "こ",
		"さ", "し", "す", "せ", "そ",
		"た", "ち", "つ", "て", "と",
		"な", "に", "ぬ", "ね", "の",
		"は", "ひ", "ふ", "へ", "ほ",
		"ま", "み", "む", "め", "も",
		"や", "ゆ", "よ",
		"ら", "り", "る", "れ", "ろ",
		"わ", "を", "ん",
		"が", "ぎ", "ぐ", "げ", "ご",
		"ざ", "じ", "ず", "ぜ", "ぞ",
		"だ", "ぢ", "づ", "で", "ど",
		"ば", "び", "ぶ", "べ", "ぼ",
		"ぱ", "ぴ", "ぷ", "ぺ", "ぽ",
		"ぁ", "ぃ", "ぅ", "ぇ", "ぉ",
		"っ", "ゃ", "ゅ", "ょ",
	}
	flag := 0
	for _, v1 := range word {
		for _, v2 := range whiteList {
			if string(v1) == v2 {
				flag++
			}
		}
	}
	if flag != utf8.RuneCountInString(word) {
		message = "平仮名で入力して下さい"
		lastLetter = correctFirstLetter.Value
		return lastLetter, message
	}

	// 最初の文字が合っているかの判断
	if firstLetter != correctFirstLetter.Value {
		message = "最初の文字が違います"
		lastLetter = correctFirstLetter.Value
		return lastLetter, message
	}
	// 最後の文字が「ん」かどうかの判断
	size := len(wordRune)
	lastLetter = string(wordRune[size-1])
	if lastLetter == "ん" {
		message = "最後が「ん」なのでゲームオーバー！　また「り」からです"
		return "り", message
	}

	return lastLetter, message
}
