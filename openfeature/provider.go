package openfeature

import (
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/tommed/ducto-featureflags/sdk"
)

// DuctoProvider implements the OpenFeature Provider interface
type DuctoProvider struct {
	Store sdk.AnyStore // your existing flag store (interface)
}

func NewProvider(store sdk.AnyStore) openfeature.FeatureProvider {
	return &DuctoProvider{Store: store}
}

func (p *DuctoProvider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{Name: "ducto-featureflags"}
}

func (p *DuctoProvider) Hooks() []openfeature.Hook {
	return nil
}
