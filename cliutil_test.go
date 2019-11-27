package cliutil_test

import (
	"strings"
	"testing"
	"time"

	"github.com/alileza/cliutil"
	"github.com/urfave/cli"
)

func TestGenerateFlags(t *testing.T) {
	var config struct {
		NodeID           string  `default:"ulala"`
		LogLevel         string  `usage:"just to see logs"`
		ListenAddress    int64   `flag:"flag-name"`
		AdvertiseAddress int     `env:"env-name"`
		MigrationPath    float64 `hidden:"true"`
		DatabaseDSN      int
		BoolTest         bool
		Test             int64
		SkipMigration    time.Duration
	}

	t.Run("test input value", func(t *testing.T) {
		_, err := cliutil.GenerateFlags(config)
		if !strings.Contains(err.Error(), "Config has to be a pointer of struct") {
			t.Error("Expecting input to be pointer of struct, got something else")
		}
	})

	t.Run("test input value", func(t *testing.T) {
		var a map[string]interface{}
		_, err := cliutil.GenerateFlags(&a)
		if !strings.Contains(err.Error(), "Config has to be a pointer of struct") {
			t.Error("Expecting input to be pointer of struct, got something else")
		}
	})

	app := cli.NewApp()

	flags, err := cliutil.GenerateFlags(&config)
	if err != nil {
		t.Fatalf("cliutil: %v", err)
	}
	app.Flags = flags

	app.Action = func(ctx *cli.Context) error {
		return cli.ShowAppHelp(ctx)
	}
	err = app.Run([]string{"cliutil"})
	if err != nil {
		t.Fatalf("cliutil: %v", err)
	}
}
