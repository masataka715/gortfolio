package blackjack

import (
	"math/rand"
	"time"
)

// スート4種 × 1〜13
func NewDeck() []Card {
	var deck []Card
	for i := 0; i < 4; i++ {
		for j := 1; j < 14; j++ {
			deck = append(deck, NewCard(i, j))
		}
	}
	return deck
}

func ShuffleDeck(deck []Card) {
	// 乱数のシード値に現在時刻を設定
	rand.Seed(time.Now().UnixNano())

	n := len(deck)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}
