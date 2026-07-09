# yaaf-common

[![Build](https://github.com/go-yaaf/yaaf-common/actions/workflows/build.yml/badge.svg)](https://github.com/go-yaaf/yaaf-common/actions/workflows/build.yml)

**YAAF-Common** is the foundational library for the **YAAF (Yet Another Application Framework)** ecosystem in Go. It provides a suite of common interfaces and utility packages designed to abstract away the boilerplate of building modern microservices.

## Core Philosophy

The primary goal of `yaaf-common` is to promote clean architecture and loose coupling by defining a set of standard interfaces for common infrastructure components like databases, caches, and message brokers.

By coding against these interfaces, your application's business logic remains decoupled from specific implementations. This allows you to:
- Easily swap infrastructure components without changing application code.
- Utilize the provided in-memory implementations for fast, dependency-free unit and integration testing.
- Maintain a consistent and predictable structure across all your microservices.

## Installation

To add `yaaf-common` to your project, use `go get`:
```bash
go get -u github.com/go-yaaf/yaaf-common
```

## Developer Guide

This guide provides an overview of the core components and how to use them.

### 1. Configuration (`BaseConfig`)

`BaseConfig` provides a simple, environment-variable-driven approach to application configuration. It's designed to work seamlessly with containerized environments like Docker and Kubernetes.

You should embed `BaseConfig` in your application's specific configuration struct. The library automatically maps environment variables in the format `SERVICE_NAME_VARIABLE_NAME`.

**Example:**
```go
package main

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/config"
)

// Define your application's configuration
type MyConfig struct {
	*config.BaseConfig
	MyParam string `json:"my_param"`
}

// NewMyConfig creates a new application configuration
func NewMyConfig() *MyConfig {
	// config.Get() returns the process-wide singleton (there is no NewConfig()).
	// It scans environment variables on first access.
	return &MyConfig{
		BaseConfig: config.Get(),
	}
}

func main() {
	conf := NewMyConfig()

	// To set this value, run: export MY_PARAM="hello from env"
	// It falls back to "default-value" if the environment variable is not set.
	conf.MyParam = conf.GetStringParamValueOrDefault("MY_PARAM", "default-value")

	fmt.Println("My Custom Param:", conf.MyParam)
	fmt.Println("Database URI:", conf.DatabaseUri()) // a built-in named accessor
	fmt.Println("Log Level:", conf.LogLevel())
}
```

---

### 2. Database & Cache Interfaces

`yaaf-common` provides interfaces for different data storage patterns. For testing purposes, it includes in-memory implementations for each.

#### Relational Databases (`IDatabase`)

The `IDatabase` interface abstracts standard CRUD and query operations for SQL-like relational databases.

**Example using `InMemoryDatabase`:**
```go
package main

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/database"
	"github.com/go-yaaf/yaaf-common/entity"
)

// User defines a sample entity
type User struct {
	entity.BaseEntity
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TABLE returns the database table name for the entity
func (u *User) TABLE() string { return "user" }

// NewUser is the entity factory used by the database layer to allocate typed instances
func NewUser() entity.Entity { return &User{} }

func main() {
	// Use the in-memory database for demonstration
	db, err := database.NewInMemoryDatabase()
	if err != nil {
		panic(err)
	}

	// Create a new user (Insert returns the stored entity)
	added, err := db.Insert(&User{Name: "John Doe", Age: 30})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted user with ID: %s\n", added.ID())

	// Get the user back — reads return an entity.Entity; type-assert to the concrete type
	e, err := db.Get(NewUser, added.ID())
	if err != nil {
		panic(err)
	}
	retrievedUser := e.(*User)
	fmt.Printf("Retrieved user: %+v\n", retrievedUser)

	// Query for users with age >= 30. Build the query from db.Query(factory),
	// then execute it with Find(), which returns (list, total, error).
	list, total, err := db.Query(NewUser).
		Filter(database.F("age").Gte(30)).
		Sort("name"). // ascending; use "name-" for descending
		Find()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d users aged 30 or more.\n", total)
	for _, item := range list {
		fmt.Printf("  - %s\n", item.(*User).Name)
	}
}
```

#### Document Stores (`IDatastore`)

The `IDatastore` interface is designed for NoSQL document-oriented databases.

**Example using `InMemoryDatastore`:**
```go
package main

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/database"
	"github.com/go-yaaf/yaaf-common/entity"
)

// Product defines a sample document
type Product struct {
	entity.BaseEntity
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// TABLE returns the collection name for the document (IDatastore reuses the Entity contract)
func (p *Product) TABLE() string { return "products" }

// NewProduct is the entity factory
func NewProduct() entity.Entity { return &Product{} }

func main() {
	// Use the in-memory datastore for demonstration
	ds, err := database.NewInMemoryDatastore()
	if err != nil {
		panic(err)
	}

	// Insert a new product (returns the stored entity)
	added, err := ds.Insert(&Product{Name: "Laptop", Price: 1200.50})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted product with ID: %s\n", added.ID())

	// Get the product back and type-assert to the concrete type
	e, err := ds.Get(NewProduct, added.ID())
	if err != nil {
		panic(err)
	}
	retrievedProduct := e.(*Product)
	fmt.Printf("Retrieved product: %+v\n", retrievedProduct)
}
```

#### Caching (`IDataCache`)

The `IDataCache` interface abstracts key-value cache operations, ideal for use with systems like Redis.

**Example using `InMemoryDataCache`:**
```go
package main

import (
	"fmt"
	"time"
	"github.com/go-yaaf/yaaf-common/database"
)

func main() {
	// Use the in-memory cache (the constructor takes no arguments)
	cache, err := database.NewInMemoryDataCache()
	if err != nil {
		panic(err)
	}

	key := "my-special-key"
	value := "hello world"

	// IDataCache is entity-oriented; use the *Raw methods to store raw bytes.
	// The optional trailing argument sets a per-key expiration.
	if err := cache.SetRaw(key, []byte(value), 5*time.Minute); err != nil {
		panic(err)
	}
	fmt.Printf("Set key '%s'\n", key)

	// Get the raw value back
	data, err := cache.GetRaw(key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved value: %s\n", string(data))

	// Delete the value (note: the method is Del, not Delete)
	if err := cache.Del(key); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted key '%s'\n", key)

	// Check existence
	if exists, _ := cache.Exists(key); !exists {
		fmt.Printf("Key '%s' no longer exists\n", key)
	}
}

// To cache full entities instead of raw bytes, use the entity-oriented methods:
//   cache.Set("user:"+u.ID(), u, 5*time.Minute)
//   e, err := cache.Get(NewUser, "user:"+u.ID()) // returns entity.Entity
```

---

### 3. Messaging (`IMessageBus`)

The `IMessageBus` interface provides an abstraction for both publish-subscribe and queueing messaging patterns.

**Example using `InMemoryMessageBus`:**
```go
package main

import (
	"fmt"
	"time"
	"github.com/go-yaaf/yaaf-common/messaging"
)

// messageFactory tells the bus how to allocate an incoming message for deserialization.
// Here messages carry a string payload.
func messageFactory() messaging.IMessage { return messaging.NewMessage[string]() }

// A subscription callback returns true to acknowledge the message.
func handleMessage(msg messaging.IMessage) bool {
	fmt.Printf("Handler received message on topic '%s': %v\n", msg.Topic(), msg.Payload())
	return true
}

func main() {
	// Use the in-memory message bus for demonstration
	bus, err := messaging.NewInMemoryMessageBus()
	if err != nil {
		panic(err)
	}

	// 1. Publish-Subscribe Example
	// Subscribe(subscriptionName, messageFactory, callback, topics...)
	topic := "my-topic"
	sub, err := bus.Subscribe("demo-subscription", messageFactory, handleMessage, topic)
	if err != nil {
		panic(err)
	}
	fmt.Println("Subscribed to topic:", topic)

	// Build a typed message with GetMessage[T](topic, payload) and publish it (variadic).
	if err := bus.Publish(messaging.GetMessage[string](topic, "Hello Pub/Sub!")); err != nil {
		panic(err)
	}
	fmt.Println("Published message to topic:", topic)

	time.Sleep(100 * time.Millisecond) // Wait for async handler
	bus.Unsubscribe(sub)

	// 2. Queueing Example
	queue := "my-queue"
	if err := bus.Push(messaging.GetMessage[string](queue, "Hello Queue!")); err != nil {
		panic(err)
	}
	fmt.Println("Pushed message to queue:", queue)

	// Pop a message from the queue: Pop(messageFactory, timeout, queue...)
	poppedMsg, err := bus.Pop(messageFactory, 1*time.Second, queue)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Popped message from queue '%s': %v\n", poppedMsg.Topic(), poppedMsg.Payload())
}
```

---

### 4. Logging (`Logger`)

`yaaf-common` includes a lightweight, structured logging wrapper around the standard library's `log/slog`.

**Example:**
```go
package main

import (
	"github.com/go-yaaf/yaaf-common/logger"
)

func main() {
	// Initialize the logger for a development environment
	logger.SetLevel("DEBUG")
	logger.EnableJsonFormat(true)
	logger.EnableStacktrace(false)
	logger.Init()

	logger.Info("Server is starting...")

	// Use structured logging with fields for context
	logger.Warn("Configuration value %s is missing, using default.", "DATABASE_URI")

	logger.Error("Failed to connect to database: %s", "postgres://user:pwd@localhost:5432/db")
}
```

---

### 5. Utilities

The `utils` package contains a collection of helpers for common tasks like:
- **`collections`**: Thread-safe maps, sets, queues, etc.
- **`pool`**: Worker pools for managing concurrent tasks.
- **`binary`**: Helpers for working with binary data.
- And more for hashing, HTTP, and time manipulation.


## License

This project is licensed under the [Apache 2.0 License](LICENSE).