package sdk

import (
	"context"
	"sync"
)

// StoreProvider defines something that can provide and watch a Store (e.g., from file, http, etc.)
type StoreProvider interface {
	Load(ctx context.Context) (*Store, error)
	Watch(ctx context.Context, onChange func(*Store))
}

// DynamicStore is a AnyStore which wraps a StoreProvider and handles live updates to the internal store.
// You can call it in the same way you call Store, except you need to call Start first.
type DynamicStore struct {
	mu     sync.RWMutex
	store  *Store
	source StoreProvider
	ctx    context.Context
}

// NewDynamicStore creates a dynamic flag store that tracks updates from the provider.
func NewDynamicStore(ctx context.Context, provider StoreProvider) *DynamicStore {
	return &DynamicStore{
		source: provider,
		ctx:    ctx,
	}
}

// Start begins watching the underlying source. Should be called once.
func (d *DynamicStore) Start() error {
	initial, err := d.source.Load(d.ctx)
	if err != nil {
		return err
	}
	d.store = initial

	go d.source.Watch(d.ctx, func(updated *Store) {
		if updated == nil {
			return
		}
		d.mu.Lock()
		d.store = updated
		d.mu.Unlock()
	})

	return nil
}

// IsEnabled evaluates a feature flag for the given context.
func (d *DynamicStore) IsEnabled(key string, ctx EvalContext) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.store.IsEnabled(key, ctx)
}

// AllFlags returns all current flag definitions.
func (d *DynamicStore) AllFlags() map[string]Flag {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.store.AllFlags()
}
