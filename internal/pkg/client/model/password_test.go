package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordGetName(t *testing.T) {
	tests := []struct {
		name         string
		password     Password
		expectedName string
	}{
		{
			name: "Check GetName for Password",
			password: Password{
				Resource: "Resource",
				UserName: "UserName",
				Password: "Password",
			},
			expectedName: "Resource",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedName, test.password.GetName())
		})
	}
}
