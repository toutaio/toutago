## ADDED Requirements

### Requirement: Router Interface Abstraction

The framework SHALL provide a Router interface that abstracts HTTP routing operations.

#### Scenario: Register HTTP routes
- **WHEN** developer registers routes with GET, POST, PUT, DELETE, PATCH methods
- **THEN** the router SHALL accept a path and handler function
- **AND** route incoming requests to the registered handler

#### Scenario: Route grouping
- **WHEN** developer creates a route group with a prefix
- **THEN** all routes in the group SHALL inherit the prefix
- **AND** group-specific middleware SHALL apply to all routes in the group

#### Scenario: Middleware registration
- **WHEN** developer registers middleware via Use()
- **THEN** middleware SHALL be applied to all subsequent routes
- **AND** middleware SHALL execute in registration order

### Requirement: Chi Router Default Implementation

The framework SHALL provide a Chi-based implementation of the Router interface.

#### Scenario: Adapt Chi router
- **WHEN** ChiRouter is instantiated
- **THEN** it SHALL wrap the Chi mux
- **AND** convert BoxIt handlers to Chi-compatible handlers
- **AND** preserve all Chi functionality

#### Scenario: Access native router
- **WHEN** developer calls Native() method
- **THEN** the underlying Chi mux SHALL be returned
- **AND** advanced Chi features SHALL be accessible

### Requirement: Request Context

The framework SHALL provide a Context interface for request-scoped data.

#### Scenario: Access request and response
- **WHEN** handler receives a Context
- **THEN** it SHALL provide access to http.Request
- **AND** it SHALL provide access to http.ResponseWriter
- **AND** it SHALL extract route parameters

#### Scenario: Request-scoped data storage
- **WHEN** handler stores data in Context via Set()
- **THEN** subsequent handlers SHALL retrieve it via Get()
- **AND** data SHALL be isolated per request

#### Scenario: Container access from Context
- **WHEN** handler accesses Context.Container()
- **THEN** the DI container SHALL be available
- **AND** handlers SHALL resolve dependencies at runtime

### Requirement: Graceful Shutdown

The framework SHALL support graceful HTTP server shutdown.

#### Scenario: Shutdown signal
- **WHEN** application receives shutdown signal (SIGINT, SIGTERM)
- **THEN** the server SHALL stop accepting new connections
- **AND** wait for active requests to complete
- **AND** exit cleanly within configured timeout
