package openfeature

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
)

func (p *DuctoProvider) StringEvaluation(
	_ context.Context,
	flagKey string,
	defaultValue string,
	evalCtx openfeature.FlattenedContext,
) openfeature.StringResolutionDetail {
	flagDef, found := p.Store.Get(flagKey)
	if !found {
		return openfeature.StringResolutionDetail{
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
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewParseErrorResolutionError("variant not found"),
			},
		}
	}

	s, ok := result.Value.(string)
	if !ok {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewTypeMismatchResolutionError("string"),
			},
		}
	}

	reason := openfeature.DefaultReason
	if result.Matched {
		reason = openfeature.TargetingMatchReason
	}

	return openfeature.StringResolutionDetail{
		Value: s,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Variant: result.Variant,
			Reason:  reason,
		},
	}
}
