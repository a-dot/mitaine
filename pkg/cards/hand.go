package cards

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type Hand struct {
	cards []*Card
}

func (h *Hand) AddCard(c *Card) {
	h.cards = append(h.cards, c)
}

func (h *Hand) ContainsRank(rank int) bool {
	for _, c := range h.cards {
		if c.Rank == rank {
			return true
		}
	}

	return false
}

func (h *Hand) Sort() {
	slices.SortFunc(h.cards, func(a, b *Card) int {
		var A, B int

		switch a.Rank {
		case 1:
			A = 100
		case 11:
			A = 200
		default:
			A = a.Rank
		}

		switch b.Rank {
		case 1:
			B = 100
		case 11:
			B = 200
		default:
			B = b.Rank
		}

		return A - B
	})
}

func (h Hand) Draw(screen *ebiten.Image, x, y float64) {
	for _, c := range h.cards {
		fmt.Printf("Draw card %d of %s at %v/%v\n", c.Rank, c.Suit.String(), x, y)
		c.Draw(screen, x, y, 2.0)
	}
}
