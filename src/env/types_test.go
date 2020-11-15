package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariables_Validate(t *testing.T) {
	tests := []struct {
		name      string
		variables Variables
		wantErr   bool
	}{
		{
			"empty variables",
			Variables{},
			true,
		},
		{
			"missing rsa private key path",
			Variables{
				StoragePath: "storage",
				LogPath:     "mnemonic",
			},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				require.Error(t, tt.variables.Validate())
			} else {
				require.NoError(t, tt.variables.Validate())
			}
		})
	}
}
