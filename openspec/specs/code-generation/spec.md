## ADDED Requirements

### Requirement: Generate unique ID using Snowflake algorithm
The system SHALL generate unique 64-bit IDs using Snowflake algorithm with configurable worker ID.

#### Scenario: Generate unique ID
- **WHEN** system calls Generate() method
- **THEN** system returns unique 64-bit integer

#### Scenario: IDs are time-ordered
- **WHEN** system generates two IDs at different times
- **THEN** later ID is greater than earlier ID

#### Scenario: Handle same-millisecond requests
- **WHEN** system generates multiple IDs within same millisecond
- **THEN** system increments sequence number (up to 4095)

#### Scenario: Sequence number overflow
- **WHEN** sequence number reaches 4095 in same millisecond
- **THEN** system waits for next millisecond before generating

### Requirement: Convert ID to Base62 short code
The system SHALL convert 64-bit integer to 6-character Base62 string using character set [0-9a-zA-Z].

#### Scenario: Convert ID to 6-character code
- **WHEN** system converts any valid Snowflake ID
- **THEN** result is exactly 6 characters long

#### Scenario: Character set validation
- **WHEN** system generates short code
- **THEN** code contains only characters from "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

#### Scenario: Pad short codes
- **WHEN** converted code is less than 6 characters
- **THEN** system pads with leading '0' character to reach 6 characters

### Requirement: Support distributed deployment
The system SHALL support multiple instances with unique worker IDs to avoid collisions.

#### Scenario: Configure worker ID
- **WHEN** system initializes CodeGenerator
- **THEN** system accepts worker ID parameter (0-1023)

#### Scenario: Different workers generate different codes
- **WHEN** two instances with different worker IDs generate codes at same time
- **THEN** generated codes are different

### Requirement: Achieve high throughput
The system SHALL support generating at least 4096 unique codes per millisecond per worker.

#### Scenario: Generate 10,000 codes
- **WHEN** system generates 10,000 codes sequentially
- **THEN** all codes are unique and generation completes in < 100ms
