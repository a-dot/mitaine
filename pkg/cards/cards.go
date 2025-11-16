package cards

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/a-dot/mitaine/pkg/images"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CardWidth  = 48
	CardHeight = 64
)

var (
	imgCards *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.CardsSheet))
	if err != nil {
		log.Fatal(err)
	}
	imgCards = ebiten.NewImageFromImage(img)
}

var Ranks = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

type Suit int

const (
	UnknownSuit Suit = iota
	Heart
	Diamond
	Spade
	Club
)

var Suits = []Suit{Heart, Diamond, Spade, Club}

type Color int

const (
	UnknownColor Color = iota
	Red
	Black
)

type Face int

const (
	UnknownFace Face = iota
	Jack
	Queen
	King
)

type Card struct {
	Rank int
	Suit Suit
}

func (s Suit) String() string {
	switch s {
	case Heart:
		return "Heart"
	case Diamond:
		return "Diamond"
	case Spade:
		return "Spade"
	case Club:
		return "Club"
	}

	return "Unknown"
}

func (c Card) String() string {
	return fmt.Sprintf("%d of %s", c.Rank, c.Suit.String())
}

func (c Card) Color() Color {
	switch c.Suit {
	case Diamond:
		return Red
	case Heart:
		return Red
	case Spade:
		return Black
	case Club:
		return Black
	}

	return UnknownColor
}

func (c Card) IsFace() bool {
	if c.Rank >= 11 && c.Rank <= 13 {
		return true
	}

	return false
}

func (c Card) Face() Face {
	switch c.Rank {
	case 11:
		return Jack
	case 12:
		return Queen
	case 13:
		return King
	}

	return UnknownFace
}

// Draw renders the card on the given screen using the sprite sheet
// The sprite sheet is expected to be organized with suits in rows and ranks in columns
// Rows (suits): 0=Hearts, 1=Diamonds, 2=Spades, 3=Clubs
// Columns (ranks): 0-12 correspond to ranks 1-13
func (c Card) Draw(screen *ebiten.Image, x, y float64, scale float64) {
	if c.Suit == UnknownSuit || c.Rank < 1 || c.Rank > 13 {
		return // Don't draw invalid cards
	}

	// Calculate the position in the sprite sheet
	// Adjust suit index (0-based) for the sprite sheet row
	suitIndex := int(c.Suit) - 1
	// Adjust rank index (0-based) for the sprite sheet column
	rankIndex := c.Rank - 1

	// Calculate the source rectangle in the sprite sheet
	sx := rankIndex * CardWidth
	sy := suitIndex * CardHeight

	// Create the source rectangle
	srcRect := image.Rect(sx, sy, sx+CardWidth, sy+CardHeight)

	// Create drawing options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)

	// Draw the card
	screen.DrawImage(imgCards.SubImage(srcRect).(*ebiten.Image), op)
}
