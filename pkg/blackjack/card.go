package blackjack

import "fmt"

type Card struct {
	Suit       int
	Number     int
	NumberName string
}

func NewCard(suit int, number int) Card {
	card := Card{Suit: suit, Number: number}
	return card
}

// 0 = ハート, 1 = スペード, 2 = ダイヤ, 3 = クラブ
func CardSuitName(card Card) string {
	suit := ""
	switch card.Suit {
	case 0:
		suit = "ハート"
	case 1:
		suit = "スペード"
	case 2:
		suit = "ダイヤ"
	case 3:
		suit = "クラブ"
	}
	return suit
}

func CardNumberName(card Card) string {
	number := ""
	switch card.Number {
	case 1:
		number = "A"
	case 11:
		number = "J"
	case 12:
		number = "Q"
	case 13:
		number = "K"
	default:
		number = fmt.Sprintf("%v", card.Number)
	}
	return number
}
