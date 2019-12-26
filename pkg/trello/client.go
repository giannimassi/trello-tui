package trello

import (
	"net/http"
	"strings"

	"github.com/VojtechVitek/go-trello"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ErrBoardNotFound is returned when a board is not found
var ErrBoardNotFound = errors.New("board not found")

// Client makes requests via the trello API to get data for the current user
type Client struct {
	l      zerolog.Logger
	cfg    *Config
	client *trello.Client
}

// NewClient returns a new instance of Client
func NewClient(cfg *Config) *Client {
	httpClient := &http.Client{
		Transport: trello.NewBearerTokenTransport(cfg.Key, &cfg.Token),
		Timeout:   cfg.Timeout,
	}
	// Error not handled since NewCustomClient implementation always returns nil
	client, _ := trello.NewCustomClient(httpClient)
	return &Client{
		l:      log.With().Str("m", "trello").Str("user", cfg.User).Logger(),
		cfg:    cfg,
		client: client,
	}
}

// Board returns a domain.Board populated with the latest info about the specified board
func (t *Client) Board(name string) (*trello.Board, []trello.List, []trello.Card, error) {
	t.l.Debug().Msg("Getting boards")
	var board trello.Board
	user, err := t.client.Member(t.cfg.User)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not get boards")
	}

	boards, err := user.Boards()
	if err != nil {
		return nil, nil, nil, ErrBoardNotFound
	}

	var boardFound bool
	for _, b := range boards {
		if strings.EqualFold(b.Name, name) {
			board = b
			boardFound = true
			break
		}
	}
	if !boardFound {
		return nil, nil, nil, ErrBoardNotFound
	}

	t.l.Debug().Interface("board", board.Name).Msg("Getting lists for board")
	lists, err := board.Lists()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "while getting list for board %s", board.Name)
	}

	cards, err := board.Cards()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "while getting cards for board %s", board.Name)
	}
	return &board, lists, cards, nil
}
