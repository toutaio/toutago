# Test-Driven Development (TDD) and Behavior-Driven Development (BDD) in Go

**Version**: 1.0  
**Last Updated**: 2025-12-20  
**Status**: Active

## Overview

This document provides comprehensive guidance on implementing Test-Driven Development (TDD) and Behavior-Driven Development (BDD) practices in Go projects. It covers methodologies, tools, patterns, and real-world examples.

---

## Table of Contents

1. [TDD Fundamentals](#1-tdd-fundamentals)
2. [TDD in Go](#2-tdd-in-go)
3. [BDD Fundamentals](#3-bdd-fundamentals)
4. [BDD in Go](#4-bdd-in-go)
5. [Testing Tools and Frameworks](#5-testing-tools-and-frameworks)
6. [Test Patterns](#6-test-patterns)
7. [Mocking and Test Doubles](#7-mocking-and-test-doubles)
8. [Integration and E2E Testing](#8-integration-and-e2e-testing)
9. [Code Coverage](#9-code-coverage)
10. [Best Practices](#10-best-practices)

---

## 1. TDD Fundamentals

### The TDD Cycle (Red-Green-Refactor)

```
1. RED    → Write a failing test
2. GREEN  → Write minimal code to pass
3. REFACTOR → Improve code quality
```

### TDD Principles

1. **Write tests first** - Before implementation
2. **Small steps** - One test at a time
3. **Immediate feedback** - Run tests frequently
4. **Design through tests** - Tests drive API design
5. **Refactor confidently** - Tests provide safety net

### Benefits

- **Better design** - Forces thinking about interfaces
- **Living documentation** - Tests describe behavior
- **Regression prevention** - Catch bugs early
- **Refactoring confidence** - Safe to improve code
- **Reduced debugging** - Issues caught immediately

---

## 2. TDD in Go

### Basic TDD Workflow

#### Step 1: Write a Failing Test (RED)

```go
// calculator_test.go
package calculator_test

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}
```

**Run and see it fail:**
```bash
$ go test
# calculator
./calculator_test.go:6:13: undefined: Add
FAIL    calculator [build failed]
```

#### Step 2: Write Minimal Code (GREEN)

```go
// calculator.go
package calculator

func Add(a, b int) int {
    return a + b
}
```

**Run and see it pass:**
```bash
$ go test
PASS
ok      calculator      0.001s
```

#### Step 3: Refactor

```go
// No refactoring needed for this simple case
// But if code was duplicated or unclear, improve it now
```

### Table-Driven TDD

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"zero", 0, 0, 0},
        {"mixed", -5, 10, 5},
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

### TDD Example: Building a User Service

#### Iteration 1: Create User

**Test (RED):**
```go
// user_service_test.go
package user_test

import (
    "testing"
    "myapp/user"
)

func TestCreateUser(t *testing.T) {
    service := user.NewService()
    
    u, err := service.CreateUser("alice@example.com")
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if u.Email != "alice@example.com" {
        t.Errorf("email = %s; want alice@example.com", u.Email)
    }
    if u.ID == "" {
        t.Error("expected non-empty ID")
    }
}
```

**Implementation (GREEN):**
```go
// user_service.go
package user

import "github.com/google/uuid"

type User struct {
    ID    string
    Email string
}

type Service struct{}

func NewService() *Service {
    return &Service{}
}

func (s *Service) CreateUser(email string) (*User, error) {
    return &User{
        ID:    uuid.New().String(),
        Email: email,
    }, nil
}
```

#### Iteration 2: Validate Email

**Test (RED):**
```go
func TestCreateUser_InvalidEmail(t *testing.T) {
    service := user.NewService()
    
    _, err := service.CreateUser("invalid-email")
    
    if err == nil {
        t.Fatal("expected error for invalid email")
    }
    if err.Error() != "invalid email format" {
        t.Errorf("error = %v; want 'invalid email format'", err)
    }
}
```

**Implementation (GREEN):**
```go
import (
    "errors"
    "regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (s *Service) CreateUser(email string) (*User, error) {
    if !emailRegex.MatchString(email) {
        return nil, errors.New("invalid email format")
    }
    
    return &User{
        ID:    uuid.New().String(),
        Email: email,
    }, nil
}
```

#### Iteration 3: Prevent Duplicate Emails

**Test (RED):**
```go
func TestCreateUser_DuplicateEmail(t *testing.T) {
    service := user.NewService()
    
    _, err := service.CreateUser("alice@example.com")
    if err != nil {
        t.Fatalf("first create failed: %v", err)
    }
    
    _, err = service.CreateUser("alice@example.com")
    if err == nil {
        t.Fatal("expected error for duplicate email")
    }
    if !errors.Is(err, user.ErrDuplicateEmail) {
        t.Errorf("error = %v; want ErrDuplicateEmail", err)
    }
}
```

**Implementation (GREEN):**
```go
var ErrDuplicateEmail = errors.New("email already exists")

type Service struct {
    users map[string]*User
}

func NewService() *Service {
    return &Service{
        users: make(map[string]*User),
    }
}

func (s *Service) CreateUser(email string) (*User, error) {
    if !emailRegex.MatchString(email) {
        return nil, errors.New("invalid email format")
    }
    
    if _, exists := s.users[email]; exists {
        return nil, ErrDuplicateEmail
    }
    
    u := &User{
        ID:    uuid.New().String(),
        Email: email,
    }
    s.users[email] = u
    
    return u, nil
}
```

---

## 3. BDD Fundamentals

### What is BDD?

BDD extends TDD by:
- Using natural language (Given-When-Then)
- Focusing on behavior, not implementation
- Collaborating with non-technical stakeholders
- Writing executable specifications

### Given-When-Then Format

```
GIVEN (context/preconditions)
WHEN  (action/event)
THEN  (expected outcome)
```

### BDD Benefits

- **Shared understanding** - Business and developers aligned
- **Clear specifications** - Behavior described in plain language
- **Living documentation** - Specs stay current with code
- **Focus on value** - Tests describe user-facing behavior

---

## 4. BDD in Go

### BDD with Standard Testing

```go
func TestUserLogin(t *testing.T) {
    t.Run("successful login with valid credentials", func(t *testing.T) {
        // GIVEN a registered user
        service := setupUserService(t)
        email := "alice@example.com"
        password := "secure123"
        service.Register(email, password)
        
        // WHEN they login with correct credentials
        token, err := service.Login(email, password)
        
        // THEN they receive an auth token
        if err != nil {
            t.Fatalf("expected successful login: %v", err)
        }
        if token == "" {
            t.Error("expected non-empty token")
        }
    })
    
    t.Run("login fails with wrong password", func(t *testing.T) {
        // GIVEN a registered user
        service := setupUserService(t)
        email := "alice@example.com"
        service.Register(email, "correct_password")
        
        // WHEN they login with wrong password
        _, err := service.Login(email, "wrong_password")
        
        // THEN they receive an authentication error
        if err == nil {
            t.Fatal("expected authentication error")
        }
        if !errors.Is(err, auth.ErrInvalidCredentials) {
            t.Errorf("error = %v; want ErrInvalidCredentials", err)
        }
    })
}
```

### BDD with Ginkgo and Gomega

**Installation:**
```bash
go get github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega
```

**Initialize:**
```bash
cd mypackage
ginkgo bootstrap
ginkgo generate user_service
```

**BDD Style Test:**
```go
// user_service_test.go
package user_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "myapp/user"
)

var _ = Describe("UserService", func() {
    var service *user.Service
    
    BeforeEach(func() {
        service = user.NewService()
    })
    
    Describe("Creating a user", func() {
        Context("with valid email", func() {
            It("creates user successfully", func() {
                u, err := service.CreateUser("alice@example.com")
                
                Expect(err).NotTo(HaveOccurred())
                Expect(u.Email).To(Equal("alice@example.com"))
                Expect(u.ID).NotTo(BeEmpty())
            })
        })
        
        Context("with invalid email", func() {
            It("returns validation error", func() {
                _, err := service.CreateUser("invalid")
                
                Expect(err).To(MatchError("invalid email format"))
            })
        })
        
        Context("with duplicate email", func() {
            It("returns duplicate error", func() {
                service.CreateUser("alice@example.com")
                
                _, err := service.CreateUser("alice@example.com")
                
                Expect(err).To(Equal(user.ErrDuplicateEmail))
            })
        })
    })
    
    Describe("Authenticating a user", func() {
        BeforeEach(func() {
            service.Register("alice@example.com", "password123")
        })
        
        Context("with correct credentials", func() {
            It("returns auth token", func() {
                token, err := service.Login("alice@example.com", "password123")
                
                Expect(err).NotTo(HaveOccurred())
                Expect(token).NotTo(BeEmpty())
            })
        })
        
        Context("with wrong password", func() {
            It("returns authentication error", func() {
                _, err := service.Login("alice@example.com", "wrong")
                
                Expect(err).To(Equal(user.ErrInvalidCredentials))
            })
        })
    })
})
```

**Run Ginkgo tests:**
```bash
ginkgo -v
```

### BDD with Godog (Cucumber for Go)

**Installation:**
```bash
go get github.com/cucumber/godog/cmd/godog@latest
```

**Feature file:**
```gherkin
# features/user_login.feature
Feature: User Authentication
  As a user
  I want to login to my account
  So that I can access protected resources

  Scenario: Successful login with valid credentials
    Given I am a registered user with email "alice@example.com"
    And my password is "secure123"
    When I attempt to login
    Then I should receive an authentication token
    And the login should succeed

  Scenario: Failed login with invalid password
    Given I am a registered user with email "alice@example.com"
    And my password is "secure123"
    When I attempt to login with password "wrong_password"
    Then I should receive an authentication error
    And the login should fail

  Scenario: Failed login with non-existent user
    Given I am not a registered user
    When I attempt to login with email "unknown@example.com"
    Then I should receive a user not found error
```

**Step definitions:**
```go
// features/user_login_test.go
package features

import (
    "context"
    "errors"
    "github.com/cucumber/godog"
    "myapp/auth"
)

type authContext struct {
    service  *auth.Service
    email    string
    password string
    token    string
    err      error
}

func (ctx *authContext) iAmARegisteredUserWithEmail(email string) error {
    ctx.email = email
    return ctx.service.Register(email, ctx.password)
}

func (ctx *authContext) myPasswordIs(password string) error {
    ctx.password = password
    return nil
}

func (ctx *authContext) iAttemptToLogin() error {
    ctx.token, ctx.err = ctx.service.Login(ctx.email, ctx.password)
    return nil
}

func (ctx *authContext) iAttemptToLoginWithPassword(password string) error {
    ctx.token, ctx.err = ctx.service.Login(ctx.email, password)
    return nil
}

func (ctx *authContext) iShouldReceiveAnAuthenticationToken() error {
    if ctx.token == "" {
        return errors.New("expected authentication token")
    }
    return nil
}

func (ctx *authContext) theLoginShouldSucceed() error {
    if ctx.err != nil {
        return ctx.err
    }
    return nil
}

func (ctx *authContext) iShouldReceiveAnAuthenticationError() error {
    if !errors.Is(ctx.err, auth.ErrInvalidCredentials) {
        return errors.New("expected authentication error")
    }
    return nil
}

func (ctx *authContext) theLoginShouldFail() error {
    if ctx.err == nil {
        return errors.New("expected login to fail")
    }
    return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    ac := &authContext{
        service: auth.NewService(),
    }
    
    ctx.Step(`^I am a registered user with email "([^"]*)"$`, ac.iAmARegisteredUserWithEmail)
    ctx.Step(`^my password is "([^"]*)"$`, ac.myPasswordIs)
    ctx.Step(`^I attempt to login$`, ac.iAttemptToLogin)
    ctx.Step(`^I attempt to login with password "([^"]*)"$`, ac.iAttemptToLoginWithPassword)
    ctx.Step(`^I should receive an authentication token$`, ac.iShouldReceiveAnAuthenticationToken)
    ctx.Step(`^the login should succeed$`, ac.theLoginShouldSucceed)
    ctx.Step(`^I should receive an authentication error$`, ac.iShouldReceiveAnAuthenticationError)
    ctx.Step(`^the login should fail$`, ac.theLoginShouldFail)
}
```

**Run Godog tests:**
```bash
godog run features/
```

---

## 5. Testing Tools and Frameworks

### Standard Library (`testing`)

**Pros:**
- Built-in, no dependencies
- Simple and straightforward
- Good for unit tests

**Cons:**
- Verbose assertions
- Limited BDD support

### Testify

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
)

func TestUser(t *testing.T) {
    user := CreateUser("alice@example.com")
    
    // Assertions that continue on failure
    assert.NotNil(t, user)
    assert.Equal(t, "alice@example.com", user.Email)
    
    // Assertions that stop test on failure
    require.NotNil(t, user)
    require.Equal(t, "alice@example.com", user.Email)
}
```

### Ginkgo + Gomega

**Best for:**
- BDD-style testing
- Complex test suites
- Teams familiar with RSpec/Jasmine

```go
Describe("Calculator", func() {
    It("adds numbers", func() {
        Expect(Add(2, 3)).To(Equal(5))
    })
    
    It("handles negatives", func() {
        Expect(Add(-2, -3)).To(Equal(-5))
    })
})
```

### Godog

**Best for:**
- Collaboration with non-developers
- Executable specifications
- Acceptance testing

### GoConvey

```go
import (
    . "github.com/smartystreets/goconvey/convey"
    "testing"
)

func TestCalculator(t *testing.T) {
    Convey("Given two numbers", t, func() {
        a, b := 2, 3
        
        Convey("When adding them", func() {
            result := Add(a, b)
            
            Convey("The result should be their sum", func() {
                So(result, ShouldEqual, 5)
            })
        })
    })
}
```

---

## 6. Test Patterns

### Setup and Teardown

```go
func TestMain(m *testing.M) {
    // Global setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Global teardown
    teardown()
    
    os.Exit(code)
}

func TestUser(t *testing.T) {
    // Per-test setup
    cleanup := setupTest(t)
    defer cleanup()
    
    // Test code
}

func setupTest(t *testing.T) func() {
    // Setup
    db := setupDatabase(t)
    
    // Return cleanup function
    return func() {
        db.Close()
    }
}
```

### Test Fixtures

```go
// testdata/users.json
[
    {"id": "1", "email": "alice@example.com"},
    {"id": "2", "email": "bob@example.com"}
]

// Test
func TestLoadUsers(t *testing.T) {
    data, err := os.ReadFile("testdata/users.json")
    require.NoError(t, err)
    
    var users []User
    err = json.Unmarshal(data, &users)
    require.NoError(t, err)
    
    assert.Len(t, users, 2)
}
```

### Helper Functions

```go
func createTestUser(t *testing.T, email string) *User {
    t.Helper()
    
    user, err := NewUser(email)
    if err != nil {
        t.Fatalf("failed to create test user: %v", err)
    }
    return user
}

func TestUserOperations(t *testing.T) {
    user := createTestUser(t, "alice@example.com")
    // Use user in test
}
```

### Subtests

```go
func TestUserValidation(t *testing.T) {
    testCases := []struct {
        name  string
        email string
        valid bool
    }{
        {"valid email", "alice@example.com", true},
        {"missing @", "aliceexample.com", false},
        {"missing domain", "alice@", false},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            err := ValidateEmail(tc.email)
            if tc.valid && err != nil {
                t.Errorf("expected valid email, got error: %v", err)
            }
            if !tc.valid && err == nil {
                t.Error("expected validation error")
            }
        })
    }
}
```

### Parallel Tests

```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name string
        data int
    }{
        {"test1", 1},
        {"test2", 2},
        {"test3", 3},
    }
    
    for _, tt := range tests {
        tt := tt // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Run in parallel
            // Test code
        })
    }
}
```

---

## 7. Mocking and Test Doubles

### Interface-Based Mocking

```go
// Production code
type UserRepository interface {
    Save(user *User) error
    FindByEmail(email string) (*User, error)
}

type UserService struct {
    repo UserRepository
}

// Test double
type mockUserRepo struct {
    users map[string]*User
    err   error
}

func (m *mockUserRepo) Save(user *User) error {
    if m.err != nil {
        return m.err
    }
    m.users[user.Email] = user
    return nil
}

func (m *mockUserRepo) FindByEmail(email string) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    user, ok := m.users[email]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

// Test
func TestUserService_CreateUser(t *testing.T) {
    repo := &mockUserRepo{
        users: make(map[string]*User),
    }
    service := NewUserService(repo)
    
    err := service.CreateUser("alice@example.com")
    
    assert.NoError(t, err)
    assert.Len(t, repo.users, 1)
}
```

### Using testify/mock

```go
import "github.com/stretchr/testify/mock"

type MockUserRepo struct {
    mock.Mock
}

func (m *MockUserRepo) Save(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Test
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepo)
    mockRepo.On("FindByEmail", "alice@example.com").Return(nil, ErrNotFound)
    mockRepo.On("Save", mock.AnythingOfType("*User")).Return(nil)
    
    service := NewUserService(mockRepo)
    err := service.CreateUser("alice@example.com")
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Using GoMock

```bash
go install github.com/golang/mock/mockgen@latest
mockgen -source=user_repository.go -destination=mocks/user_repository_mock.go
```

```go
import (
    "github.com/golang/mock/gomock"
    "myapp/mocks"
)

func TestUserService(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockRepo.EXPECT().
        FindByEmail("alice@example.com").
        Return(nil, ErrNotFound)
    mockRepo.EXPECT().
        Save(gomock.Any()).
        Return(nil)
    
    service := NewUserService(mockRepo)
    err := service.CreateUser("alice@example.com")
    
    assert.NoError(t, err)
}
```

### HTTP Mocking with httptest

```go
func TestHTTPClient(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "/users", r.URL.Path)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode([]User{{Email: "alice@example.com"}})
    }))
    defer server.Close()
    
    client := NewClient(server.URL)
    users, err := client.GetUsers()
    
    assert.NoError(t, err)
    assert.Len(t, users, 1)
}
```

---

## 8. Integration and E2E Testing

### Database Integration Tests

```go
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewUserRepository(db)
    
    // Test
    user := &User{Email: "alice@example.com"}
    err := repo.Save(user)
    require.NoError(t, err)
    
    found, err := repo.FindByEmail("alice@example.com")
    require.NoError(t, err)
    assert.Equal(t, user.Email, found.Email)
}

func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    
    db, err := sql.Open("postgres", "postgres://localhost/test?sslmode=disable")
    require.NoError(t, err)
    
    // Run migrations
    runMigrations(t, db)
    
    t.Cleanup(func() {
        db.Exec("DROP TABLE users")
        db.Close()
    })
    
    return db
}
```

**Run integration tests:**
```bash
# Skip integration tests
go test -short ./...

