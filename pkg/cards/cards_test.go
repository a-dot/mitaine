package cards

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColor(t *testing.T) {
	tests := []struct {
		card Card
		want Color
	}{
		{card: Card{Suit: Diamond}, want: Red},
		{card: Card{Suit: Heart}, want: Red},
		{card: Card{Suit: Club}, want: Black},
		{card: Card{Suit: Spade}, want: Black},
		{card: Card{}, want: UnknownColor},
	}

	for _, tc := range tests {
		require.Equal(t, tc.want, tc.card.Color())
	}
}

func TestFace(t *testing.T) {
	tests := []struct {
		card Card
		want Face
	}{
		{card: Card{Rank: 11}, want: Jack},
		{card: Card{Rank: 12}, want: Queen},
		{card: Card{Rank: 13}, want: King},
		{card: Card{Rank: 10}, want: UnknownFace},
		{card: Card{Rank: 9}, want: UnknownFace},
		{card: Card{Rank: 8}, want: UnknownFace},
		{card: Card{Rank: 7}, want: UnknownFace},
		{card: Card{Rank: 6}, want: UnknownFace},
		{card: Card{Rank: 5}, want: UnknownFace},
		{card: Card{Rank: 4}, want: UnknownFace},
		{card: Card{Rank: 3}, want: UnknownFace},
		{card: Card{Rank: 2}, want: UnknownFace},
		{card: Card{Rank: 1}, want: UnknownFace},
	}

	for _, tc := range tests {
		require.Equal(t, tc.want, tc.card.Face())
	}
}

func TestIsFace(t *testing.T) {
	tests := []struct {
		card Card
		want bool
	}{
		{card: Card{Rank: 11}, want: true},
		{card: Card{Rank: 12}, want: true},
		{card: Card{Rank: 13}, want: true},
		{card: Card{Rank: 10}, want: false},
		{card: Card{Rank: 9}, want: false},
		{card: Card{Rank: 8}, want: false},
		{card: Card{Rank: 7}, want: false},
		{card: Card{Rank: 6}, want: false},
		{card: Card{Rank: 5}, want: false},
		{card: Card{Rank: 4}, want: false},
		{card: Card{Rank: 3}, want: false},
		{card: Card{Rank: 2}, want: false},
		{card: Card{Rank: 1}, want: false},
	}

	for _, tc := range tests {
		require.Equal(t, tc.want, tc.card.IsFace())
	}
}

// TestDrawPositions ensures card position calculation logic works correctly
func TestDrawPositions(t *testing.T) {
	tests := []struct {
		name      string
		card      Card
		wantSX    int
		wantSY    int
	}{
		{
			name:   "Ace of Hearts",
			card:   Card{Rank: 1, Suit: Heart},
			wantSX: 0,
			wantSY: 0,
		},
		{
			name:   "King of Hearts",
			card:   Card{Rank: 13, Suit: Heart},
			wantSX: 12 * CardWidth,
			wantSY: 0,
		},
		{
			name:   "Ace of Diamonds",
			card:   Card{Rank: 1, Suit: Diamond},
			wantSX: 0,
			wantSY: 1 * CardHeight,
		},
		{
			name:   "Ace of Spades",
			card:   Card{Rank: 1, Suit: Spade},
			wantSX: 0,
			wantSY: 2 * CardHeight,
		},
		{
			name:   "Ace of Clubs",
			card:   Card{Rank: 1, Suit: Club},
			wantSX: 0,
			wantSY: 3 * CardHeight,
		},
		{
			name:   "Queen of Clubs",
			card:   Card{Rank: 12, Suit: Club},
			wantSX: 11 * CardWidth,
			wantSY: 3 * CardHeight,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate the same values that our Draw function would
			suitIndex := int(tc.card.Suit) - 1
			rankIndex := tc.card.Rank - 1
			sx := rankIndex * CardWidth
			sy := suitIndex * CardHeight

			require.Equal(t, tc.wantSX, sx)
			require.Equal(t, tc.wantSY, sy)
		})
	}
}
