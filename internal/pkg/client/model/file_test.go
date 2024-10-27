package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileGetName(t *testing.T) {
	tests := []struct {
		name         string
		file         File
		expectedName string
	}{
		{
			name: "Check GetName for File",
			file: File{
				Name:        "Name",
				Description: "Description",
				Body:        []byte("Body"),
			},
			expectedName: "Name",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedName, test.file.GetName())
		})
	}
}