# Run only integration tests
go test -run Integration ./...
```

### Docker-based Integration Tests

```go
import "github.com/testcontainers/testcontainers-go"

func TestWithPostgres(t *testing.T) {
    ctx := context.Background()
    
    postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image:        "postgres:15",
            ExposedPorts: []string{"5432/tcp"},
            Env: map[string]string{
                "POSTGRES_PASSWORD": "test",
                "POSTGRES_DB":       "testdb",
            },
            WaitingFor: wait.ForLog("database system is ready"),
        },
        Started: true,
    })
    require.NoError(t, err)
    defer postgres.Terminate(ctx)
    
    // Get connection details
    host, _ := postgres.Host(ctx)
    port, _ := postgres.MappedPort(ctx, "5432")
    
    // Connect and test
    db := connectDB(t, host, port.Port())
    // Run tests...
}
```

### API End-to-End Tests

```go
func TestAPI_E2E(t *testing.T) {
    // Start test server
    server := setupTestServer(t)
    defer server.Close()
    
    t.Run("create and retrieve user", func(t *testing.T) {
        // Create user
        createResp := createUser(t, server.URL, "alice@example.com")
        assert.Equal(t, http.StatusCreated, createResp.StatusCode)
        
        var user User
        json.NewDecoder(createResp.Body).Decode(&user)
        
        // Retrieve user
        getResp := getUser(t, server.URL, user.ID)
        assert.Equal(t, http.StatusOK, getResp.StatusCode)
        
        var retrieved User
        json.NewDecoder(getResp.Body).Decode(&retrieved)
        assert.Equal(t, user.Email, retrieved.Email)
    })
}
```

---

## 9. Code Coverage

### Generate Coverage Report

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage percentage
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Coverage in CI/CD

```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests with coverage
        run: go test -v -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

