package userid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	tests := []struct {
		name       string
		id         int64
		expectedID int64
	}{
		{
			name:       "Check Set and Get",
			id:         64,
			expectedID: 64,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := Set(context.Background(), test.id)
			id, ok := Get(ctx)
			require.True(t, ok)
			require.Equal(t, test.expectedID, id)
		})
	}
}
