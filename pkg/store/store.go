package store

import (
	"encoding/json"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/giannimassi/trello-tui/pkg/gui/state"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	DefaultFileName = ".trello-tui"
)

type JSONFileStore struct {
	l     zerolog.Logger
	path  string // Path to the file used for this store
	state *state.State
	m     *sync.RWMutex // Mutex to take care of race conditions
}

func NewJSONFileStore(l zerolog.Logger, path string) *JSONFileStore {
	if path == "" {
		path = defaultStorePath()
	}

	return &JSONFileStore{
		l:     l.With().Str("path", path).Logger(),
		path:  path,
		state: state.NewState(),
		m:     new(sync.RWMutex),
	}
}

func (s *JSONFileStore) StartStateRead()  { s.m.RLock() }
func (s *JSONFileStore) EndStateRead()    { s.m.RUnlock() }
func (s *JSONFileStore) StartStateWrite() { s.m.Lock() }
func (s *JSONFileStore) EndStateWrite()   { s.m.Unlock() }

func (s *JSONFileStore) Init() error {
	s.StartStateWrite()
	defer s.EndStateWrite()
	err := s.readFromFileUnguarded(s.state)
	if err != nil {
		return errors.Wrapf(err, "while reading from file store")
	}
	s.l.Info().Msg("Initialized")
	return nil
}

func (s *JSONFileStore) Write(state *state.State) error {
	s.StartStateWrite()
	defer s.EndStateWrite()
	s.l.Debug().Msg("Saving state to file store")
	if err := s.writeToFileUnguarded(state); err != nil {
		return errors.Wrapf(err, "while writing to file store")
	}
	s.state = state
	return nil
}

func (s *JSONFileStore) State() *state.State {
	s.StartStateRead()
	defer s.EndStateRead()
	return s.state
}

func (s *JSONFileStore) readFromFileUnguarded(store *state.State) error {
	f, err := openOrCreate(s.path)
	if err != nil {
		return errors.Wrapf(err, "while opening store at %state", s.path)
	}
	defer f.Close()

	s.l.Debug().Msg("Decoding store")
	err = decodeStore(f, store)
	if err != nil {
		return errors.Wrapf(err, "while decoding store at %state", s.path)
	}
	return nil
}

func (s *JSONFileStore) writeToFileUnguarded(store *state.State) error {
	f, err := openOrCreate(s.path)
	if err != nil {
		return errors.Wrapf(err, "while opening file store at %state", s.path)
	}
	defer f.Close()

	// Write to file
	err = encodeStore(f, store)
	if err != nil {
		return errors.Wrapf(err, "while encoding store at %state", s.path)
	}

	return nil
}

func decodeStore(r io.Reader, s *state.State) error {
	d := json.NewDecoder(r)
	err := d.Decode(&s)
	if err != nil && err != io.EOF {
		return errors.Wrapf(err, "while decoding store")
	}

	return nil
}

func encodeStore(w io.Writer, s *state.State) error {
	var e = json.NewEncoder(w)
	err := e.Encode(s)
	if err != nil {
		return errors.Wrapf(err, "while encoding store")
	}

	return nil
}

func openOrCreate(path string) (*os.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create
		f, err := os.Create(path)
		if err != nil {
			return nil, errors.Wrapf(err, "while creating %state", path)
		}
		return f, nil
	}
	// Open existing file
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrapf(err, "while opening %state", path)
	}
	return f, nil
}

func defaultStorePath() string {
	usr, err := user.Current()
	if err != nil {
		return "./" + DefaultFileName
	}
	return filepath.Join(usr.HomeDir, DefaultFileName)
}
