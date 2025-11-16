package mitaine

import (
	"fmt"

	"github.com/a-dot/mitaine/pkg/cards"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerType int

const (
	CPU PlayerType = iota
	Human
)

type Player struct {
	typ PlayerType

	name string

	x        float64
	y        float64
	xx       float64
	position cards.PlayerPosition

	hand            cards.Hand
	selectedIndices []int // Indices of the selected cards (empty slice if none selected)
}

func NewPlayers(names ...string) []Player {
	ret := make([]Player, 0, len(names))

	nPlayers := 4

	for _, name := range names {
		x, y, xx, pos := playerCoords(nPlayers)

		p := Player{
			name:            name,
			typ:             Human,
			x:               x,
			y:               y,
			xx:              xx,
			position:        pos,
			selectedIndices: []int{}, // No cards selected initially
		}

		nPlayers--

		ret = append(ret, p)
	}

	for i := 0; i < 4-len(names); i++ {
		x, y, xx, pos := playerCoords(nPlayers)

		p := Player{
			name:            fmt.Sprintf("CPU%d", i+1),
			typ:             CPU,
			x:               x,
			y:               y,
			xx:              xx,
			position:        pos,
			selectedIndices: []int{},
		}

		nPlayers--

		ret = append(ret, p)
	}

	return ret
}

func playerCoords(i int) (x, y, xx float64, position cards.PlayerPosition) {
	const (
		cardWidth     = 48 * 2.0 // 96 pixels when scaled 2x
		cardHeight    = 64 * 2.0 // 128 pixels when scaled 2x
		margin        = 10.0     // margin from screen edges
		overlapFactor = 0.3      // Cards overlap by 70%, showing 30%
		maxCards      = 7        // Maximum cards a player can have
	)

	// Calculate the total width needed for a full hand
	handWidth := cardWidth + float64(maxCards-1)*(cardWidth*overlapFactor)

	switch i {
	case 4: // Bottom player
		x = (1024 - handWidth) / 2 // Center horizontally
		y = 768 - cardHeight - margin
		xx = 600
		position = cards.PositionBottom
	case 3: // Left player
		x = margin
		y = 300
		xx = 200
		position = cards.PositionLeft
	case 2: // Top player
		x = (1024 - handWidth) / 2 // Center horizontally
		y = margin
		xx = 600
		position = cards.PositionTop
	case 1: // Right player
		x = 1024 - handWidth - margin // Position so hand ends at right edge
		y = 300
		xx = 800
		position = cards.PositionRight
	default:
		// TODO die
		panic("invalid")
	}

	return
}

func (p *Player) PlayCard() {
	switch p.typ {
	case CPU:
		p.CPUPlayCard()
	case Human:
		p.HumanPlayCard()
	}
}

// Returns first Human player in the slice
func HumanPlayer(players []Player) *Player {
	for _, p := range players {
		if p.typ == Human {
			return &p
		}
	}

	return nil
}

// HumanPlayCard waits for the human to select cards by clicking
// The card selection is handled by HandleClick, this method just waits
func (p *Player) HumanPlayCard() {
	// Human players select cards via mouse clicks (HandleClick)
	// The game state should wait until selectedIndices has at least one card
	// This is a placeholder - actual implementation depends on game flow
	if len(p.selectedIndices) > 0 {
		fmt.Printf("Human player selected %d card(s) at indices %v\n", len(p.selectedIndices), p.selectedIndices)
		// 	// TODO: Actually play the selected cards and reset selection
		// 	p.selectedIndices = []int{}
	}
}

// CPU chooses a card and plays it
func (p *Player) CPUPlayCard() {
	fmt.Println("what to do, what to do?")
}

func (p Player) Draw(screen *ebiten.Image) {
	p.hand.Draw(screen, p.x, p.y, p.position, p.selectedIndices)
}

// HandleClick checks if a card was clicked and updates the selection
// Returns true if a card was clicked
func (p *Player) HandleClick(state State, mouseX, mouseY int) bool {
	if p.typ != Human {
		return false
	}

	const (
		cardWidth     = 48 * 2.0 // 96 pixels when scaled 2x
		cardHeight    = 64 * 2.0 // 128 pixels when scaled 2x
		overlapFactor = 0.3
		horizOffset   = cardWidth * overlapFactor
	)

	numCards := p.hand.CardCount()
	mx, my := float64(mouseX), float64(mouseY)

	// Check cards in reverse order (rightmost/topmost first) for proper overlap handling
	for i := numCards - 1; i >= 0; i-- {
		cardX := p.x + float64(i)*horizOffset
		cardY := p.y

		// Check if mouse is within card bounds
		if mx >= cardX && mx < cardX+cardWidth &&
			my >= cardY && my < cardY+cardHeight {
			// Toggle selection: if card is already selected, remove it; otherwise add it
			found := false
			for idx, selected := range p.selectedIndices {
				if selected == i {
					// Remove from selection
					p.selectedIndices = append(p.selectedIndices[:idx], p.selectedIndices[idx+1:]...)
					found = true
					break
				}
			}
			if !found {
				if state == ExchangeCards && len(p.selectedIndices) == 2 {
					// can't have more than 2 cards selected when Exchanging cards
					return false
				}

				// Add to selection
				p.selectedIndices = append(p.selectedIndices, i)
			}
			return true
		}
	}

	return false
}
