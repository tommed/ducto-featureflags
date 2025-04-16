<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="./assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Ducto Feature Flags

[![CI](https://github.com/tommed/ducto-featureflags/actions/workflows/ci.yml/badge.svg)](https://github.com/tommed/ducto-featureflags/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/tommed/ducto-featureflags/branch/main/graph/badge.svg)](https://codecov.io/gh/tommed/ducto-featureflags)

> Lightweight, embeddable, and pluggable **[OpenFeature](https://openfeature.dev/) compatible** feature flag engine for 
> Go. Designed for pipelines, microservices, and event-driven systems.

---
## ✨ What is Ducto Feature Flags?

`ducto-featureflags` is a minimalist feature flag manager built in Go. It was designed as a reusable component for:
- Data transformation pipelines (like [Ducto](https://github.com/tommed))
- Serverless functions
- Microservices and APIs
- CLI tools and test harnesses

It supports both static file-based flags and dynamic backends, and can be used as:
- A **Go SDK** (directly, or as an [OpenFeature Provider](https://github.com/open-feature/go-sdk/openfeature))
- A **CLI for testing**
- A **preprocessor plugin** in the [ducto-orchestrator](https://github.com/tommed/ducto-orchestrator)

Please note that whilst Ducto Feature Flags is compatible with OpenFeature, 
it is **not** compliant with [flagd](https://flagd.dev/) because we support nested conditional
statements, which _cannot_ be reduced to flagd's simpler conditional system.
Once flagd can support nested conditionals like 'and' and 'or', we will provide support. 

---
## ✅ Features

- 🔍 Evaluate flags at runtime
- 🧩 Simple JSON/YAML flag format
- ♻️ Optional hot-reloading (fsnotify)
- 🤝 OpenFeature compatible
- 🌐 Future: HTTP / Redis / Consul backends
- 🔓 MIT licensed and reusable in other OSS projects

[View the Specifications here](./docs/specs.md).

---
## 🔧 Example Flag File

The simplest possible flag file is static like so:
```json
{
  "ui": {
    "variants": {
      "beta": true,
      "stable": false
    },
    "defaultVariant": "stable"
  }
}
```

To make this more dynamic, you can add rules based on an `EvalContext`:
```json
{
    "new_ui": {
      "variants": {
        "beta": true,
        "stable": false
      },
      "defaultVariant": "stable",
      "rules": [
        { "if": { "env": "prod", "group": "beta" }, "variant": "beta" },
        { "if": { "env": "prod" }, "variant": "stable" }
      ]
    }
}
```

You can also make use of YAML files, like this [example here](./examples/04-with_rules.yaml).

---
## 🧑‍💻 Usage (SDK)

```golang
store, err := featureflags.NewStoreFromFile("flags.json")
// assert no error

if store.IsEnabled("new_ui", featureflags.EvalContext{}) {
    // Enable experimental flow
}
```

---
## 📦 Use as an OpenFeature Provider

```bash
go get -u github.com/tommed/ducto-featureflags
```

[Check our OpenFeature example here](./openfeature/example_provider_test.go)

---
## 📦 Install CLI

```bash
go install github.com/tommed/ducto-featureflags/cmd/ducto-flags@latest
```

### CLI Usage

```bash
# Check if a single flag is enabled
ducto-flags -file flags.json -key new_ui -ctx env=prod -ctx region=EU

# Print all flags
ducto-flags -file flags.json -list

# Host a flags server (optional auth token)
ducto-flags serve -file flags.json [-token secret-123]
```

---
## 🛠️ Planned Backends
- [x] JSON file
- [x] YAML file
- [x] HTTP endpoint
- [x] OpenFeature compatibility
- [x] OpenFeature provider
- [ ] Redis
- [ ] Google Firestore
- [ ] Env var overrides
- [ ] Versioned flag API with auditing

---
## 🤖 Part of the Ducto Ecosystem
- [ducto-dsl](https://github.com/tommed/ducto-dsl) – declarative data transformation engine
- [ducto-orchestrator](https://github.com/tommed/ducto-orchestrator) – pluggable streaming runtime
- `ducto-featureflags` – this repo

---
## 🧰 License
- Code licensed under [MIT](./LICENSE)
- Logos and illustrations (and their likeness) are (c) Copyright 2025 Tom Medhurst, all rights reserved.
