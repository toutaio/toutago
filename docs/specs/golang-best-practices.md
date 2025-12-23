# Go Programming Best Practices and Recommendations

**Version**: 1.0  
**Last Updated**: 2025-12-20  
**Status**: Active

## Overview

This document provides comprehensive best practices and recommendations for Go (Golang) programming, covering code organization, idioms, performance, testing, and common patterns.

---

## 1. Code Organization

### Project Structure

```
project/
├── cmd/                    # Main applications
│   └── myapp/
│       └── main.go
├── internal/               # Private application code
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/                    # Public library code
│   └── util/
├── api/                    # API definitions (OpenAPI, protobuf)
├── configs/                # Configuration files
├── scripts/                # Build and deployment scripts
├── test/                   # Additional test data and utilities
├── go.mod
├── go.sum
└── README.md
```

**Best Practices:**
- Use `internal/` for code that should not be imported by other projects
- Use `pkg/` for reusable libraries
- Keep `cmd/` lightweight - mainly initialization and wiring
- One `main.go` per application in separate `cmd/` subdirectories

### Package Naming

- Use short, lowercase, single-word names
- Avoid underscores, dashes, or mixed caps
- Package name should be the base name of its import path
- Avoid generic names like `util`, `common`, `base`

**Good:**
```go
package handler
package repository
package auth
```

**Avoid:**
```go
package handlerPackage
package repo_utils
package common
```

---

## 2. Code Style and Idioms

### Variable Naming

- Use camelCase for local variables
- Use PascalCase for exported identifiers
- Prefer short names for short scopes (`i`, `err`, `ok`)
- Use descriptive names for larger scopes

```go
// Good
func (s *Server) HandleRequest(ctx context.Context, req *Request) error {
    for i, item := range req.Items {
        processedItem := s.processItem(item)
        // ...
    }
}

// Avoid
func (s *Server) HandleRequest(ctx context.Context, requestObject *Request) error {
    for index, currentItem := range requestObject.Items {
        // ...
    }
}
```

### Error Handling

**Always check errors:**
```go
// Good
file, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()

// Never ignore errors
file, _ := os.Open("file.txt") // Bad!
```

**Error wrapping (Go 1.13+):**
```go
if err != nil {
    return fmt.Errorf("processing user %s: %w", userID, err)
}
```

**Custom errors:**
```go
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
)

// Check with errors.Is
if errors.Is(err, ErrNotFound) {
    // handle not found
}
```

### Defer, Panic, and Recover

**Defer for cleanup:**
```go
func processFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close() // Runs when function returns
    
    // Process file...
    return nil
}
```

**Use panic sparingly:**
- Only for unrecoverable errors in initialization
- Never panic in library code
- Use `recover()` only in specific cases (e.g., HTTP servers)

```go
// Acceptable panic usage
func init() {
    if critical_resource == nil {
        panic("critical resource not initialized")
    }
}
```

### Interface Best Practices

**Accept interfaces, return structs:**
```go
// Good
func ProcessData(r io.Reader) (*Result, error) {
    // ...
}

// Avoid
func ProcessData(f *os.File) (*Result, error) {
    // ...
}
```

**Keep interfaces small:**
```go
// Good - single responsibility
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Avoid - too many methods
type DataProcessor interface {
    Read() error
    Write() error
    Validate() error
    Transform() error
    // ...
}
```

**Define interfaces at usage point:**
```go
// In consumer package
type UserStore interface {
    GetUser(id string) (*User, error)
}

func NewService(store UserStore) *Service {
    return &Service{store: store}
}
```

---

## 3. Concurrency

### Goroutines

**Don't leak goroutines:**
```go
// Good - controlled lifecycle
func process(ctx context.Context, items []Item) {
    for _, item := range items {
        go func(i Item) {
            select {
            case <-ctx.Done():
                return
            default:
                processItem(i)
            }
        }(item)
    }
}
```

**Use WaitGroups for coordination:**
```go
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
        processItem(i)
    }(item)
}
wg.Wait()
```

### Channels

**Buffered vs unbuffered:**
```go
// Unbuffered - synchronous
ch := make(chan int)

// Buffered - asynchronous up to capacity
ch := make(chan int, 100)
```

