package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoteGetName(t *testing.T) {
	tests := []struct {
		name         string
		note         Note
		expectedName string
	}{
		{
			name: "Check GetName for Note",
			note: Note{
				Name: "Name",
				Body: "Body",
			},
			expectedName: "Name",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedName, test.note.GetName())
		})
	}
}
