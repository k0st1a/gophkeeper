package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCardGetName(t *testing.T) {
	tests := []struct {
		name         string
		card         Card
		expectedName string
	}{
		{
			name: "Check GetName for Card",
			card: Card{
				Number:  "Number",
				Expires: "Expires",
				Holder:  "Holder",
			},
			expectedName: "Number",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedName, test.card.GetName())
		})
	}
}
