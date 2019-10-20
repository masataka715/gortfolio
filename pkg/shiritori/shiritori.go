package shiritori

import (
	"encoding/base64"
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
	data := map[string]interface{}{}
	if r.Method == http.MethodPost {
		r.ParseForm()
		shiritoriWord := r.FormValue("shiritoriWord")
		lastLetter, message := judge(r, shiritoriWord, "り")
		if lastLetter != "" {
			base64Value := base64.StdEncoding.EncodeToString([]byte(lastLetter))
			http.SetCookie(w, &http.Cookie{
				Name:  "lastLetter",
				Value: base64Value,
				Path:  "/",
			})
			data["lastLetter"] = lastLetter
		}
		data["shiritoriMessage"] = message
	}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/shiritori.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func judge(r *http.Request, word string, correctFirstLetter string) (string, string) {
	var message string
	var lastLetter string
	wordRune := []rune(word)
	firstLetter := string(wordRune[0])

	cookie, err := r.Cookie("lastLetter")
	if err == nil {
		v, _ := base64.StdEncoding.DecodeString(cookie.Value)
		correctFirstLetter = string(v)
	} else {
		correctFirstLetter = "り"
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
		lastLetter = correctFirstLetter
		return lastLetter, message
	}

	// 最初の文字が合っているかの判断
	if firstLetter != correctFirstLetter {
		message = "最初の文字が違います"
		lastLetter = correctFirstLetter
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
