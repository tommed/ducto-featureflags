# Contributing

## Rules

Please make sure you read our [Code of Conduct](./CODE_OF_CONDUCT.md) before engaging with this project.

### Rules for Developers

#### 🟣 1. Export only when absolutely necessary
> Avoid exposing structs, only expose interfaces and constructors.

#### 🟣 2. Prefer interface-oriented design to struct exposure
- Keep orchestrator, event sources, output writers defined by interfaces.
- Allows tests, mocks, fakes, and alternative backends.

#### 🟣 3. Package Layout
- `internal/` = private to repo
- `pkg/` = avoid. This is an older practice which isn't used anymore. The root directory is fine
- Avoid dumping everything into `cmd/` this is for entrypoints only and are usually not testable, so one-liners only please
- `model/` for shared data models (Program, Instruction)
- Keep tests next to the code not in a separate `tests/` package wherever possible

#### 🟣 4. Constructor Pattern
```go
func NewFoo(...) FooInterface
```

#### 🟣 5. Consistent Option Structs over long parameter lists
```go
type HTTPOptions struct {
    Addr          string
    MetadataField string
}
```

#### 🟣 6. EventSource: Return error from `Start()` immediately when applicable
> Always report errors immediately, don't bury them.

#### 🟣 7. Use context.Context everywhere for:
- Cancellation
- Deadlines
- Propagation across layers

#### 🟣 8. Avoid:
- Global state
- Hidden side effects in constructors
- `panic()` except when truly fatal (e.g., bad build tags)

#### ✅ Ducto Suggested Convention:
| Situation                               | Recommended                                 |
|-----------------------------------------|---------------------------------------------|
| DSL execution	via transform.Transformer |
| Inputs (Events)                         | via EventSource interface                   |
| Outputs                                 | via OutputWriter interface                  |
| CLI Entrypoint                          | isolated under internal/cli                 |
| Versioning                              | stick to docs/spec-vN.md and validate early |
| Debugging / Meta fields                 | reserved @keyname convention                |
| Metrics                                 | postponed until later                       |

