package openfeature

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
)

func (p *DuctoProvider) BooleanEvaluation(
	_ context.Context,
	flagKey string,
	defaultValue bool,
	evalCtx openfeature.FlattenedContext,
) openfeature.BoolResolutionDetail {
	flagDef, found := p.Store.Get(flagKey)
	if !found {
		return openfeature.BoolResolutionDetail{
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
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewParseErrorResolutionError("variant not found"),
			},
		}
	}

	b, ok := result.Value.(bool)
	if !ok {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewTypeMismatchResolutionError("bool"),
			},
		}
	}

	reason := openfeature.DefaultReason
	if result.Matched {
		reason = openfeature.TargetingMatchReason
	}

	return openfeature.BoolResolutionDetail{
		Value: b,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Variant: result.Variant,
			Reason:  reason,
		},
	}
}
