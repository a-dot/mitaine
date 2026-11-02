package mitaine

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"log/slog"
	"time"

	"github.com/a-dot/mitaine/pkg/cards"
	"github.com/a-dot/mitaine/pkg/deck"
	"github.com/a-dot/mitaine/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768
)

var (
	imgTable *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.TableGreen))
	if err != nil {
		log.Fatal(err)
	}
	imgTable = ebiten.NewImageFromImage(img)
}

type Game struct {
	State        State
	ScoringState ScoringState
	Turn         int

	ActivePlayer int // offset of the player who is active

	Dealer  deck.Deck
	Table   Table     // cards played in the center of the table
	Discard deck.Deck // cards that were swept off the table

	DealerIdx int
	Players   []Player

	// Track previous mouse state to detect clicks (not just presses)
	prevMousePressed bool

	WaitingForPlayerInput bool

	WaitStart    time.Time
	WaitDuration time.Duration // how long the active CPU player waits before playing

	ExchangeButton *Button
	PlayButton     *Button

	// Set when the human clicked the Play button, consumed by the PlayerChoosesCard state
	PlayRequested bool
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Active returns the player whose turn it is
func (g *Game) Active() *Player {
	return &g.Players[g.ActivePlayer%len(g.Players)]
}

// CanPlay reports whether the human is on turn with exactly one card selected
func (g *Game) CanPlay() bool {
	if g.State != PlayerChoosesCard {
		return false
	}

	p := g.Active()

	return p.typ == Human && len(p.selectedIndices) == 1
}

// PlayToTable applies the rules for player p playing card c on the table
func (g *Game) PlayToTable(p *Player, c *cards.Card) {
	prev := g.Table.Top()

	switch {
	case prev == nil:
		// nothing on the table, the card just stays in the middle
		g.Table.Place(c)

	case prev.Rank == c.Rank && g.Table.Len() == 1:
		// sweep: the played card goes face up in the player's points pile,
		// the card it matched is scrapped
		slog.Info("sweep: played card goes to the points pile, matched card is scrapped", "player", p.name, "played", c, "matched", prev)

		g.Table.Take()
		p.points = append(p.points, c)

	case prev.Rank == c.Rank:
		// same card on a pile of several: everything is discarded
		slog.Info("clears the table", "player", p.name, "played", c, "matched", prev, "discarded", g.Table.Len()+1)

		g.Discard.Cards = append(g.Discard.Cards, g.Table.Take()...)
		g.Discard.Cards = append(g.Discard.Cards, c)

	default:
		// TODO handle the other conditions (e.g. higher/lower cards, face cards)
		g.Table.Place(c)
	}
}

// OnExchangeClick handles the Exchange button click event
func (g *Game) OnExchangeClick() {
	// Choose which cards will be exchanged
	exchanges := make([][]*cards.Card, 4)

	for i, p := range g.Players {
		if len(p.selectedIndices) > 0 {
			exchanges[i] = []*cards.Card{p.hand.Pos(p.selectedIndices[0]), p.hand.Pos(p.selectedIndices[1])}
		} else {
			if i == g.DealerIdx || i == g.DealerIdx+1 {
				// lower cards
				exchanges[i] = []*cards.Card{}

				for j := 0; j < p.hand.CardCount() && len(exchanges[i]) < 2; j++ {
					c := p.hand.Pos(j)

					if c == nil {
						// TODO should never happen
					}

					if j < p.hand.CardCount()-1 && c.Rank == p.hand.Pos(j+1).Rank {
						// we don't want to break a hand
						continue
					}

					exchanges[i] = append(exchanges[i], p.hand.Pos(j))
				}
			} else {
				// higher cards
				exchanges[i] = []*cards.Card{}

				for j := p.hand.CardCount() - 1; j >= 0 && len(exchanges[i]) < 2; j-- {
					c := p.hand.Pos(j)

					if c == nil {
						// TODO should never happen
					}

					if j > 0 && c.Rank == p.hand.Pos(j-1).Rank {
						// we don't want to break a hand
						continue
					}

					exchanges[i] = append(exchanges[i], p.hand.Pos(j))
				}
			}
		}

		slog.Info("exchanges cards", "player", p.name, "seat", i, "cards", exchanges[i])
	}

	// Player 0 exchanges with Player 2
	g.Players[2].hand.AddCard(exchanges[0][0])
	g.Players[2].hand.AddCard(exchanges[0][1])
	g.Players[0].hand.RemoveCard(exchanges[0][1])
	g.Players[0].hand.RemoveCard(exchanges[0][0])

	// Player 2 exchanges with Player 0
	g.Players[0].hand.AddCard(exchanges[2][0])
	g.Players[0].hand.AddCard(exchanges[2][1])
	g.Players[2].hand.RemoveCard(exchanges[2][1])
	g.Players[2].hand.RemoveCard(exchanges[2][0])

	// Player 1 exchanges with Player 3
	g.Players[3].hand.AddCard(exchanges[1][0])
	g.Players[3].hand.AddCard(exchanges[1][1])
	g.Players[1].hand.RemoveCard(exchanges[1][1])
	g.Players[1].hand.RemoveCard(exchanges[1][0])

	// Player 3 exchanges with Player 1
	g.Players[1].hand.AddCard(exchanges[3][0])
	g.Players[1].hand.AddCard(exchanges[3][1])
	g.Players[3].hand.RemoveCard(exchanges[3][1])
	g.Players[3].hand.RemoveCard(exchanges[3][0])

	for i := range g.Players {
		g.Players[i].hand.Sort()

		// unselect cards
		g.Players[i].selectedIndices = []int{}
	}

	g.ActivePlayer = g.DealerIdx

	g.State = PlayerChoosesCard
}

func (g *Game) Update() error {
	// Handle mouse clicks for card selection (detect click = press then release)
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Detect click: was pressed and now released
	if g.prevMousePressed && !mousePressed {
		mouseX, mouseY := ebiten.CursorPosition()

		// Check if Exchange button was clicked (only in ExchangeCards state)
		if g.State == ExchangeCards {
			humanPlayer := HumanPlayer(g.Players)
			if len(humanPlayer.selectedIndices) >= 2 && g.ExchangeButton.Contains(mouseX, mouseY) {
				g.OnExchangeClick()
			}
		}

		// Check if Play button was clicked (only when the human has selected a card)
		if g.CanPlay() && g.PlayButton.Contains(mouseX, mouseY) {
			g.PlayRequested = true
		}

		// Check if any human player clicked on their cards
		for i := range g.Players {
			if g.Players[i].HandleClick(g.State, mouseX, mouseY) {
				break // Only handle one click at a time
			}
		}
	}

	g.prevMousePressed = mousePressed

	return g.State.Update(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the message.
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)

	// Draw the table background stretched to fit the window
	op := &ebiten.DrawImageOptions{}

	// Get original image dimensions
	imgWidth, imgHeight := imgTable.Bounds().Dx(), imgTable.Bounds().Dy()

	// Calculate scale factors to fill the screen
	scaleX := float64(ScreenWidth) / float64(imgWidth)
	scaleY := float64(ScreenHeight) / float64(imgHeight)

	// Apply scaling
	op.GeoM.Scale(scaleX, scaleY)

	// Draw centered
	screen.DrawImage(imgTable, op)

	g.Table.Draw(screen)

	for _, p := range g.Players {
		p.Draw(screen)
	}

	// Draw Play button when the human has selected a card to play
	if g.CanPlay() {
		g.PlayButton.Draw(screen)
	}

	// Draw Exchange button if in ExchangeCards state and 2+ cards selected
	if g.State == ExchangeCards {
		humanPlayer := HumanPlayer(g.Players)
		if len(humanPlayer.selectedIndices) >= 2 {
			g.ExchangeButton.Draw(screen)
		}
	}

	// // Draw some cards as examples
	// if len(g.Dealer.Cards) > 0 {
	// 	// Draw first card
	// 	g.Dealer.Cards[0].Draw(screen, imgCards, 100, 100, 2.0)

	// 	// Draw second card if available
	// 	if len(g.Dealer.Cards) > 1 {
	// 		g.Dealer.Cards[1].Draw(screen, imgCards, 200, 100, 2.0)
	// 	}

	// 	// Draw third card if available
	// 	if len(g.Dealer.Cards) > 2 {
	// 		g.Dealer.Cards[2].Draw(screen, imgCards, 300, 100, 2.0)
	// 	}
	// }

	// // Create and draw a specific card as an example
	// aceHearts := cards.Card{Rank: 1, Suit: cards.Heart}
	// kingClubs := cards.Card{Rank: 13, Suit: cards.Club}

	// aceHearts.Draw(screen, imgCards, 100, 300, 2.0)
	// kingClubs.Draw(screen, imgCards, 200, 300, 2.0)
}
