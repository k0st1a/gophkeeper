package traceid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		name            string
		traceID         string
		expectedTraceID string
	}{
		{
			name:            "Check success add and get from ctx",
			traceID:         "my-trace-id",
			expectedTraceID: "my-trace-id",
		},
		{
			name:            "Check unsuccess get from ctx",
			expectedTraceID: "undefined",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			if test.traceID != "" {
				ctx = Add(ctx, test.traceID)
			}

			traceID := Get(ctx)

			assert.Equal(t, test.expectedTraceID, traceID)
		})
	}
}
