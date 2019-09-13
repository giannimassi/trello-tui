package trello

import (
	"strings"

	"github.com/VojtechVitek/go-trello"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	ErrBoardNotFound = errors.New("board not found")
)

type Config struct {
	User, Key, Token string
}

type Client struct {
	l      zerolog.Logger
	cfg    *Config
	client *trello.Client
}

func NewClient(log zerolog.Logger, cfg *Config) *Client {
	return &Client{
		l:   log.With().Str("user", cfg.User).Logger(),
		cfg: cfg,
	}
}

func (t *Client) Init() error {
	client, err := trello.NewAuthClient(t.cfg.Key, &t.cfg.Token)
	if err != nil {
		t.l.Error().Str("user", t.cfg.User).Err(err).Msg("Could not initialize auth client")
		return errors.Wrapf(err, "could not initialize auth client for member %s", t.cfg.User)
	}
	t.client = client
	t.l.Info().Msg("Initialized")
	return nil
}

func (t *Client) BoardInfo(name string) (trello.Board, []trello.List, []trello.Card, error) {
	t.l.Debug().Msg("Getting boards info")
	var board trello.Board

	user, err := t.client.Member(t.cfg.User)
	if err != nil {
		return board, nil, nil, errors.Wrap(err, "could not get boards")
	}

	boards, err := user.Boards()
	if err != nil {
		return board, nil, nil, ErrBoardNotFound
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
		return board, nil, nil, ErrBoardNotFound
	}

	t.l.Debug().Str("board", name).Msg("Getting lists for board")
	lists, err := board.Lists()
	if err != nil {
		return board, nil, nil, errors.Wrapf(err, "while getting list for board %s", board.Name)
	}

	cards, err := board.Cards()
	if err != nil {
		return board, nil, nil, errors.Wrapf(err, "while getting cards for board %s", board.Name)
	}

	return board, lists, cards, nil
}
