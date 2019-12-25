package trello

import (
	"net/http"
	"strings"

	"github.com/VojtechVitek/go-trello"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/giannimassi/trello-tui/pkg/domain"
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
	return &Client{
		l:   log.With().Str("m", "trello").Str("user", cfg.User).Logger(),
		cfg: cfg,
	}
}

// Init setups the client with the configuration provided, connection to trello's backend
// is initialized on the first use
func (t *Client) Init() error {
	rr := trello.NewBearerTokenTransport(t.cfg.Key, &t.cfg.Token)
	httpClient := &http.Client{
		Transport: rr,
		Timeout:   t.cfg.Timeout,
	}

	client, err := trello.NewCustomClient(httpClient)
	if err != nil {
		t.l.Error().Str("user", t.cfg.User).Err(err).Msg("Could not initialize auth client")
		return errors.Wrapf(err, "could not initialize auth client for member %s", t.cfg.User)
	}
	t.client = client
	t.l.Info().Msg("Initialized")
	return nil
}

// Board returns a domain.Board populated with the latest info about the specified board
func (t *Client) Board(name string) (*domain.Board, error) {
	t.l.Debug().Msg("Getting boards")
	var board trello.Board

	user, err := t.client.Member(t.cfg.User)
	if err != nil {
		return nil, errors.Wrap(err, "could not get boards")
	}

	boards, err := user.Boards()
	if err != nil {
		return nil, ErrBoardNotFound
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
		return nil, ErrBoardNotFound
	}

	t.l.Debug().Interface("board", board.Name).Msg("Getting lists for board")
	lists, err := board.Lists()
	if err != nil {
		return nil, errors.Wrapf(err, "while getting list for board %s", board.Name)
	}

	cards, err := board.Cards()
	if err != nil {
		return nil, errors.Wrapf(err, "while getting cards for board %s", board.Name)
	}

	return domain.NewBoard(board.Name, board.Desc, listsByID(lists, cardsByListID(cards)), len(cards) == 0), nil
}

func listsByID(trelloCards []trello.List, cards map[string]map[int]domain.Card) []domain.List {
	lists := make([]domain.List, 0, len(trelloCards))
	for _, c := range trelloCards {
		lists = append(lists, domain.NewList(c.Id, c.Name, cards[c.Id]))
	}
	return lists
}

func cardsByListID(trelloCards []trello.Card) map[string]map[int]domain.Card {
	cards := make(map[string]map[int]domain.Card)
	for _, c := range trelloCards {
		if cards[c.IdList] == nil {
			cards[c.IdList] = make(map[int]domain.Card)
		}
		labels := make([]domain.CardLabel, len(c.Labels))
		for i, lbl := range c.Labels {
			labels[i] = domain.CardLabel{Name: lbl.Name, Color: lbl.Color}
		}

		cards[c.IdList][c.IdShort] = domain.NewCard(c.Id, c.Name, c.Desc, c.Pos, labels)
	}
	return cards
}
