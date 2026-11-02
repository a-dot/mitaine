package mitaine

import (
	"fmt"
	"log/slog"
	"math/rand/v2"

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

	points []*cards.Card // pile of won cards, shown face up next to the player
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

func (p *Player) PlayCard(g *Game) {
	switch p.typ {
	case CPU:
		p.CPUPlayCard(g)
	case Human:
		p.HumanPlayCard(g)
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

// HumanPlayCard plays the selected card on the table.
// The card selection is handled by HandleClick and the Play button; the game
// only calls this once the human asked to play.
func (p *Player) HumanPlayCard(g *Game) {
	if len(p.selectedIndices) != 1 {
		return
	}

	c := p.hand.Pos(p.selectedIndices[0])
	if c == nil {
		return
	}

	slog.Info("plays", "player", p.name, "card", c)

	p.hand.RemoveCard(c)
	p.selectedIndices = []int{}

	g.PlayToTable(p, c)
}

// CPU chooses a card and plays it
func (p *Player) CPUPlayCard(g *Game) {
	n := p.hand.CardCount()
	if n == 0 {
		return
	}

	c := p.cpuChooseCard(g.Table.Top())

	slog.Info("plays", "player", p.name, "card", c)

	p.hand.RemoveCard(c)

	g.PlayToTable(p, c)
}

// cpuChooseCard picks the card the CPU plays given the card on top of the pile
// (nil if the table is empty).
//
// It plays a card of the same rank as the top one, as long as it doesn't hold
// several of that rank: with a pair or more in hand it keeps them instead.
// Otherwise it plays a random card, avoiding that rank when it can.
//
// TODO look at more than the top of the table (what's left in its hand, ...)
func (p *Player) cpuChooseCard(top *cards.Card) *cards.Card {
	n := p.hand.CardCount()
	if n == 0 {
		return nil
	}

	if top == nil {
		return p.hand.Pos(rand.IntN(n))
	}

	// a lone card matching the top one is played right away
	if p.hand.CountRank(top.Rank) == 1 {
		for i := 0; i < n; i++ {
			if c := p.hand.Pos(i); c.Rank == top.Rank {
				return c
			}
		}
	}

	// otherwise play any card, except the ones of the top rank that we hold
	// several of (unless the hand is made only of those)
	candidates := make([]*cards.Card, 0, n)
	for i := 0; i < n; i++ {
		if c := p.hand.Pos(i); c.Rank != top.Rank {
			candidates = append(candidates, c)
		}
	}

	if len(candidates) == 0 {
		return p.hand.Pos(rand.IntN(n))
	}

	return candidates[rand.IntN(len(candidates))]
}

func (p Player) Draw(screen *ebiten.Image) {
	p.hand.Draw(screen, p.x, p.y, p.position, p.selectedIndices)
	p.drawPoints(screen)
}

// drawPoints draws the points pile face up near the player's hand:
// just above it for the left and right players, off to the side for the
// top (left of the hand) and bottom (right of the hand) players
func (p Player) drawPoints(screen *ebiten.Image) {
	const (
		scale         = 2.0
		cardWidth     = cards.CardWidth * scale
		cardHeight    = cards.CardHeight * scale
		overlapFactor = 0.3 // must match the hand layout
		maxCards      = 7
		handWidth     = cardWidth + (maxCards-1)*cardWidth*overlapFactor
		gap           = 20.0 // space between hand and points pile
		stackOffset   = 3.0  // each card in the pile sits slightly offset from the previous
	)

	var x, y float64

	switch p.position {
	case cards.PositionLeft, cards.PositionRight:
		x = p.x
		y = p.y - cardHeight - gap
	case cards.PositionTop:
		x = p.x - cardWidth - gap
		y = p.y
	default: // bottom
		x = p.x + handWidth + gap
		y = p.y
	}

	for i, c := range p.points {
		offset := float64(i) * stackOffset
		c.Draw(screen, x+offset, y-offset, scale)
	}
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

				if state == PlayerChoosesCard {
					// only one card can be played, so picking another replaces the selection
					p.selectedIndices = []int{i}
					return true
				}

				// Add to selection
				p.selectedIndices = append(p.selectedIndices, i)
			}
			return true
		}
	}

	return false
}
