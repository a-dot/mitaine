package main

import (
	"log"

	"github.com/a-dot/mitaine/pkg/mitaine"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(mitaine.ScreenWidth, mitaine.ScreenHeight)
	ebiten.SetWindowTitle("mitaine")

	g := &mitaine.Game{}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
