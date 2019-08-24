package trello

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	User, Key, Token, Board string
	RefreshInterval         time.Duration
}

type Client struct {
	l      zerolog.Logger
	cfg    *Config
	client *trello.Client
	user   *trello.Member
	store  *JSONFileStore
	errs   []error
}

func NewClient(log zerolog.Logger, cfg *Config) *Client {
	return &Client{
		l:     log,
		cfg:   cfg,
		store: NewJSONFileStore("./.trello-cli"),
	}
}

func (t *Client) Init() error {
	t.l.Debug().Msg("Initializing client")
	err := t.initClient()
	if err != nil {
		t.l.Error().Err(err).Msg("Could not initialize client")
		return err
	}

	err = t.initUser()
	if err != nil {
		t.l.Error().Err(err).Msg("Could not get user info")
		return err
	}

	err = t.initStore()
	if err != nil {
		t.l.Error().Err(err).Msg("Could not initialize store")
		return err
	}

	t.l.Debug().Msg("Client initialized")
	return nil
}

func (t *Client) RefreshLoop(ctx context.Context) {
	t.l.Debug().Msg("Starting refresh loop")
	var timer = time.NewTicker(t.cfg.RefreshInterval)
infiniteLoop:
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			break infiniteLoop
		case <-timer.C:
			err := t.refreshStore()
			if err != nil {

				t.l.Error().Err(err).Msg("Unexpected error while refreshing board")
				t.errs = append(t.errs, err)
				continue
			}
		}
	}
	t.l.Debug().Msg("Exiting refresh loop")
}

func (t *Client) Name() string {
	t.store.startRead()
	defer t.store.endRead()
	return t.store.Board.Name
}

func (t *Client) Description() string {
	t.store.startRead()
	defer t.store.endRead()
	return t.store.Board.Desc
}

func (t *Client) ListsNames() []string {
	t.store.startRead()
	defer t.store.endRead()
	var listNames = make([]string, len(t.store.Lists))
	for i, l := range t.store.Lists {
		listNames[i] = l.Name
	}
	return listNames
}

func (t *Client) initClient() error {
	client, err := trello.NewAuthClient(t.cfg.Key, &t.cfg.Token)
	if err != nil {
		return errors.Wrapf(err, "could not initialize auth client for member %s", t.cfg.User)
	}
	t.client = client
	return nil
}

func (t *Client) initUser() error {
	user, err := t.client.Member(t.cfg.User)
	if err != nil {
		return errors.Wrap(err, "could not get user info")
	}
	t.user = user
	return nil
}

func (t *Client) initStore() error {
	err := t.store.Init()
	if err != nil {
		return errors.Wrap(err, "while initializing store")
	}

	err = t.refreshStoreIfOld()
	if err != nil {
		return errors.Wrap(err, "while refreshing store")
	}

	return nil
}

func (t *Client) refreshStoreIfOld() error {
	if time.Since(t.store.LastUpdated) > t.cfg.RefreshInterval {
		err := t.refreshStore()
		if err != nil {
			return errors.Wrapf(err, "while refreshing old store")
		}
	}
	return nil
}

func (t *Client) refreshStore() error {
	boards, err := t.user.Boards()
	if err != nil {
		return errors.Wrap(err, "could not get boards")
	}

	var board *trello.Board
	for _, b := range boards {
		if strings.EqualFold(b.Name, t.cfg.Board) {
			board = &b
			break
		}
	}
	if board == nil {
		return errors.Wrapf(errors.New(fmt.Sprintf("board %s not found", t.cfg.Board)), "while getting board %s", t.cfg.Board)
	}

	lists, err := board.Lists()
	if err != nil {
		return errors.Wrapf(err, "while getting list for board %s", board.Name)
	}
	err = t.store.Update(*board, lists)
	if err != nil {
		return errors.Wrapf(err, "while updating store", board.Name)
	}

	t.l.Debug().Msg("Store refreshed")
	return nil
}
