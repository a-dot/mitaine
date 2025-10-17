package mitaine

import (
	"fmt"
	"time"

	"github.com/a-dot/mitaine/pkg/cards"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerType int

const (
	CPU PlayerType = iota
	Human
)

const THINK_TIME_MS = 1000

type Player struct {
	typ PlayerType

	name string

	x  float64
	y  float64
	xx float64

	hand cards.Hand
}

func NewPlayers(names ...string) []Player {
	ret := make([]Player, 0, len(names))

	nPlayers := 4

	for _, name := range names {
		x, y, xx := playerCoords(nPlayers)

		p := Player{
			name: name,
			typ:  Human,
			x:    x,
			y:    y,
			xx:   xx,
		}

		nPlayers--

		ret = append(ret, p)
	}

	for i := 0; i < 4-len(names); i++ {
		x, y, xx := playerCoords(nPlayers)

		p := Player{
			name: fmt.Sprintf("CPU%d", i+1),
			typ:  CPU,
			x:    x,
			y:    y,
			xx:   xx,
		}

		nPlayers--

		ret = append(ret, p)
	}

	return ret
}

func playerCoords(i int) (x, y, xx float64) {
	switch i {
	case 4:
		x = 400
		y = 768
		xx = 600
	case 3:
		x = 0
		y = 300
		xx = 200
	case 2:
		x = 400
		y = 0
		xx = 600
	case 1:
		x = 1024
		y = 300
		xx = 800
	default:
		// TODO die
		panic("invalid")
	}

	return
}

func (p Player) Think() {
	if p.typ == CPU {
		time.Sleep(THINK_TIME_MS * time.Millisecond)
	}
}

func (p *Player) PlayCard() {
	switch p.typ {
	case CPU:
		p.CPUPlayCard()
	case Human:
		//TODO
	}
}

// CPU chooses a card and plays it
func (p *Player) CPUPlayCard() {
	fmt.Println("what to do, what to do?")
}

func (p Player) Draw(screen *ebiten.Image) {
	p.hand.Draw(screen, p.x, p.y)
}