### Coverage Best Practices

- **Aim for 80%+** - Good coverage baseline
- **Don't obsess over 100%** - Some code is hard to test
- **Focus on critical paths** - Business logic, edge cases
- **Exclude generated code** - Use build tags or separate packages

```go
// Exclude from coverage
//go:build !test

// Or use separate test builds
//go:build integration
```

---

## 10. Best Practices

### General Testing Principles

1. **Tests should be independent** - No shared state
2. **Tests should be repeatable** - Same result every time
3. **Tests should be fast** - Quick feedback loop
4. **Tests should be clear** - Easy to understand failures
5. **One assertion per test** - Or closely related assertions

### Naming Conventions

```go
// Test function names
func TestUserService_CreateUser_WithValidEmail_ReturnsUser(t *testing.T)
func TestUserService_CreateUser_WithInvalidEmail_ReturnsError(t *testing.T)

// Table test names
{"valid email", "alice@example.com", nil},
{"invalid email", "invalid", ErrInvalidEmail},
{"empty email", "", ErrEmptyEmail},
```

### Assertion Messages

```go
// Good - descriptive error messages
if result != expected {
    t.Errorf("Add(%d, %d) = %d; want %d", a, b, result, expected)
}

// Bad - unclear error
if result != expected {
    t.Error("wrong result")
}
```

