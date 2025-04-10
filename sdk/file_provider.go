package sdk

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func NewFileProvider(path string) StoreProvider {
	return NewFileProviderWithLog(path, nil)
}

func NewFileProviderWithLog(path string, writer io.Writer) StoreProvider {
	return &fileProvider{path: path, writer: writer}
}

// fileProvider implements StoreProvider by watching a file on disk.
type fileProvider struct {
	path     string
	last     *Store
	lastLock sync.RWMutex
	writer   io.Writer
}

func (f *fileProvider) logEvent(format string, args ...any) {
	if f.writer != nil {
		_, _ = f.writer.Write([]byte(strings.TrimSpace(fmt.Sprintf(format, args...)) + "\n"))
	}
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
	f.logEvent("Store updated")

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
			f.logEvent("error watching file: %w", err)
		}
	}
}
