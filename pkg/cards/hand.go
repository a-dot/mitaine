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

func (h *Hand) CardCount() int {
	return len(h.cards)
}

func (h *Hand) RemoveCard(rm *Card) {
	idx := -1

	for i := range h.cards {
		if h.cards[i] == rm {
			idx = i
			break
		}
	}

	if idx > -1 {
		h.cards = slices.Delete(h.cards, idx, idx+1)
	}
}

func (h *Hand) ContainsRank(rank int) bool {
	for _, c := range h.cards {
		if c.Rank == rank {
			return true
		}
	}

	return false
}

func (h *Hand) Display() {
	fmt.Printf("%s\n", h.cards)
}

func (h *Hand) Pos(i int) *Card {
	if i >= h.CardCount() {
		return nil
	}

	if i < 0 {
		return nil
	}

	return h.cards[i]
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

// PlayerPosition represents where a player is seated at the table
type PlayerPosition int

const (
	PositionBottom PlayerPosition = iota
	PositionLeft
	PositionTop
	PositionRight
)

// Draw renders all cards in the hand with proper spacing based on player position
// All players now display cards horizontally
// selectedIndices contains indices of selected cards, which will be offset upward
func (h Hand) Draw(screen *ebiten.Image, x, y float64, position PlayerPosition, selectedIndices []int) {
	const (
		scale          = 2.0
		cardWidth      = CardWidth * scale  // 48 * 2 = 96
		cardHeight     = CardHeight * scale // 64 * 2 = 128
		overlapFactor  = 0.3                // Cards overlap by 70%, showing 30%
		horizOffset    = cardWidth * overlapFactor
		selectedOffset = 20.0 // How many pixels to move selected card up
	)

	// Create a set for quick lookup of selected indices
	selectedSet := make(map[int]bool)
	for _, idx := range selectedIndices {
		selectedSet[idx] = true
	}

	for i, c := range h.cards {
		// All positions use horizontal layout
		cardX := x + float64(i)*horizOffset
		cardY := y

		// If this card is selected, move it up
		if selectedSet[i] {
			cardY -= selectedOffset
		}

		c.Draw(screen, cardX, cardY, scale)
	}
}