### Test Organization

```
mypackage/
├── user.go
├── user_test.go          # Unit tests
├── user_integration_test.go  # Integration tests
├── testdata/             # Test fixtures
│   ├── users.json
│   └── config.yaml
└── mocks/                # Mock implementations
    └── user_repo_mock.go
```

### When to Write Tests

**Always test:**
- Business logic
- Data transformations
- Edge cases and error handling
- Public APIs

**Consider not testing:**
- Trivial getters/setters
- Third-party library wrappers (test integration instead)
- Generated code
- Main functions (test components instead)

### TDD vs BDD Decision Matrix

| Scenario | Approach | Reason |
|----------|----------|--------|
| Internal library | TDD | Focus on technical correctness |
| User-facing feature | BDD | Align with business requirements |
| Algorithm implementation | TDD | Precise technical specs |
| API endpoint | BDD | Describe user behavior |
| Data transformation | TDD | Clear inputs/outputs |
| Workflow/process | BDD | Multi-step scenarios |

---

## Practical Examples

### Example 1: TDD - String Calculator Kata

```go
// Step 1: Empty string returns 0
func TestAdd_EmptyString_ReturnsZero(t *testing.T) {
    assert.Equal(t, 0, Add(""))
}

func Add(numbers string) int {
    return 0
}

// Step 2: Single number returns that number
func TestAdd_SingleNumber_ReturnsThatNumber(t *testing.T) {
    assert.Equal(t, 1, Add("1"))
    assert.Equal(t, 5, Add("5"))
}

func Add(numbers string) int {
    if numbers == "" {
        return 0
    }
    num, _ := strconv.Atoi(numbers)
    return num
}

// Step 3: Two numbers return sum
func TestAdd_TwoNumbers_ReturnsSum(t *testing.T) {
    assert.Equal(t, 3, Add("1,2"))
    assert.Equal(t, 10, Add("5,5"))
}

func Add(numbers string) int {
    if numbers == "" {
        return 0
    }
    
    parts := strings.Split(numbers, ",")
    if len(parts) == 1 {
        num, _ := strconv.Atoi(numbers)
        return num
    }
    
    sum := 0
    for _, part := range parts {
        num, _ := strconv.Atoi(part)
        sum += num
    }
    return sum
}
```

