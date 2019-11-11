package blackjack

import "strconv"

type Player struct {
	Name         string
	Cards        []Card
	Score        int
	ScoreMessage string
	Stanted      bool
	Victory      int
	Defeat       int
}

func NewPlayer(name string) Player {
	var cards []Card
	player := Player{
		Name:    name,
		Cards:   cards,
		Score:   0,
		Stanted: false,
		Victory: 0,
		Defeat:  0,
	}
	return player
}

func RenewalPlayer(player Player, result string) Player {
	var victory int
	var defeat int
	if result == "win" {
		victory = player.Victory + 1
	}
	if result == "lose" {
		defeat = player.Defeat + 1
	}

	var cards []Card
	renewalPlayer := Player{
		Name:    "あなた",
		Cards:   cards,
		Score:   0,
		Stanted: false,
		Victory: victory,
		Defeat:  defeat,
	}
	return renewalPlayer
}

func VictoryDefeatMes(renewalPlayer Player) string {
	return strconv.Itoa(renewalPlayer.Victory) + "勝" + strconv.Itoa(renewalPlayer.Defeat) + "敗です"
}

func DrawCard(player *Player, deck *[]Card) Card {
	card := (*deck)[0]
	*deck = (*deck)[1:]
	card.NumberName = CardNumberName(card)
	player.Cards = append(player.Cards, card)
	return card
}

func PrintDrawCard(player Player, card Card) (string, string) {
	return CardSuitName(card), CardNumberName(card)
}
