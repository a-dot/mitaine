package mitaine

import (
	"bytes"
	"fmt"
	"image"
	"log"

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

	Active int // offset of the player who is active

	Dealer deck.Deck
	Play   deck.Deck

	Players []Player
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
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

	for _, p := range g.Players {
		p.Draw(screen)
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