### Example 2: BDD - Shopping Cart

```go
var _ = Describe("Shopping Cart", func() {
    var cart *ShoppingCart
    
    BeforeEach(func() {
        cart = NewShoppingCart()
    })
    
    Describe("Adding items", func() {
        Context("when cart is empty", func() {
            It("adds the item successfully", func() {
                err := cart.AddItem("apple", 2)
                
                Expect(err).NotTo(HaveOccurred())
                Expect(cart.ItemCount()).To(Equal(1))
            })
        })
        
        Context("when item already exists", func() {
            BeforeEach(func() {
                cart.AddItem("apple", 2)
            })
            
            It("increases the quantity", func() {
                cart.AddItem("apple", 3)
                
                Expect(cart.Quantity("apple")).To(Equal(5))
            })
        })
    })
    
    Describe("Calculating total", func() {
        Context("with multiple items", func() {
            BeforeEach(func() {
                cart.AddItem("apple", 2)   // $1 each
                cart.AddItem("banana", 3)  // $0.50 each
            })
            
            It("returns correct total", func() {
                total := cart.Total()
                Expect(total).To(Equal(3.50))
            })
        })
        
        Context("with discount code", func() {
            It("applies 10% discount", func() {
                cart.AddItem("apple", 10)
                cart.ApplyDiscount("SAVE10")
                
                Expect(cart.Total()).To(Equal(9.00))
            })
        })
    })
})
```

---

## Tools Summary

| Tool | Purpose | Best For |
|------|---------|----------|
| `testing` | Standard library | All projects |
| `testify` | Assertions | Cleaner test code |
| `ginkgo/gomega` | BDD framework | Complex test suites |
| `godog` | Cucumber-style BDD | Stakeholder collaboration |
| `gomock` | Mocking | Interface-heavy code |
| `httptest` | HTTP testing | API testing |
| `testcontainers` | Integration testing | Database/service tests |
| `go-sqlmock` | Database mocking | Unit testing with DB |

---

## References

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify](https://github.com/stretchr/testify)
- [Ginkgo](https://onsi.github.io/ginkgo/)
- [Godog](https://github.com/cucumber/godog)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Test Fixtures](https://dave.cheney.net/2016/05/10/test-fixtures-in-go)

---

**Remember**: Tests are code too. Keep them clean, maintainable, and valuable. Focus on testing behavior, not implementation details.
