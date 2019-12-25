// trello-tui
// A terminal ui for trello
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-tui/pkg/app"
	"github.com/giannimassi/trello-tui/pkg/gui"
	"github.com/giannimassi/trello-tui/pkg/state"
	"github.com/giannimassi/trello-tui/pkg/trello"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	// TrelloUser is a the environment variable for storing the trello user
	TrelloUser = "TRELLO_USER"
	// TrelloKey is a the environment variable for storing the trello key
	TrelloKey = "TRELLO_KEY"
	// TrelloToken is a the environment variable for storing the trello token
	TrelloToken = "TRELLO_TOKEN"

	minRefreshInterval     = time.Second * 1
	defaultRefreshInterval = time.Second * 10
)

// setup parses the configuration from flags and envronment and setups the global logger
func setup() (app.Config, func()) {
	boardName := flag.String("board", "", "board name")
	refresh := flag.Duration("refresh", defaultRefreshInterval, fmt.Sprintf("refresh interval (min=%v)", minRefreshInterval))
	logFlag := flag.Bool("log", false, "Log to file")
	v := flag.Bool("vv", false, "Increase verbosity level")
	flag.Parse()

	cleanup := func() {}
	if *logFlag {
		f, err := os.OpenFile("trello-tui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err).Msg("Unexpected error while opening file for logging. Stopping application")
		}
		_, _ = f.Write([]byte("\n"))
		cleanup = func() { f.Close() }
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: f})
		if !*v {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	} else {
		log.Logger = log.Output(ioutil.Discard)
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	if *refresh < minRefreshInterval {
		log.Warn().Msg("Minimum value for refresh interval is 10 s")
		*refresh = minRefreshInterval
	}

	return app.Config{
		State: state.Config{
			Trello: trello.Config{
				User:    os.Getenv(TrelloUser),
				Key:     os.Getenv(TrelloKey),
				Token:   os.Getenv(TrelloToken),
				Timeout: time.Second * 10,
			},
			SelectedBoard:        *boardName,
			BoardRefreshInterval: *refresh,
		},

		Gui: gui.Config{
			Dev: *v,
		},
	}, cleanup
}

func main() {
	var (
		cfg, cleanup = setup()
		a            = app.NewApp(&cfg)
	)
	defer func() {
		log.Info().Msg("Quitting trello-tui")
		cleanup()
	}()

	log.Info().Str("version", version).Str("commit", commit).Str("build", date).Msg("Starting trello-tui")
	if err := a.Init(); err != nil {
		log.Error().Err(err).Msg("Unexpected error while initializing. Stopping application")
		return
	}

	if err := a.Run(); err != nil {
		log.Error().Err(err).Msg("Unexpected error while running. Stopping application")
	}
	a.Close()
}
