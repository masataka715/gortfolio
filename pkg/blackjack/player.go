package blackjack

type Player struct {
	Name         string
	Cards        []Card
	Score        int
	ScoreMessage string
	Stanted      bool
}

func NewPlayer(name string) Player {
	var cards []Card
	player := Player{Name: name, Cards: cards, Score: 0, Stanted: false}
	return player
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
