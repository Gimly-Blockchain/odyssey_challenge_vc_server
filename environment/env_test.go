package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	writeEnv := func(e map[string]string) {
		for k, v := range e {
			require.NoError(t, os.Setenv(k, v))
		}
	}

	tests := []struct {
		name    string
		want    Variables
		env     map[string]string
		wantErr bool
	}{
		{
			"no environment",
			Variables{
				StoragePath: defaultStoragePath,
				LogPath:     defaultLogPath,
			},
			map[string]string{},
			false,
		},
		{
			"okay environment",
			Variables{
				StoragePath: "a",
				LogPath:     "b",
			},
			map[string]string{
				"DSB_STORAGE_PATH": "a",
				"DSB_LOG_PATH":     "b",
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			writeEnv(tt.env)
			defer os.Clearenv()

			me, err := Get()
			if tt.wantErr {
				require.Error(t, err)
			}

			require.Equal(t, tt.want, me)
		})
	}
}

func Test_parseBool(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			"appropriate value returns an appropriate boolean",
			"false",
			false,
		},
		{
			"non-appropriate value returns false",
			"non-appropriate",
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, parseBool(tt.str))
		})
	}
}
