package mitaine

import (
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
		game.Players = NewPlayers()

		game.State = NewGame

	case NewGame:
		game.ScoringState = Zero
		game.Turn = 1
		game.Active = 0

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

		game.State = ExchangeCards

	case ExchangeCards:
		//TODO
		game.State = PlayerChoosesCard

	case PlayerChoosesCard:
		activePlayer := &game.Players[game.Active%4]

		activePlayer.Think()

		activePlayer.PlayCard()

		if game.Turn == 2 && game.Active == 3 {
			// TODO that's it
			time.Sleep(1 * time.Second)

			game.State = Scoring
		}

		game.Active++

	case Scoring:
		//TODO

		game.State = NewGame
	}

	return nil
}
