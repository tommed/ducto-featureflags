<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="../assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Feature Flags: Specification v1

This specification defines the **v1 structure and behavior** of the feature flag engine used by `ducto-featureflags`. It documents the design as originally implemented prior to OpenFeature compatibility and multi-variant support.

## ✅ Goals
- Provide simple boolean feature toggles
- Support conditional evaluation based on dynamic context
- Allow percentage rollouts and rule precedence
- Support file, embedded, and HTTP-served flags

---

## 📦 Flag Model

```json
{
  "new_ui": {
    "enabled": true,
    "rules": [
      { "if": { "env": "prod", "group": "beta" }, "value": true },
      { "if": { "env": "prod" }, "value": false }
    ]
  }
}
```

### Flag Fields
| Field     | Type             | Description                                      |
|-----------|------------------|--------------------------------------------------|
| `enabled` | `bool` (pointer) | Default value if no rule matches (not a blocker) |
| `rules`   | `[]Rule`         | Ordered rules to evaluate against context        |

### Rule Fields
| Field       | Type                  | Description                                  |
|-------------|-----------------------|----------------------------------------------|
| `if`        | `map[string]string`   | Match context values by key                  |
| `value`     | `bool`                | Boolean value to return if rule matches      |
| `percent`   | `int` (0–100)         | Optional: percent rollout gate (uses `seed`) |
| `seed`      | `string`              | Context key to seed hash for percent logic   |
| `seed_hash` | `"sha256"` (optional) | Use SHA256 instead of default hash           |

---

## 🧠 Rule Evaluation

1. `enabled` acts as a fallback if no rule matches (or is nil)
2. Rules are evaluated **in order**
3. For each rule:
    - If `if` is present: all key-values must match the evaluation context
    - If `percent` is present: a seeded hash is used to gate by rollout percentage
    - If the rule matches, return its `value`
4. If no rule matches:
    - If `enabled != nil`, return `*enabled`
    - Else return `false`

Note: `enabled: false` does **not** short-circuit rule evaluation. Rules can override it.

---

## 🧪 Evaluation Context

The context is a `map[string]string` provided at evaluation time:

```json
{
  "env": "prod",
  "group": "beta",
  "user_id": "12345"
}
```

This is used by rules to determine eligibility based on traits.

---

## 🔁 YAML Format Example
```yaml
new_ui:
  enabled: false
  rules:
    - if:
        env: "dev"
      value: true
    - if:
        group: "beta"
      value: true
    - if:
        env: "prod"
      value: false
```

---

## 🎯 Percentage Rollout Example

```json
{
  "experiment": {
    "enabled": false,
    "rules": [
      {
        "percent": 20,
        "seed": "user_id",
        "seed_hash": "sha256",
        "value": true
      }
    ]
  }
}
```

This flag enables `true` for ~20% of users based on a consistent hash of their `user_id`.

---

## 🖥️ Serving

The engine can:
- Evaluate flags in memory (via `IsEnabled(key, context)`)
- Be embedded in CLI tools
- Be served via HTTP API:

```http
GET /api/flags?key=new_ui
```

Returns:
```json
{ "enabled": true }
```

Or:
```http
GET /api/flags
```
Returns raw flag definitions.

---

## 🧩 Limitations in v1
- Only supports **boolean** flags
- No typed resolution (string, int, object)
- No named variants (e.g., "blue", "green")
- No hooks, metadata, reason codes
- No OpenFeature provider support (planned for v2)

---

## ✅ Summary

This spec formalizes the **v1 feature flag engine** inside `ducto-featureflags`, which supports boolean flags, conditional evaluation via simple `if` rules, and controlled rollout via percentage-based gating. Rules may override the static `enabled` field, and `enabled` serves as a fallback only. This forms the foundation for v2, which will support variants and OpenFeature compliance.