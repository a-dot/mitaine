package deck

import (
	"reflect"
	"testing"

	"github.com/a-dot/mitaine/pkg/cards"
	"github.com/stretchr/testify/require"
)

func TestStandard(t *testing.T) {
	_ = Standard()

	//TODO
}

func TestShuffle(t *testing.T) {
	dA := Standard()

	dB := Deck{
		Cards: make([]*cards.Card, len(dA.Cards)),
	}
	copy(dB.Cards, dA.Cards)

	dA.Shuffle()

	require.Equal(t, 52, len(dA.Cards))
	require.Equal(t, 52, len(dB.Cards))
	require.False(t, reflect.DeepEqual(dA, dB))
}
