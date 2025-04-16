<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="../assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Feature Flags: Specification v2 (OpenFeature-Compatible)

This document defines the **v2 feature flag model** used by `ducto-featureflags`, designed to align with the [OpenFeature specification](https://openfeature.dev) and support typed, multi-variant, and rule-based flag evaluation.

## ✅ Goals
- Support OpenFeature Go SDK provider interface
- Enable typed resolution: boolean, string, number, object
- Introduce named variants and defaultVariant
- Allow conditional resolution through rule targeting
- Preserve compatibility with v1 flags when possible

---
## 📦 Flag Model

```json
{
  "new_ui": {
    "enabled": true,
    "defaultVariant": "off",
    "variants": {
      "on": true,
      "off": false
    },
    "rules": [
      {
        "if": { "env": "prod", "group": "beta" },
        "variant": "on"
      },
      {
        "if": { "env": "prod" },
        "variant": "off"
      }
    ]
  }
}
```

### Flag Fields
| Field            | Type                     | Description                         |
|------------------|--------------------------|-------------------------------------|
| `disabled`       | `bool` (default false)   | Whether the flag is active          |
| `defaultVariant` | `string`                 | Fallback variant if no rule matches |
| `variants`       | `map[string]interface{}` | Named, typed variant values         |
| `rules`          | `[]VariantRule`          | Targeted resolution logic           |

### VariantRule Fields
| Field       | Type                      | Description                                      |
|-------------|---------------------------|--------------------------------------------------|
| `if`        | `map[string]string`       | Context matchers for conditional activation      |
| `percent`   | `int` (0–100)             | Optional: percent rollout gate                   |
| `seed`      | `string`                  | Seed key from context                            |
| `seed_hash` | `"sha256"` (optional)     | Optional hash function                           |
| `variant`   | `string`                  | Name of the variant to return if matched         |

---
## 🧠 Rule Evaluation

1. If `enabled` is `false`, return `defaultVariant`
2. Evaluate rules in order:
    - If `if` matches, and/or `percent` check passes → return `variant`
3. If no rules match, return `defaultVariant`

---
## ✅ Supported Types

OpenFeature-compatible SDKs expect typed resolution:
- Boolean: `true` / `false`
- String: `"green"`, `"beta"`
- Number: `10`, `3.14`
- Object: `{...}`

The value returned from `variants[variant]` must be coercible to the requested type.

---
## 🧪 Evaluation Context

```json
{
  "env": "prod",
  "group": "beta",
  "user_id": "12345"
}
```

---
## 🔁 YAML Example

```yaml
new_ui:
  defaultVariant: off
  variants:
    on: true
    off: false
  rules:
    - if:
        env: prod
        group: beta
      variant: on
    - if:
        env: prod
      variant: off
```

---
## 💪 Supporting Object Variants

Values are no longer restricted to booleans, so variants can include numbers, strings
and even objects like so:

```json
{
  "checkout_config": {
    "defaultVariant": "standard",
    "variants": {
      "standard": {
        "checkout_timeout": 30,
        "retry": false
      },
      "beta": {
        "checkout_timeout": 10,
        "retry": true
      }
    },
    "rules": [
      {
        "if": { "group": "beta" },
        "variant": "beta"
      }
    ]
  }
}
```

---
## 📌 Compatibility with v1

- v1 boolean flags can be upgraded by mapping:
    - `enabled` → `defaultVariant`
    - `rules[].value: true|false` → `variant: "on"|"off"`
    - `variants: { on: true, off: false }`

---
## 🧩 Limitations
- `rules` are flat: no nested `and`/`or`/`not` (same as v1)
- One variant per matching rule
- Requires OpenFeature consumers to respect the type expectations

---
## 🧭 Summary

The v2 feature flag engine introduces typed variant support and full OpenFeature compatibility. Evaluation resolves named variants based on rule context, with fallback to a default. This builds on v1’s boolean-only engine to support modern use cases like experimentation, segmentation, and multi-type resolution across platforms.

This specification supports:
- `BooleanEvaluation`
- `StringEvaluation`
- `IntEvaluation`
- `FloatEvaluation`
- `ObjectEvaluation`

It is the foundation for the `ducto-featureflags/openfeature` provider module.