**Close channels from sender:**
```go
// Good
func producer(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 10; i++ {
        ch <- i
    }
}

func consumer(ch <-chan int) {
    for val := range ch { // Exits when channel closed
        process(val)
    }
}
```

**Use select for multiple channels:**
```go
select {
case msg := <-ch1:
    handleMessage(msg)
case <-ctx.Done():
    return ctx.Err()
case <-time.After(5 * time.Second):
    return errors.New("timeout")
}
```

### Context Usage

**Always pass context as first parameter:**
```go
func ProcessRequest(ctx context.Context, req *Request) error {
    // ...
}
```

**Propagate context through call chain:**
```go
func (s *Service) HandleUser(ctx context.Context, id string) error {
    user, err := s.repo.GetUser(ctx, id) // Pass ctx down
    if err != nil {
        return err
    }
    return s.notifier.Notify(ctx, user) // Pass ctx down
}
```

**Use context for cancellation and timeouts:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := slowOperation(ctx)
```

---

## 4. Testing

### Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed signs", -2, 3, 1},
        {"zeros", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Test Helpers

```go
func TestUserCreation(t *testing.T) {
    user := createTestUser(t, "john@example.com")
    // ...
}

func createTestUser(t *testing.T, email string) *User {
    t.Helper() // Mark as helper for better error reporting
    user, err := NewUser(email)
    if err != nil {
        t.Fatalf("failed to create test user: %v", err)
    }
    return user
}
```

### Mocking and Interfaces

```go
// Production code
type UserRepository interface {
    GetUser(ctx context.Context, id string) (*User, error)
}

// Test code
type mockUserRepo struct {
    users map[string]*User
}

func (m *mockUserRepo) GetUser(ctx context.Context, id string) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}
```

### Benchmarks

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

// Run: go test -bench=. -benchmem
```

---

## 5. Performance Best Practices

### Memory Allocation

**Reuse slices with capacity:**
```go
// Good
results := make([]Result, 0, len(inputs))
for _, input := range inputs {
    results = append(results, process(input))
}

// Avoid - causes reallocations
var results []Result
for _, input := range inputs {
    results = append(results, process(input))
}
```

**Use sync.Pool for temporary objects:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process() {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    buf.Reset()
    // Use buffer...
}
```

### String Building

```go
// Good - for multiple concatenations
var sb strings.Builder
sb.WriteString("Hello")
sb.WriteString(" ")
sb.WriteString("World")
result := sb.String()

// Avoid - creates intermediate strings
result := "Hello" + " " + "World"
```

### Map Preallocation

```go
// Good
m := make(map[string]int, expectedSize)

// Less optimal
m := make(map[string]int)
```

---

## 6. Common Patterns

### Functional Options

```go
type Server struct {
    host string
    port int
    timeout time.Duration
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        host: "localhost",
        port: 8080,
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(
    WithHost("0.0.0.0"),
    WithPort(9000),
)
```

### Builder Pattern

```go
type QueryBuilder struct {
    query  strings.Builder
    params []interface{}
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
    qb.query.WriteString("SELECT ")
    qb.query.WriteString(strings.Join(fields, ", "))
    return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
    qb.query.WriteString(" FROM ")
    qb.query.WriteString(table)
    return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
    return qb.query.String(), qb.params
}

// Usage
query, params := NewQueryBuilder().
    Select("id", "name").
    From("users").
    Build()
```

### Singleton (with sync.Once)

```go
type Config struct {
    // config fields
}

var (
    instance *Config
    once     sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        instance = &Config{
            // initialize
        }
    })
    return instance
}
```

---

## 7. Dependencies and Modules

### Go Modules Best Practices

```bash
# Initialize module
go mod init github.com/username/project

# Add dependencies
go get package@version

# Update dependencies
go get -u ./...

# Tidy up
go mod tidy

# Vendor dependencies
go mod vendor
```

**Version pinning:**
```go
// go.mod
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4
)
```

### Dependency Injection

```go
// Good - explicit dependencies
type Service struct {
    repo   UserRepository
    logger Logger
    cache  Cache
}

func NewService(repo UserRepository, logger Logger, cache Cache) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
        cache:  cache,
    }
}

// Avoid - hidden dependencies
var globalDB *sql.DB

