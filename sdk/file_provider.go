package sdk

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func NewFileProvider(path string) StoreProvider {
	return &fileProvider{path: path}
}

// fileProvider implements StoreProvider by watching a file on disk.
type fileProvider struct {
	path     string
	last     *Store
	lastLock sync.RWMutex
}

// Load loads the current store from disk.
func (f *fileProvider) Load(_ context.Context) (*Store, error) {
	absPath, err := filepath.Abs(f.path)
	if err != nil {
		return nil, err
	}
	store, err := NewStoreFromFile(absPath)
	if err != nil {
		return nil, err
	}

	f.lastLock.Lock()
	f.last = store
	f.lastLock.Unlock()

	return store, nil
}

func (f *fileProvider) Watch(ctx context.Context, onChange func(*Store)) {
	absPath, err := filepath.Abs(f.path)
	if err != nil {
		return
	}

	dir := filepath.Dir(absPath)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	//goland:noinspection GoUnhandledErrorResult
	defer watcher.Close()
	_ = watcher.Add(dir)

	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-watcher.Events:
			evtPath := filepath.Clean(evt.Name)
			if evtPath != absPath {
				continue
			}
			if evt.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				time.Sleep(50 * time.Millisecond) // debounce
				store, err := f.Load(ctx)
				if err == nil {
					onChange(store)
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("file watcher error: %v\n", err) // could expose as metric/log
		}
	}
}
