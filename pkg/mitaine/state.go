package mitaine

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/a-dot/mitaine/pkg/deck"
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

		game.State = NewGame

	case NewGame:
		game.ScoringState = Zero
		game.Turn = 1
		game.ActivePlayer = 0

		slog.Debug("empty all the previous decks (discard, play)")
		game.Play.Empty()

		slog.Debug("shuffling a new deck for start of a new game")
		game.Dealer = deck.Standard()
		game.Dealer.Shuffle()

		slog.Debug("switch state to DealCards")
		game.State = DealCards

	case DealCards:
		var n int
		if game.Turn == 1 {
			n = 6
		} else if game.Turn == 2 {
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

		fmt.Printf("waiting for player input, %d cards selected\n", len(p.selectedIndices))

		// wait for the click of the button

		return nil

	case PlayerChoosesCard:
		activePlayer := &game.Players[game.ActivePlayer%4]

		// For CPU players: wait 1 second before allowing them to play
		if activePlayer.typ == CPU {
			// If this is the first frame for this CPU player, start the timer
			if game.WaitStart.IsZero() {
				game.WaitStart = time.Now()
				return nil
			}

			// Wait until 1 second has passed
			if time.Since(game.WaitStart) < time.Second*1 {
				return nil
			}

			// Reset timer for next player
			game.WaitStart = time.Time{}
		}

		// For human players: check if they've selected a card
		if activePlayer.typ == Human && len(activePlayer.selectedIndices) == 0 {
			return nil
		}

		activePlayer.PlayCard()

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
