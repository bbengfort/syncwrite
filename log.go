package syncwrite

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

// Log defines the methods for an append-only log that must be synchronized.
type Log interface {
	Open(path string) error
	Append(value []byte) error
	Get(index int) (*Entry, error)
	Close() error
}

//===========================================================================
// Base log (not thread-safe)
//===========================================================================

type log struct {
	entries []*Entry
}

// Initialize the internal entries data structure
func (l *log) init() {
	l.entries = make([]*Entry, 0)
}

// Creates an entry for the specified value and appends it to the entries.
func (l *log) create(value []byte) *Entry {
	entry := &Entry{
		Index: len(l.entries),
		Value: value,
	}
	l.entries = append(l.entries, entry)
	return entry
}

// Get an entry at the index or return nil if there is no entry at index.
func (l *log) get(index int) *Entry {
	if index < len(l.entries) {
		return l.entries[index]
	}
	return nil
}

//===========================================================================
// In-Memory Log
//===========================================================================

// InMemoryLog performs no writes to disk.
type InMemoryLog struct {
	sync.RWMutex
	log
}

// Open initializes the log
func (l *InMemoryLog) Open(path string) error {
	l.init()
	return nil
}

// Close is a no-op
func (l *InMemoryLog) Close() error {
	return nil
}

// Append a value and immediately return
func (l *InMemoryLog) Append(value []byte) error {
	l.Lock()
	defer l.Unlock()
	l.create(value)
	return nil
}

// Get a value by its index
func (l *InMemoryLog) Get(index int) (*Entry, error) {
	l.RLock()
	defer l.RUnlock()

	e := l.get(index)
	if e == nil {
		return nil, fmt.Errorf("no entry at index %d", index)
	}

	if e.Index != index {
		return nil, errors.New("log is not consistent with entries")
	}
	return e, nil
}

//===========================================================================
// File Log
//===========================================================================

// FileLog writes log entries one per line in append mode.
type FileLog struct {
	sync.RWMutex
	log
	file *os.File
}

// Open the file and read the entries from disk.
func (l *FileLog) Open(path string) (err error) {
	l.Lock()
	defer l.Unlock()

	// Initialize entries and read them from disk
	l.read(path)

	// Open the file in append, ready for writing
	l.file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	return err
}

// Read the entries from disk
func (l *FileLog) read(path string) {
	l.init() // start with a fresh entries log

	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		entry := new(Entry)
		if err = entry.Load(sc.Bytes()); err != nil {
			return
		}
		l.entries = append(l.entries, entry)
	}
}

// Close the file and set the file object to nil to indicate no more writes.
func (l *FileLog) Close() (err error) {
	l.Lock()
	defer l.Unlock()

	if err = l.file.Close(); err != nil {
		return err
	}
	l.file = nil
	return nil
}

// Append a value and immediately return
func (l *FileLog) Append(value []byte) error {
	l.Lock()
	defer l.Unlock()

	if l.file == nil {
		return errors.New("log file has been closed ")
	}

	// Create and dump entry, appending a newline
	entry := l.create(value)
	data, err := entry.Dump()
	if err != nil {
		return err
	}
	data = append(data, byte('\n'))

	if _, err = l.file.Write(data); err != nil {
		return err
	}
	return nil
}

// Get a value by its index
func (l *FileLog) Get(index int) (*Entry, error) {
	l.RLock()
	defer l.RUnlock()

	e := l.get(index)
	if e == nil {
		return nil, fmt.Errorf("no entry at index %d", index)
	}

	if e.Index != index {
		return nil, errors.New("log is not consistent with entries")
	}
	return e, nil
}

//===========================================================================
// LevelDB Log
//===========================================================================

// LevelDBLog writes log entries to a leveldb store.
type LevelDBLog struct {
	sync.RWMutex
	db      *leveldb.DB
	lastIdx int
}

// Open the file and read the entries from disk.
func (l *LevelDBLog) Open(path string) (err error) {
	l.Lock()
	defer l.Unlock()

	l.db, err = leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}

	// Find the biggest key in the database
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		index := l.readKey(iter.Key())
		if index > l.lastIdx {
			l.lastIdx = index
		}
	}
	iter.Release()
	return iter.Error()
}

// Close the file and set the file object to nil to indicate no more writes.
func (l *LevelDBLog) Close() (err error) {
	l.Lock()
	defer l.Unlock()

	if err = l.db.Close(); err != nil {
		return err
	}
	l.db = nil
	return nil
}

// Append a value and immediately return
func (l *LevelDBLog) Append(value []byte) error {
	l.Lock()
	defer l.Unlock()

	if l.db == nil {
		return errors.New("log database has been closed ")
	}

	// Create and dump entry to bytes, incrementing last index
	l.lastIdx++
	entry := &Entry{
		Index: l.lastIdx,
		Value: value,
	}
	data, err := entry.Dump()
	if err != nil {
		return err
	}

	// Put the entry into the database
	key := l.makeKey(entry.Index)
	return l.db.Put(key, data, nil)
}

// Get a value by its index
func (l *LevelDBLog) Get(index int) (*Entry, error) {
	l.RLock()
	defer l.RUnlock()

	if l.db == nil {
		return nil, errors.New("log database has been closed ")
	}

	data, err := l.db.Get(l.makeKey(index), nil)
	if err != nil {
		return nil, err
	}

	entry := new(Entry)
	if err = entry.Load(data); err != nil {
		return nil, err
	}

	if entry.Index != index {
		return nil, errors.New("log is not consistent with entries")
	}
	return entry, nil
}

func (l *LevelDBLog) makeKey(index int) []byte {
	return []byte(strconv.Itoa(index))
}

func (l *LevelDBLog) readKey(key []byte) int {
	index, _ := strconv.Atoi(string(key))
	return index
}
