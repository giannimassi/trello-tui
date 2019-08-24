package trello

import (
	"encoding/json"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/pkg/errors"
)

const (
	DefaultFileName = ".trello-cli"
)

type JSONFileStore struct {
	store
	path string        // Path to the file used for this store
	m    *sync.RWMutex // Mutex to take care of race conditions
}

func NewJSONFileStore(path string) *JSONFileStore {
	if path == "" {
		path = defaultStorePath()
	}

	return &JSONFileStore{
		path: path,
		m:    new(sync.RWMutex),
	}
}

func (s *JSONFileStore) startRead()  { s.m.RLock() }
func (s *JSONFileStore) endRead()    { s.m.RUnlock() }
func (s *JSONFileStore) startWrite() { s.m.Lock() }
func (s *JSONFileStore) endWrite()   { s.m.Unlock() }

func (s *JSONFileStore) Init() error {
	s.startWrite()
	defer s.endWrite()
	err := s.readAllUnlocked()
	if err != nil {
		return errors.Wrapf(err, "while reading from store file")
	}

	return nil
}

func (s *JSONFileStore) Update(board trello.Board, lists []trello.List) error {
	s.startWrite()
	defer s.endWrite()
	s.Board = board
	s.Lists = lists
	s.LastUpdated = time.Now()
	return s.writeAllUnlocked()
}

func (s *JSONFileStore) readAllUnlocked() error {
	f, err := openOrCreate(s.path)
	if err != nil {
		return errors.Wrapf(err, "while opening store at %s", s.path)
	}
	defer f.Close()

	err = decodeUserStore(f, &s.store)
	if err != nil {
		return errors.Wrapf(err, "while decoding store at %s", s.path)
	}
	return nil
}

func (s *JSONFileStore) writeAllUnlocked() error {
	f, err := openOrCreate(s.path)
	if err != nil {
		return errors.Wrapf(err, "while opening file store at %s", s.path)
	}
	defer f.Close()

	// Write to file
	err = encodeUserStore(f, &s.store)
	if err != nil {
		return errors.Wrapf(err, "while encoding store at %s", s.path)
	}

	return nil
}

type store struct {
	LastUpdated time.Time
	Board       trello.Board
	Lists       []trello.List
}

func decodeUserStore(r io.Reader, s *store) error {
	d := json.NewDecoder(r)
	err := d.Decode(&s)
	if err == io.EOF {
		return nil
	}

	if err != nil && err != io.EOF {
		return errors.Wrapf(err, "while decoding store")
	}

	return nil
}

func encodeUserStore(w io.Writer, s *store) error {
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
			return nil, errors.Wrapf(err, "while creating %s", path)
		}
		return f, nil
	}
	// Open existing file
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrapf(err, "while opening %s", path)
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