func (s *Service) GetUser(id string) (*User, error) {
    return globalDB.Query(...) // Hidden dependency
}
```

---

## 8. Security Best Practices

### Input Validation

```go
func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("email cannot be empty")
    }
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
    }
    if len(email) > 254 {
        return errors.New("email too long")
    }
    return nil
}
```

### SQL Injection Prevention

```go
// Good - parameterized query
func GetUser(db *sql.DB, id string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    // ...
}

// Bad - string concatenation
func GetUser(db *sql.DB, id string) (*User, error) {
    query := "SELECT * FROM users WHERE id = " + id // DON'T DO THIS!
    row := db.QueryRow(query)
    // ...
}
```

### Secrets Management

```go
// Good - environment variables
dbPassword := os.Getenv("DB_PASSWORD")
if dbPassword == "" {
    return errors.New("DB_PASSWORD not set")
}

// Bad - hardcoded
const dbPassword = "super_secret_123" // NEVER DO THIS
```

---

## 9. Logging and Observability

### Structured Logging

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

logger.Info("user created",
    "user_id", userID,
    "email", email,
    "timestamp", time.Now(),
)

logger.Error("failed to process request",
    "error", err,
    "request_id", reqID,
)
```

### Error Context

```go
func processUser(id string) error {
    user, err := getUser(id)
    if err != nil {
        return fmt.Errorf("processUser %s: %w", id, err)
    }
    
    if err := validateUser(user); err != nil {
        return fmt.Errorf("processUser %s validation: %w", id, err)
    }
    
    return nil
}
```

---

## 10. Documentation

### Package Documentation

```go
// Package auth provides authentication and authorization utilities.
//
// This package implements JWT-based authentication and role-based
// access control for HTTP services.
package auth
```

### Function Documentation

```go
// Authenticate validates the provided credentials and returns a JWT token.
//
// The token is valid for 24 hours and includes the user's ID and roles.
// Returns ErrInvalidCredentials if authentication fails.
//
// Example:
//
//	token, err := Authenticate("user@example.com", "password123")
//	if err != nil {
//	    log.Fatal(err)
//	}
func Authenticate(email, password string) (string, error) {
    // ...
}
```

### Example Tests

```go
func ExampleAuthenticate() {
    token, err := Authenticate("user@example.com", "password")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Token received:", len(token) > 0)
    // Output: Token received: true
}
```

---

## 11. Tools and Commands

### Essential Tools

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Static analysis
go vet ./...

# Generate code coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Profile memory
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Recommended Linters

```yaml
# .golangci.yml
linters:
  enable:
    - errcheck      # Check error handling
    - gosimple      # Simplify code
    - govet         # Report suspicious constructs
    - ineffassign   # Detect ineffectual assignments
    - staticcheck   # Advanced static analysis
    - unused        # Find unused code
    - gofmt         # Check formatting
    - goimports     # Check import formatting
```

---

## 12. Common Pitfalls to Avoid

### Loop Variable Capture

```go
// Wrong
for _, item := range items {
    go func() {
        process(item) // All goroutines see last item!
    }()
}

// Correct
for _, item := range items {
    item := item // Create new variable
    go func() {
        process(item)
    }()
}

// Or pass as parameter
for _, item := range items {
    go func(i Item) {
        process(i)
    }(item)
}
```

### Nil Interface Gotcha

```go
type MyError struct{}
func (e *MyError) Error() string { return "error" }

func returnsError() error {
    var err *MyError // nil pointer
    return err       // Returns non-nil interface!
}

// Fix
func returnsError() error {
    var err *MyError
    if err == nil {
        return nil // Return nil explicitly
    }
    return err
}
```

### Shadowing Variables

```go
// Wrong
var err error
data, err := fetchData()
if err != nil {
    // This shadows err in inner scope
    err := handleError(err)
    if err != nil {
        return err
    }
}
// Original err is unchanged here!

// Correct
var err error
data, err := fetchData()
if err != nil {
    err = handleError(err) // Use = not :=
    if err != nil {
        return err
    }
}
```

---

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

---

**Note**: These are guidelines, not strict rules. Always consider your specific context and team conventions. When in doubt, prioritize clarity and maintainability over cleverness.
