package mitaine

import (
	"image/color"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/a-dot/mitaine/pkg/deck"
)

// Bounds of the random delay before a CPU player plays its card
const (
	cpuMinWait = 1500 * time.Millisecond
	cpuMaxWait = 4 * time.Second
)

type State int

const (
	Initial State = iota
	NewGame
	DealCards
	ExchangeCards
	PlayerChoosesCard
	Scoring
)

type ScoringState int

const (
	Zero ScoringState = iota
	Mitaine
	Gant
	Chausson
)

func (s State) Update(game *Game) error {
	switch s {
	case Initial:
		game.Players = NewPlayers("Alice")

		// Initialize Exchange button positioned above the bottom player's cards
		const (
			buttonWidth  = 100.0
			buttonHeight = 30.0
			cardHeight   = 64 * 2.0 // 128 pixels when scaled 2x
			margin       = 10.0
			buttonMargin = 10.0 // Space between button and cards
			screenWidth  = 1024
			screenHeight = 768
		)
		buttonX := (screenWidth - buttonWidth) / 2 // Center horizontally
		buttonY := screenHeight - cardHeight - margin - buttonHeight - buttonMargin
		game.ExchangeButton = NewButton(buttonX, buttonY, buttonWidth, buttonHeight, "Exchange")
		game.ExchangeButton.BgColor = color.RGBA{46, 139, 87, 255}

		// Play button sits in the same spot as the Exchange button, they are never shown together
		game.PlayButton = NewButton(buttonX, buttonY, buttonWidth, buttonHeight, "Play")
		game.PlayButton.BgColor = color.RGBA{178, 34, 34, 255}

		game.State = NewGame

	case NewGame:
		game.ScoringState = Zero
		game.Turn = 1
		game.ActivePlayer = 0

		slog.Debug("empty all the previous decks (discard, table) and points piles")
		game.Table.Empty()
		game.Discard.Empty()
		game.PlayRequested = false

		for i := range game.Players {
			game.Players[i].points = nil
		}

		slog.Debug("shuffling a new deck for start of a new game")
		game.Dealer = deck.Standard()
		game.Dealer.Shuffle()

		slog.Debug("switch state to DealCards")
		game.State = DealCards

	case DealCards:
		var n int

		switch game.Turn {
		case 1:
			n = 6
		default:
			n = 7
		}

		for range n {
			for i := range game.Players {
				game.Players[i].hand.AddCard(game.Dealer.Deal())
			}
		}

		for i := range game.Players {
			game.Players[i].hand.Sort()
		}

		game.WaitingForPlayerInput = true
		game.State = ExchangeCards

	case ExchangeCards:
		if !game.WaitingForPlayerInput {
			return nil
		}

		p := HumanPlayer(game.Players)

		slog.Debug("waiting for player input", "selected", len(p.selectedIndices))

		// wait for the click of the button

		return nil

	case PlayerChoosesCard:
		activePlayer := game.Active()

		// For CPU players: wait a random time before allowing them to play
		if activePlayer.typ == CPU {
			// If this is the first frame for this CPU player, start the timer
			if game.WaitStart.IsZero() {
				game.WaitStart = time.Now()
				game.WaitDuration = cpuMinWait + time.Duration(rand.Float64()*float64(cpuMaxWait-cpuMinWait))
				slog.Info("cpu waiting before playing", "player", activePlayer.name, "duration", game.WaitDuration)
				return nil
			}

			// Wait until the random delay has passed
			if time.Since(game.WaitStart) < game.WaitDuration {
				return nil
			}

			// Reset timer for next player
			game.WaitStart = time.Time{}
		}

		// For human players: wait until they select a card and click Play
		if activePlayer.typ == Human && !game.PlayRequested {
			return nil
		}

		game.PlayRequested = false

		activePlayer.PlayCard(game)

		if game.Turn == 2 && ((game.ActivePlayer % 4) == 3) {
			// TODO that's it
			time.Sleep(1 * time.Second)

			game.State = Scoring
		}

		game.ActivePlayer++

	case Scoring:
		game.DealerIdx++
		if game.DealerIdx == 4 {
			game.DealerIdx = 0
		}

		game.State = NewGame
	}

	return nil
}
