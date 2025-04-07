package sdk

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchingStore wraps Store and reloads it when the underlying file changes.
type WatchingStore struct {
	mu     sync.RWMutex
	store  *Store
	path   string
	events chan struct{}
	errors chan error
	done   chan struct{}
}

// NewWatchingStore creates a WatchingStore that watches the given file.
func NewWatchingStore(path string) (*WatchingStore, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	s := &WatchingStore{
		path:   absPath,
		events: make(chan struct{}, 1),
		errors: make(chan error, 1),
		done:   make(chan struct{}),
	}

	store, err := NewStoreFromFile(absPath)
	if err != nil {
		return nil, err
	}
	s.store = store

	go s.watch()
	return s, nil
}

func (s *WatchingStore) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.errors <- fmt.Errorf("watcher init failed: %w", err)
		return
	}
	defer watcher.Close()

	dir := s.watchDir()
	if err := watcher.Add(dir); err != nil {
		s.errors <- fmt.Errorf("watch add failed: %w", err)
		return
	}

	for {
		select {
		case <-s.done:
			return
		case evt := <-watcher.Events:
			evtPath := filepath.Clean(evt.Name)
			watchPath := filepath.Clean(s.path)

			if evtPath == watchPath &&
				evt.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				s.reload()
			}
		case err := <-watcher.Errors:
			s.errors <- err
		}
	}
}

func (s *WatchingStore) watchDir() string {
	return filepath.Dir(s.path)
}

func (s *WatchingStore) reload() {
	time.Sleep(50 * time.Millisecond) // debounce file locks
	store, err := NewStoreFromFile(s.path)
	if err != nil {
		s.errors <- fmt.Errorf("reload failed: %w", err)
		return
	}

	s.mu.Lock()
	s.store = store
	s.mu.Unlock()

	select {
	case s.events <- struct{}{}:
	default:
	}
}

// IsEnabled implements the Store interface.
func (s *WatchingStore) IsEnabled(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store.IsEnabled(key)
}

func (s *WatchingStore) AllFlags() map[string]Flag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store.AllFlags()
}

func (s *WatchingStore) Close() {
	close(s.done)
}

// Events exposes reload notifications (optional).
func (s *WatchingStore) Events() <-chan struct{} {
	return s.events
}

// Errors exposes watcher errors (optional).
func (s *WatchingStore) Errors() <-chan error {
	return s.errors
}
