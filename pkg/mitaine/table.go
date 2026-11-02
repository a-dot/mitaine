package mitaine

import (
	"math/rand/v2"

	"github.com/a-dot/mitaine/pkg/cards"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tableCenterX = ScreenWidth / 2
	tableCenterY = ScreenHeight / 2

	maxPileOffset = 14.0 // max pixels a card can sit off-center
	maxPileAngle  = 0.3  // max rotation in radians (~17 degrees)
)

// tableCard is a card lying on the table, slightly off from the pile's center
type tableCard struct {
	card  *cards.Card
	dx    float64
	dy    float64
	angle float64
}

// Table is the pile of cards played in the center of the table
type Table struct {
	cards []tableCard
}

// Place puts a card on top of the pile, randomly offset and rotated
// so the pile looks thrown rather than neatly stacked
func (t *Table) Place(c *cards.Card) {
	t.cards = append(t.cards, tableCard{
		card:  c,
		dx:    (rand.Float64()*2 - 1) * maxPileOffset,
		dy:    (rand.Float64()*2 - 1) * maxPileOffset,
		angle: (rand.Float64()*2 - 1) * maxPileAngle,
	})
}

// Top returns the card on top of the pile, or nil if the table is empty
func (t *Table) Top() *cards.Card {
	if len(t.cards) == 0 {
		return nil
	}

	return t.cards[len(t.cards)-1].card
}

func (t *Table) Len() int {
	return len(t.cards)
}

// Take removes and returns every card on the table, bottom to top
func (t *Table) Take() []*cards.Card {
	ret := make([]*cards.Card, 0, len(t.cards))

	for _, tc := range t.cards {
		ret = append(ret, tc.card)
	}

	t.Empty()

	return ret
}

func (t *Table) Empty() {
	t.cards = nil
}

func (t Table) Draw(screen *ebiten.Image) {
	const scale = 2.0

	for _, tc := range t.cards {
		tc.card.DrawRotated(screen, tableCenterX+tc.dx, tableCenterY+tc.dy, scale, tc.angle)
	}
}
