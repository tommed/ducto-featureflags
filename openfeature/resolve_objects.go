package openfeature

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
)

func (p *DuctoProvider) ObjectEvaluation(
	_ context.Context,
	flagKey string,
	defaultValue interface{},
	evalCtx openfeature.FlattenedContext,
) openfeature.InterfaceResolutionDetail {
	flagDef, found := p.Store.Get(flagKey)
	if !found {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(flagKey),
			},
		}
	}

	internalCtx := convertFlattenedContext(evalCtx)
	variant, val, ok, matched := flagDef.Evaluate(internalCtx)
	if !ok {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewParseErrorResolutionError("variant not found"),
			},
		}
	}

	reason := openfeature.DefaultReason
	if matched {
		reason = openfeature.TargetingMatchReason
	}

	return openfeature.InterfaceResolutionDetail{
		Value: val,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Variant: variant,
			Reason:  reason,
		},
	}
}
