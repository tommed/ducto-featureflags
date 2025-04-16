package openfeature

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
)

func (p *DuctoProvider) IntEvaluation(
	_ context.Context,
	flagKey string,
	defaultValue int64,
	evalCtx openfeature.FlattenedContext,
) openfeature.IntResolutionDetail {
	flagDef, found := p.Store.Get(flagKey)
	if !found {
		return openfeature.IntResolutionDetail{
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
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewParseErrorResolutionError("variant not found"),
			},
		}
	}

	var n int64
	switch v := result.Value.(type) {
	case int:
		n = int64(v)
	case int64:
		n = v
	case float64:
		n = int64(v)
	default:
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewTypeMismatchResolutionError("int"),
			},
		}
	}

	reason := openfeature.DefaultReason
	if result.Matched {
		reason = openfeature.TargetingMatchReason
	}

	return openfeature.IntResolutionDetail{
		Value: n,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Variant: result.Variant,
			Reason:  reason,
		},
	}
}

func (p *DuctoProvider) FloatEvaluation(
	_ context.Context,
	flagKey string,
	defaultValue float64,
	evalCtx openfeature.FlattenedContext,
) openfeature.FloatResolutionDetail {
	flagDef, found := p.Store.Get(flagKey)
	if !found {
		return openfeature.FloatResolutionDetail{
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
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewParseErrorResolutionError("variant not found"),
			},
		}
	}

	var f float64
	switch v := result.Value.(type) {
	case float64:
		f = v
	case float32:
		f = float64(v)
	case int:
		f = float64(v)
	case int64:
		f = float64(v)
	default:
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Variant:         result.Variant,
				Reason:          openfeature.DefaultReason,
				ResolutionError: openfeature.NewTypeMismatchResolutionError("float"),
			},
		}
	}

	reason := openfeature.DefaultReason
	if result.Matched {
		reason = openfeature.TargetingMatchReason
	}

	return openfeature.FloatResolutionDetail{
		Value: f,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Variant: result.Variant,
			Reason:  reason,
		},
	}
}
