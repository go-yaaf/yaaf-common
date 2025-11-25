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
	// The service name ("my-app") will be used as a prefix for environment variables
	return &MyConfig{
		BaseConfig: config.NewConfig("my-app"),
	}
}

func main() {
	conf := NewMyConfig()

	// To set this value, run: export MY_APP_MY_PARAM="hello from env"
	// It will use "default-value" if the environment variable is not set.
	conf.MyParam = conf.Get("my_param", "default-value").(string)

	fmt.Println("Service Name:", conf.Service())
	fmt.Println("My Custom Param:", conf.MyParam)
	fmt.Println("HTTP Port:", conf.Port()) // Accessing a built-in variable
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

func main() {
	// Use the in-memory database for demonstration
	db, err := database.NewInMemoryDatabase()
	if err != nil {
		panic(err)
	}

	// Create a new user
	user := &User{Name: "John Doe", Age: 30}
	id, err := db.Insert(user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted user with ID: %s\n", id)

	// Get the user back
	var retrievedUser User
	if err := db.Get(id, &retrievedUser); err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved user: %+v\n", retrievedUser)

	// Query for users with age >= 30
	users := make([]*User, 0)
	q := database.NewQuery("user").WithFilter("age", database.Gte, 30)
	if err := db.Find(q, &users); err != nil {
		panic(err)
	}
	fmt.Printf("Found %d users aged 30 or more.\n", len(users))
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

// COLLECTION returns the collection name for the document
func (p *Product) COLLECTION() string { return "products" }

func main() {
	// Use the in-memory datastore for demonstration
	ds, err := database.NewInMemoryDatastore()
	if err != nil {
		panic(err)
	}

	// Insert a new product
	product := &Product{Name: "Laptop", Price: 1200.50}
	id, err := ds.Insert(product)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted product with ID: %s\n", id)

	// Get the product
	var retrievedProduct Product
	if err := ds.Get(id, &retrievedProduct); err != nil {
		panic(err)
	}
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
	// Use the in-memory cache with a 5-minute default expiration
	cache := database.NewInMemoryDataCache(5*time.Minute, 10*time.Minute)

	key := "my-special-key"
	value := "hello world"

	// Set a value in the cache
	if err := cache.Set(key, []byte(value)); err != nil {
		panic(err)
	}
	fmt.Printf("Set key '%s'\n", key)

	// Get a value from the cache
	data, err := cache.Get(key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved value: %s\n", string(data))

	// Delete the value
	if err := cache.Delete(key); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted key '%s'\n", key)

	// Trying to get the key again will result in an error
	_, err = cache.Get(key)
	if err != nil {
		fmt.Printf("Could not get key '%s': %s\n", key, err)
	}
}
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

// A simple message handler for pub/sub
func handleMessage(msg messaging.IMessage) error {
	fmt.Printf("Handler received message on topic '%s': %s\n", msg.Topic(), string(msg.Payload()))
	return nil
}

func main() {
	// Use the in-memory message bus for demonstration
	bus, err := messaging.NewInMemoryMessageBus()
	if err != nil {
		panic(err)
	}

	// 1. Publish-Subscribe Example
	topic := "my-topic"
	if _, err := bus.Subscribe(topic, handleMessage); err != nil {
		panic(err)
	}
	fmt.Println("Subscribed to topic:", topic)

	msg := messaging.NewMessage(topic, []byte("Hello Pub/Sub!"))
	if err := bus.Publish(msg); err != nil {
		panic(err)
	}
	fmt.Println("Published message to topic:", topic)

	time.Sleep(100 * time.Millisecond) // Wait for async handler

	// 2. Queueing Example
	queue := "my-queue"
	queueMsg := messaging.NewMessage(queue, []byte("Hello Queue!"))

	if err := bus.Push(queue, queueMsg); err != nil {
		panic(err)
	}
	fmt.Println("Pushed message to queue:", queue)

	// Pop message from queue (with a 1-second timeout)
	poppedMsg, err := bus.Pop(queue, 1*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Popped message from queue '%s': %s\n", poppedMsg.Topic(), string(poppedMsg.Payload()))
}
```

---

### 4. Logging (`Logger`)

`yaaf-common` includes a lightweight, structured logging wrapper around `zap`.

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