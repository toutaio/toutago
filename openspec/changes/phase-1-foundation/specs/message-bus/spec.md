## ADDED Requirements

### Requirement: Message Interface

The framework SHALL provide a Message interface for typed communication.

#### Scenario: Message identification
- **WHEN** a message is created
- **THEN** it SHALL have a unique slug identifier
- **AND** it SHALL have a type categorization
- **AND** it SHALL support optional metadata

#### Scenario: Message validation
- **WHEN** a message is created with struct validation tags
- **THEN** the framework SHALL validate the message fields
- **AND** return validation errors if constraints are violated

### Requirement: Message Bus Publish-Subscribe

The framework SHALL provide a MessageBus for pub/sub communication.

#### Scenario: Subscribe to message type
- **WHEN** a handler subscribes to a message type
- **THEN** the bus SHALL register the handler for that type
- **AND** the handler SHALL receive all published messages of that type

#### Scenario: Publish message
- **WHEN** a message is published to the bus
- **THEN** all subscribed handlers for that message type SHALL receive it
- **AND** handlers SHALL be invoked concurrently in separate goroutines

#### Scenario: Unsubscribe handler
- **WHEN** a handler unsubscribes from a message type
- **THEN** the handler SHALL no longer receive messages of that type

### Requirement: Config-Based Message Routing

The framework SHALL support YAML-based message routing configuration.

#### Scenario: Load routing configuration
- **WHEN** routes.yaml is parsed
- **THEN** message slugs SHALL be mapped to handler identifiers
- **AND** conditional routing rules SHALL be supported
- **AND** async flags SHALL be respected

#### Scenario: Route message by configuration
- **WHEN** a message is published
- **THEN** the bus SHALL look up configured routes for that message slug
- **AND** invoke the configured handlers in specified order
- **AND** respect async/sync dispatch flags

### Requirement: Code-Based Message Routing

The framework SHALL support programmatic message routing.

#### Scenario: Subscribe programmatically
- **WHEN** handler subscribes via Subscribe() method
- **THEN** the handler SHALL receive messages without YAML configuration
- **AND** code-based routes SHALL take precedence over config-based routes

### Requirement: Handler Error Handling

The framework SHALL handle errors from message handlers gracefully.

#### Scenario: Handler returns error
- **WHEN** a message handler returns an error
- **THEN** the error SHALL be logged
- **AND** other handlers SHALL continue processing
- **AND** optionally an error message SHALL be published
