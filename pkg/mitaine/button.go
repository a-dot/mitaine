package mitaine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	x      float64
	y      float64
	width  float64
	height float64
	text   string
}

func NewButton(x, y, width, height float64, label string) *Button {
	return &Button{
		x:      x,
		y:      y,
		width:  width,
		height: height,
		text:   label,
	}
}

// Contains checks if a point is within the button's bounds
func (b *Button) Contains(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= b.x && fx < b.x+b.width &&
		fy >= b.y && fy < b.y+b.height
}

// Draw renders the button
func (b *Button) Draw(screen *ebiten.Image) {
	// Draw button background
	vector.DrawFilledRect(screen,
		float32(b.x), float32(b.y),
		float32(b.width), float32(b.height),
		color.RGBA{70, 130, 180, 255}, // Steel blue
		false)

	// Draw button border
	vector.StrokeRect(screen,
		float32(b.x), float32(b.y),
		float32(b.width), float32(b.height),
		2.0,
		color.RGBA{255, 255, 255, 255}, // White border
		false)

	// Draw button text (centered)
	// Using ebitenutil for simple text rendering
	textX := int(b.x + 15)
	textY := int(b.y + 10)
	ebitenutil.DebugPrintAt(screen, b.text, textX, textY)
}
