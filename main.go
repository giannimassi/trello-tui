package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-cli/pkg/gui"
	"github.com/giannimassi/trello-cli/pkg/trello"
)

const (
	TrelloUser  = "TRELLO_USER"
	TrelloKey   = "TRELLO_KEY"
	TrelloToken = "TRELLO_TOKEN"
)

type config struct {
	Trello *trello.Config
	Gui    *gui.Config
}

func setup() config {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if len(os.Args) < 2 {
		log.Fatal().Msg("Board name not provided")
	}

	return config{
		Trello: &trello.Config{
			User:            os.Getenv(TrelloUser),
			Key:             os.Getenv(TrelloKey),
			Token:           os.Getenv(TrelloToken),
			Board:           os.Args[1],
			RefreshInterval: time.Second * 20,
		},
		Gui: &gui.Config{},
	}
}

func main() {
	var (
		cfg = setup()
		t   = trello.NewClient(log.Logger.With().Str("m", "trello").Logger(), cfg.Trello)
		a   = gui.NewApp(log.Logger.With().Str("m", "gui").Logger(), cfg.Gui, t)
	)

	log.Info().Msg("Starting trello-cli")

	// Init client and start refresh loop in go routine
	if err := t.Init(); err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	go t.RefreshLoop(ctx)
	defer cancel()

	// Init gui and run it
	err := a.Init()
	defer a.Close()
	if err != nil {
		return
	}
	if err := a.Run(); err != nil {
		return
	}
}
