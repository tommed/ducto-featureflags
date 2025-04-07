<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="./assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Ducto Feature Flags

[![CI](https://github.com/tommed/ducto-featureflags/actions/workflows/ci.yml/badge.svg)](https://github.com/tommed/ducto-featureflags/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/tommed/ducto-featureflags/branch/main/graph/badge.svg)](https://codecov.io/gh/tommed/ducto-featureflags)

> Lightweight, embeddable, and pluggable feature flag engine for Go — designed for pipelines, microservices, and event-driven systems.

---
## ✨ What is Ducto-FeatureFlags?

`ducto-featureflags` is a minimalist feature flag manager built in Go. It was designed as a reusable component for:
- Data transformation pipelines (like [Ducto](https://github.com/tommed))
- Serverless functions
- Microservices and APIs
- CLI tools and test harnesses

It supports both static file-based flags and dynamic backends (coming soon), and can be used as:
- A **Go SDK**
- A **CLI for testing**
- A **preprocessor plugin** in the [ducto-orchestrator](https://github.com/tommed/ducto-orchestrator)

---
## ✅ Features

- 🔍 Evaluate flags at runtime
- 🧩 Simple JSON/YAML flag format
- ♻️ Optional hot-reloading (fsnotify)
- 🌐 Future: HTTP / Redis / Consul backends
- 🔓 MIT licensed and reusable in other OSS projects

---
## 🔧 Example Flag File

The simplest possible flag file is static like so:
```json
{
  "flags": {
    "new_ui": {
      "enabled": true
    },
    "beta_mode": {
      "enabled": false
    }
  }
}
```

To make this more dynamic, you can add rules based on an `EvalContext`:
```json
{
  "flags": {
    "new_ui": {
      "rules": [
        { "if": { "env": "prod", "group": "beta" }, "value": true },
        { "if": { "env": "prod" }, "value": false }
      ],
      "enabled": true
    }
  }
}
```

You can also make use of YAML files, like this [example here](./examples/with_rules.yaml).

---
## 🧑‍💻 Usage (SDK)

```golang
store, _ := featureflags.NewStoreFromFile("flags.json")

if store.IsEnabled("new_ui", featureflags.EvalContext{}) {
    // Enable experimental flow
}
```

---
## 📦 Install

```bash
go install github.com/tommed/ducto-featureflags/cmd/ducto-flags@latest
```

### CLI Usage

```bash
# Check if a single flag is enabled
ducto-flags -file flags.json -key new_ui -ctx env=prod -ctx region=EU

# Print all flags
ducto-flags -file flags.json -list
```

---
## 🛠️ Planned Backends
- [x] JSON file
- [x] YAML file
- [ ] HTTP endpoint
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
