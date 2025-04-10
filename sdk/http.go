package sdk

import "context"

// NewStoreFromURL loads a flag file from an HTTP(S) endpoint, but does not update the flags once acquired.
// Ideally, you would use NewHTTPProvider instead inside a DynamicStore to have an always up-to-date
// copy of the store from a remote location, but this convenience function exists for one-offs if needed.
func NewStoreFromURL(ctx context.Context, url string, token string) (*Store, error) {
	var provider = httpProvider{URL: url, Token: token}
	return provider.Load(ctx)
}
