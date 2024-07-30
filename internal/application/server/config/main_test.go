package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFromEnv(t *testing.T) {
	origStateFun := func() {
		func(args []string) {
			//nolint:reassign //for tests only
			os.Args = args
		}(os.Args)

		func(cl *flag.FlagSet) {
			//nolint:reassign //for tests only
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origStateFun()

	tests := []struct {
		name string
		args []string
		env  map[string]string
		cfg  Config
	}{
		{
			name: "Check config from env",
			args: []string{""},
			env: map[string]string{
				"ADDRESS":   "localhost:8080",
				"LOG_LEVEL": "LOG_LEVEL_FROM_ENV",
			},
			cfg: Config{
				Address:  "localhost:8080",
				LogLevel: "LOG_LEVEL_FROM_ENV",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//nolint:reassign //for tests only
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			for k, v := range test.env {
				t.Setenv(k, v)
			}
			//nolint:reassign //for tests only
			os.Args = test.args

			cfg, err := New()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origStateFun()
		})
	}
}

func TestConfigFromFlags(t *testing.T) {
	origStateFun := func() {
		func(args []string) {
			//nolint:reassign //for tests only
			os.Args = args
		}(os.Args)

		func(cl *flag.FlagSet) {
			//nolint:reassign //for tests only
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origStateFun()

	tests := []struct {
		name string
		args []string
		cfg  Config
	}{
		{
			name: "Check config from flags",
			args: []string{
				"cmd",
				"-address", "localhost:8081",
				"-log-level", "LOG_LEVEL_FROM_FLAG",
			},
			cfg: Config{
				Address:  "localhost:8081",
				LogLevel: "LOG_LEVEL_FROM_FLAG",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//nolint:reassign //for tests only
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			//nolint:reassign //for tests only
			os.Args = test.args

			cfg, err := New()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origStateFun()
		})
	}
}

func TestConfig(t *testing.T) {
	origStateFun := func() {
		func(args []string) {
			//nolint:reassign //for tests only
			os.Args = args
		}(os.Args)

		func(cl *flag.FlagSet) {
			//nolint:reassign //for tests only
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origStateFun()

	tests := []struct {
		name string
		env  map[string]string
		args []string
		cfg  Config
	}{
		{
			name: "Check config from env and args",
			env: map[string]string{
				"ADDRESS":   "localhost:8080",
				"LOG_LEVEL": "LOG_LEVEL_FROM_FLAG",
			},
			args: []string{
				"cmd",
				"-address", "localhost:8081",
				"-log-level", "LOG_LEVEL_FROM_FLAG",
			},
			cfg: Config{
				Address:  "localhost:8081",
				LogLevel: "LOG_LEVEL_FROM_FLAG",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//nolint:reassign //for tests only
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			for k, v := range test.env {
				t.Setenv(k, v)
			}
			//nolint:reassign //for tests only
			os.Args = test.args

			cfg, err := New()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origStateFun()
		})
	}
}
