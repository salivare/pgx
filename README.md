
# sqlx-pgx

A lightweight wrapper that integrates **sqlx** with the **pgx** driver for effortless PostgreSQL connectivity.

---

## Installation

Install the library with:

```bash
go get github.com/salivare/pgx@latest
```

---

## Configuration

By default the library reads its settings from:

```
./config/dbconfig.yaml
```

To override, export the `CONFIG_PATH` environment variable:

```bash
export CONFIG_PATH=/path/to/your/dbconfig.yaml
```

Your YAML should define your database parameters under a top-level `database` key, for example:

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: appuser
  password: secret
  name: appdb
  sslmode: disable
  pool:
    maxOpenConns: 10
    maxIdleConns: 5
  retry:
    maxAttempts: 3
    delay: 2s
    maxDelay: 10s
```

---

## Quick Start

```go
package main

import (
  "log"

  "github.com/salivare/pgx/sqlx"
)

func main() {
  db, err := sqlx.Connect()
  if err != nil {
    log.Fatalf("connection failed: %v", err)
  }
  defer db.Close()

  // Use db...
}
```

---

## Available Options

You can pass one or more `Option` functions to `Connect` to override default settings.

### WithDSN

Provide a full connection string. If omitted, the library uses your YAML/ENV config.

```go
import "github.com/salivare/pgx/sqlx"

func main() {
  custom := "postgres://u:p@myhost:5432/mydb?sslmode=disable"
  db, err := sqlx.Connect(sqlx.WithDSN(custom))
  if err != nil {
    log.Fatalf("connection failed: %v", err)
  }
  defer db.Close()

  // Use db...
}
```

### WithMaxOpenConns

Override the maximum number of open connections in the pool. If omitted, the value from config is used.

```go
import "github.com/salivare/pgx/sqlx"

func main() {
  db, err := sqlx.Connect(sqlx.WithMaxOpenConns(100))
  if err != nil {
    log.Fatalf("connection failed: %v", err)
  }
  defer db.Close()

  // Use db...
}
```

### WithRetryAttempts

Specify how many times to retry failed connection attempts. If omitted, the value from config is used.

```go
import "github.com/salivare/pgx/sqlx"

func main() {
  db, err := sqlx.Connect(sqlx.WithRetryAttempts(5))
  if err != nil {
    log.Fatalf("connection failed: %v", err)
  }
  defer db.Close()

  // Use db...
}
```

---

## Closing the Connection

Always call `defer db.Close()` immediately after a successful `Connect` to ensure all resources are released when your application terminates.

```go
db, err := sqlx.Connect()
if err != nil {
  log.Fatalf("connection failed: %v", err)
}
defer db.Close()
```
```