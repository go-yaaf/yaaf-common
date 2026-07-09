# CLAUDE.md — yaaf-common

Guidance for AI coding assistants working in this repository. Read this before generating or reviewing code.

## What this is

`yaaf-common` is the foundational Go library for the YAAF ecosystem: interface-driven abstractions for infrastructure (relational DB, document store, cache, message bus), an entity model, and utility packages. Every interface ships with an in-memory implementation for tests. Production adapters (PostgreSQL, MongoDB, Redis, Kafka, Google Pub/Sub, BigQuery, ...) live in **separate** `go-yaaf/*` modules that implement these interfaces.

- Go 1.24.0+. Dependencies: `google/uuid`, `wneessen/go-mail` (+ `golang.org/x/text` indirect).
- Program to the interfaces (`IDatabase`, `IDatastore`, `IDataCache`, `IMessageBus`); only reference `NewInMemory*` in tests/wiring.

## Source of truth

The full, code-verified API reference is in **`llms-full.txt`** (and a concise index in `llms.txt`). When unsure about a signature, read the interface file directly instead of guessing:

- `database/database.go` — `IDatabase`
- `database/datastore.go` — `IDatastore`
- `database/datacache.go` — `IDataCache`, `ILocker`
- `database/query.go` — `IQuery`
- `database/query_filter.go` — `QueryFilter` builder (`F(field).Eq(...)`)
- `entity/entity.go` — `Entity`, `BaseEntity`, ID generators, `EntityFactory`
- `messaging/message_bus.go`, `messaging/messages.go`
- `config/base_config.go`, `logger/logger.go`

## API shapes assistants get wrong (all forms on the left do NOT compile)

| Don't write | Do write |
|-------------|----------|
| `database.NewQuery("users")` | `db.Query(NewUser)` — start from `db.Query(factory)` |
| `.WithFilter("age", database.Gte, 30)` | `.Filter(database.F("age").Gte(30))` |
| `.WithSort("name", true)` | `.Sort("name")` (asc) / `.Sort("name-")` (desc) |
| `db.Execute(query, factory)` / `db.Find(q, &out)` | `query.Find()` → `(out []Entity, total int64, err error)` |
| `db.Get(id, &user)` | `e, err := db.Get(NewUser, id)`; then `user := e.(*User)` |
| `database.NewInMemoryDataCache(5*time.Minute, ...)` | `database.NewInMemoryDataCache()` — **no args** |
| `cache.Set(key, []byte(v))` / `cache.Delete(k)` / `cache.Keys(p)` | `cache.Set(key, entity, exp...)` or `cache.SetRaw(key, bytes, exp...)`; `cache.Del(k)`; `cache.Scan(from, match, count)` |
| `cache.Increment(key, 1)` | not on `IDataCache` — use hash/list ops or a concrete adapter |
| `bus.Subscribe(topic, handler)` | `bus.Subscribe(sub, factory, callback, topics...)` |
| `messaging.NewMessage(topic, payload)` | `messaging.GetMessage[T](topic, payload)`; factory body: `messaging.NewMessage[T]()` |
| `config.NewConfig("app")` | `config.Get()` (singleton) |
| logging via `zap` | Go `log/slog`: `logger.Info(format, args...)` |

## Core rules

1. **Entities** implement `Entity` (`ID/TABLE/NAME/KEY`) — usually by embedding `BaseEntity` and implementing `TABLE()`. Always provide a factory: `func NewUser() entity.Entity { return &User{} }`.
2. **Factories everywhere.** Database/cache/messaging methods that must allocate a typed instance take an `EntityFactory` (`func() Entity`) or `MessageFactory` (`func() IMessage`).
3. **Reads return `Entity`** — type-assert to the concrete type: `user := e.(*User)`.
4. **`Timestamp`** (`entity/timestamp.go`) is `int64` epoch **milliseconds**.
5. **`IDataCache` is entity-oriented**, not a raw Redis client. Values are entities (with parallel `*Raw` byte methods). Method names: `Del` (not `Delete`), `Scan` (not `Keys`).
6. **Close resources**: `defer db.Close()`, `defer cache.Close()`, `defer bus.Close()`.
7. **Sharding**: `TABLE()` may contain `{{tenantId}}`/`{{accountId}}`/`{{year}}`/`{{month}}`; pass the shard key via the trailing `keys ...string` argument.

## Canonical snippet

```go
type User struct {
    entity.BaseEntity
    Name   string `json:"name"`
    Age    int    `json:"age"`
    Status string `json:"status"`
}
func (u *User) TABLE() string { return "user" }
func NewUser() entity.Entity  { return &User{} }

db, _ := database.NewInMemoryDatabase()

added, _ := db.Insert(&User{Name: "John", Age: 30, Status: "active"})
e, _ := db.Get(NewUser, added.ID())
user := e.(*User)

list, total, err := db.Query(NewUser).
    MatchAll(
        database.F("status").Eq("active"),
        database.F("age").Gte(30),
    ).
    Sort("name").Page(0).Limit(20).
    Find() // (out []Entity, total int64, err error)
```

## Testing

Use in-memory implementations — no external services needed. `go test ./...`. Shared sample entities live in `test/`.

## When making changes

- Match the surrounding style; keep the interface-driven, factory-based conventions.
- If you change a public interface, update `llms.txt`, `llms-full.txt`, `README.md`, and this file to match.
