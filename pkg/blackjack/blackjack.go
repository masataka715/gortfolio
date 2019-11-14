package blackjack

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gortfolio/pkg/footprint"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	when := time.Now().Format("2006年01月02日 15時04分")
	footprint.Insert("ブラックジャック", when)

	data := map[string]interface{}{}
	data["Blackjack"] = GetAll()

	if r.Method == http.MethodPost {
		r.ParseForm()
		stanted := r.FormValue("stanted")
		deck := NewDeck()
		ShuffleDeck(deck)

		// プレイヤー
		var player Player
		BlackPlayerCookie, err := r.Cookie("blackPlayer")
		player = NewPlayer("あなた")
		if err == nil {
			bytes, _ := base64.StdEncoding.DecodeString(BlackPlayerCookie.Value)
			json.Unmarshal(bytes, &player)
		}
		player.Stanted, _ = strconv.ParseBool(stanted)
		if player.Stanted == false {
			PlayCard(&player, &deck)
		}
		data["player"] = player

		// ディーラー
		var dealer Player
		BlackDealerCookie, err := r.Cookie("blackDealer")
		dealer = NewPlayer("ディーラー")
		if err == nil {
			bytes, _ := base64.StdEncoding.DecodeString(BlackDealerCookie.Value)
			json.Unmarshal(bytes, &dealer)
		}
		if dealer.Stanted == false && player.Score < 22 {
			if dealer.Score > 16 {
				dealer.Stanted = true
				data["mainMessage"] = "ディーラーは引き終えました"
			} else {
				PlayCard(&dealer, &deck)
			}
		}
		data["dealer"] = dealer

		SetBlackjackCookie(w, r, player, dealer)
		if player.Score > 21 {
			data["mainMessage"] = player.Name + "は21を超えてしまいました"
			data["finishMessage"] = dealer.Name + "の勝ちです"

			renewalPlayer := RenewalPlayer(player, "lose")
			data["victoryDefeat"] = VictoryDefeatMes(renewalPlayer)
			RenewBlackjackCookie(w, renewalPlayer)
		}
		if dealer.Score > 21 {
			data["mainMessage"] = dealer.Name + "は21を超えてしまいました"
			data["finishMessage"] = player.Name + "の勝ちです"

			renewalPlayer := RenewalPlayer(player, "win")
			data["victoryDefeat"] = VictoryDefeatMes(renewalPlayer)
			RenewBlackjackCookie(w, renewalPlayer)
		}
		if player.Stanted == true && dealer.Stanted == true {
			data["mainMessage"] = "どちらも引き終えました"
			var renewalPlayer Player
			if player.Score > dealer.Score {
				data["finishMessage"] = player.Name + "の勝ちです"
				renewalPlayer = RenewalPlayer(player, "win")
			} else if player.Score == dealer.Score {
				data["finishMessage"] = "引き分けです"
				renewalPlayer = RenewalPlayer(player, "win")
			} else {
				data["finishMessage"] = player.Name + "の負けです"
				renewalPlayer = RenewalPlayer(player, "lose")
			}
			data["victoryDefeat"] = VictoryDefeatMes(renewalPlayer)
			RenewBlackjackCookie(w, renewalPlayer)
		}
	}

	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/blackjack.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func PlayCard(player *Player, deck *[]Card) Card {
	card := DrawCard(player, deck)
	CalcScore(player)
	return card
}

// A(1) = 1, J(11)〜K(13) = 10, それ以外は数字どおり
func CardScore(card Card) int {
	score := 0
	switch card.Number {
	case 11, 12, 13:
		score = 10
	default:
		score = card.Number
	}
	return score
}

func CalcScore(player *Player) {
	score := 0
	n := len(player.Cards)
	for i := 0; i < n; i++ {
		score += CardScore(player.Cards[i])
	}
	player.Score = score
	player.ScoreMessage = PrintScore(*player)
}

func PrintScore(player Player) string {
	return player.Name + "の得点は" + fmt.Sprintf("%v", player.Score) + "点です。"
}

func SetBlackjackCookie(w http.ResponseWriter, r *http.Request, player Player, dealer Player) {
	byte, _ := json.Marshal(player)
	http.SetCookie(w, &http.Cookie{
		Name:  "blackPlayer",
		Value: base64.StdEncoding.EncodeToString(byte),
		Path:  "/",
	})

	byte, _ = json.Marshal(dealer)
	http.SetCookie(w, &http.Cookie{
		Name:  "blackDealer",
		Value: base64.StdEncoding.EncodeToString(byte),
		Path:  "/",
	})
}

func RenewBlackjackCookie(w http.ResponseWriter, player Player) {
	byte, _ := json.Marshal(player)
	http.SetCookie(w, &http.Cookie{
		Name:  "blackPlayer",
		Value: base64.StdEncoding.EncodeToString(byte),
		Path:  "",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "blackDealer",
		Value: "",
		Path:  "",
	})
}
