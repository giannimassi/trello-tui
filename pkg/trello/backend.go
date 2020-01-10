package trello

import (
	"context"
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/giannimassi/trello-tui/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is the state configuration, including trello authentication details
type Config struct {
	User, Key, Token     string
	Timeout              time.Duration
	SelectedBoard        string
	BoardRefreshInterval time.Duration
}

// Backend ensures the state is updated as required
type Backend struct {
	l      zerolog.Logger
	cfg    *Config
	client *Client
	store  *store.Store

	lastActionUpdate map[int]time.Time

	// Buffered channel drained in Run
	refreshCardActionsCh chan *trello.Card
}

// NewBackend returns a new instance of Backend
func NewBackend(cfg *Config) *Backend {
	u := Backend{
		l:                    log.Logger.With().Str("m", "state-update").Logger(),
		cfg:                  cfg,
		client:               NewClient(cfg),
		refreshCardActionsCh: make(chan *trello.Card, 100),
		lastActionUpdate:     make(map[int]time.Time),
	}
	u.store = store.NewStore(&boardLoading{
		boardName:             cfg.SelectedBoard,
		cardCommentsRequested: u.cardCommentsRequested,
	})

	return &u
}

// should not block for
func (u *Backend) cardCommentsRequested(card *trello.Card) {
	// send on channel with timeout
	select {
	case u.refreshCardActionsCh <- card:
		return
	case <-time.After(time.Second):
		return
	}
}

// Run executes the loop exuting the requests via the trello client
func (u *Backend) Run(ctx context.Context) {
	var (
		t = time.NewTimer(0)
	)
	u.l.Debug().Msg("Refresh state started")
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return

		case now := <-t.C:
			u.l.Debug().Msg("Refreshing board, lists and cards")
			if err := u.updateBoard(); err != nil {
				u.l.Error().Err(err).Msg("Could not update board")
			}
			t.Reset(u.cfg.BoardRefreshInterval - time.Since(now))
		case id := <-u.refreshCardActionsCh:
			u.l.Debug().Msg("Refreshing card actions")
			if err := u.updateCardActions(id); err != nil {
				u.l.Error().Err(err).Msg("Could not update card actions")
			}
		}
	}
}

// Store returns the Store instance to which new data is stored
func (u *Backend) Store() *store.Store {
	return u.store
}

// FIXME: update comment
// updateBoard returns a copy of the current state with updated board / err
func (u *Backend) updateBoard() error {
	b, lists, cards, err := u.client.Board(u.cfg.SelectedBoard)
	if err != nil {
		u.offline(err)
		return err
	}
	u.online(b, lists, cards)
	return nil
}

// FIXME: add comment
func (u *Backend) updateCardActions(card *trello.Card) error {
	// check when card actions where last updated
	lastUpdated, found := u.lastActionUpdate[card.IdShort]
	// TODO: separate refresh interval for actions?
	if found && time.Since(lastUpdated) < u.cfg.BoardRefreshInterval {
		return nil
	}

	actions, err := u.client.CardActions(card)
	if err != nil {
		return err
	}
	updated := time.Now()
	u.store.BeginWrite()
	switch t := u.store.State.(type) {
	case *boardOnline:
		u.store.State = t.setCardActions(card.IdShort, actions)
	case *boardOffline:
		u.store.State = t.setCardActions(card.IdShort, actions)
	default:
		log.Error().Msg("unexpected board state")
		u.store.EndWrite(false)
		return nil
	}
	u.store.EndWrite(true)
	u.lastActionUpdate[card.IdShort] = updated
	return nil
}

func (u *Backend) online(newBoard *trello.Board, lists []trello.List, cards []trello.Card) {
	u.store.BeginWrite()
	switch t := u.store.State.(type) {
	case *boardLoading:
		u.store.State = t.online(newBoard, lists, cards)
	case *boardOnline:
		u.store.State = t.online(newBoard, lists, cards)
	case *boardOffline:
		u.store.State = t.online(newBoard, lists, cards)
	case *boardLoadingOffline:
		u.store.State = t.online(newBoard, lists, cards)
	default:
		log.Error().Msg("wrong type")
		u.store.EndWrite(false)
		return
	}
	u.store.EndWrite(true)
}

func (u *Backend) offline(err error) {
	u.store.BeginWrite()
	switch t := u.store.State.(type) {
	case *boardLoading:
		u.store.State = t.offline(err)
	case *boardOnline:
		u.store.State = t.offline(err)
	case *boardOffline:
		u.store.State = t.offline(err)
	case *boardLoadingOffline:
		u.store.State = t.offline(err)
	default:
		log.Error().Msg("wrong type")
		u.store.EndWrite(false)
		return
	}
	u.store.EndWrite(true)
}
