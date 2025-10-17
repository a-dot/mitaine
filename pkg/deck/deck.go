package deck

import (
	"math/rand/v2"

	"github.com/a-dot/mitaine/pkg/cards"
)

type Deck struct {
	Cards []*cards.Card
}

func Standard() Deck {
	d := Deck{
		Cards: make([]*cards.Card, 0, 52),
	}

	for _, s := range cards.Suits {
		for _, r := range cards.Ranks {
			d.Cards = append(d.Cards, &cards.Card{
				Suit: s,
				Rank: r,
			})
		}
	}

	return d
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(a, b int) {
		d.Cards[a], d.Cards[b] = d.Cards[b], d.Cards[a]
	})
}

// Removes and returns a card from a deck
func (d *Deck) Deal() *cards.Card {
	var card *cards.Card

	card, d.Cards = d.Cards[0], d.Cards[1:]

	return card
}

func (d *Deck) Empty() {
	d.Cards = []*cards.Card{}
}
