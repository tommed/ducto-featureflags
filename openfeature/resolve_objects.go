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
	result := flagDef.Evaluate(internalCtx)
	if !result.OK {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewParseErrorResolutionError("variant not found"),
			},
		}
	}

	reason := openfeature.DefaultReason
	if result.Matched {
		reason = openfeature.TargetingMatchReason
	}

	return openfeature.InterfaceResolutionDetail{
		Value: result.Value,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Variant: result.Variant,
			Reason:  reason,
		},
	}
}
