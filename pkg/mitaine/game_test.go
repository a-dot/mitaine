package mitaine

import (
	"testing"

	"github.com/a-dot/mitaine/pkg/cards"
	"github.com/stretchr/testify/require"
)

func TestPlayToTable(t *testing.T) {
	var (
		sevenHeart = func() *cards.Card { return &cards.Card{Rank: 7, Suit: cards.Heart} }
		sevenSpade = func() *cards.Card { return &cards.Card{Rank: 7, Suit: cards.Spade} }
		kingClub   = func() *cards.Card { return &cards.Card{Rank: 13, Suit: cards.Club} }
	)

	tests := []struct {
		name        string
		onTable     []*cards.Card
		played      *cards.Card
		wantTable   int
		wantPoints  int
		wantDiscard int
	}{
		{
			name:      "empty table",
			played:    sevenHeart(),
			wantTable: 1,
		},
		{
			name:       "sweep: same rank on a single card",
			onTable:    []*cards.Card{sevenHeart()},
			played:     sevenSpade(),
			wantTable:  0,
			wantPoints: 1,
		},
		{
			name:        "same rank on several cards discards everything",
			onTable:     []*cards.Card{kingClub(), sevenHeart()},
			played:      sevenSpade(),
			wantTable:   0,
			wantDiscard: 3,
		},
		{
			name:      "different rank goes on top",
			onTable:   []*cards.Card{sevenHeart()},
			played:    kingClub(),
			wantTable: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{}
			p := &Player{}

			for _, c := range tt.onTable {
				g.Table.Place(c)
			}

			g.PlayToTable(p, tt.played)

			require.Equal(t, tt.wantTable, g.Table.Len())
			require.Len(t, p.points, tt.wantPoints)
			require.Len(t, g.Discard.Cards, tt.wantDiscard)

			if tt.wantTable > 0 {
				require.Same(t, tt.played, g.Table.Top())
			}

			if tt.wantPoints > 0 {
				require.Same(t, tt.played, p.points[0])
			}
		})
	}
}

func TestTablePlaceIsOffCenter(t *testing.T) {
	var table Table

	for range 20 {
		table.Place(&cards.Card{Rank: 2, Suit: cards.Heart})
	}

	var offset bool

	for _, tc := range table.cards {
		require.LessOrEqual(t, tc.dx, maxPileOffset)
		require.GreaterOrEqual(t, tc.dx, -maxPileOffset)
		require.LessOrEqual(t, tc.dy, maxPileOffset)
		require.GreaterOrEqual(t, tc.dy, -maxPileOffset)
		require.LessOrEqual(t, tc.angle, maxPileAngle)
		require.GreaterOrEqual(t, tc.angle, -maxPileAngle)

		if tc.dx != 0 || tc.dy != 0 || tc.angle != 0 {
			offset = true
		}
	}

	require.True(t, offset)
}

func TestCPUChooseCard(t *testing.T) {
	var (
		twoHeart = &cards.Card{Rank: 2, Suit: cards.Heart}
		twoSpade = &cards.Card{Rank: 2, Suit: cards.Spade}
		fiveClub = &cards.Card{Rank: 5, Suit: cards.Club}
		kingClub = &cards.Card{Rank: 13, Suit: cards.Club}
		topTwo   = &cards.Card{Rank: 2, Suit: cards.Diamond}
		handOf   = func(cs ...*cards.Card) cards.Hand {
			var h cards.Hand
			for _, c := range cs {
				h.AddCard(c)
			}
			return h
		}
	)

	tests := []struct {
		name string
		hand cards.Hand
		top  *cards.Card
		want []*cards.Card // any of these is acceptable
	}{
		{
			name: "plays the matching card when it has only one",
			hand: handOf(fiveClub, twoHeart, kingClub),
			top:  topTwo,
			want: []*cards.Card{twoHeart},
		},
		{
			name: "keeps a pair instead of matching",
			hand: handOf(twoHeart, fiveClub, twoSpade),
			top:  topTwo,
			want: []*cards.Card{fiveClub},
		},
		{
			name: "no matching card plays any card",
			hand: handOf(fiveClub, kingClub),
			top:  topTwo,
			want: []*cards.Card{fiveClub, kingClub},
		},
		{
			name: "only a pair left has to play one of them",
			hand: handOf(twoHeart, twoSpade),
			top:  topTwo,
			want: []*cards.Card{twoHeart, twoSpade},
		},
		{
			name: "empty table plays any card",
			hand: handOf(fiveClub, twoHeart),
			want: []*cards.Card{fiveClub, twoHeart},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{hand: tt.hand}

			// random fallbacks, so try many times
			for range 50 {
				require.Contains(t, tt.want, p.cpuChooseCard(tt.top))
			}
		})
	}
}
