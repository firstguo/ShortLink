## ADDED Requirements

### Requirement: Redirect short code to original URL
The system SHALL redirect users to the original URL when they access a valid short code via HTTP 302 status.

#### Scenario: Successful redirect
- **WHEN** user sends GET /{code} with valid short code
- **THEN** system returns 302 Found with Location header set to original URL

#### Scenario: Short code not found
- **WHEN** user sends GET /{code} with non-existent short code
- **THEN** system returns 404 status with error message "Short link not found"

#### Scenario: Short link disabled
- **WHEN** user sends GET /{code} with disabled short link (status=0)
- **THEN** system returns 404 status with error message "Short link not found"

### Requirement: Query cache first for redirect
The system SHALL check Redis cache before querying the database to achieve P99 response time < 50ms.

#### Scenario: Cache hit
- **WHEN** short code exists in Redis cache
- **THEN** system returns original URL from cache without database query

#### Scenario: Cache miss
- **WHEN** short code does not exist in Redis cache
- **THEN** system queries database and populates cache

### Requirement: Populate cache on database hit
The system SHALL write to Redis cache when a short code is found in database but not in cache.

#### Scenario: Cache population on miss
- **WHEN** database query returns valid short link
- **THEN** system sets Redis key with TTL 24 hours ± 10% random offset

#### Scenario: Random TTL prevents cache stampede
- **WHEN** multiple requests cause cache population
- **THEN** each cache entry has slightly different TTL (±2.4 hours)

### Requirement: Cache non-existent codes to prevent穿透
The system SHALL cache a "NULL" value for non-existent short codes to prevent cache penetration attacks.

#### Scenario: Cache NULL for missing code
- **WHEN** database query returns no record for short code
- **THEN** system sets Redis key "shortlink:{code}" to "NULL" with TTL 5 minutes

#### Scenario: Return 404 for cached NULL
- **WHEN** Redis cache contains "NULL" for short code
- **THEN** system returns 404 without database query

### Requirement: Handle Redis failure gracefully
The system SHALL degrade to database queries when Redis is unavailable.

#### Scenario: Redis connection failure
- **WHEN** Redis is down or unreachable
- **THEN** system queries database directly and returns result

#### Scenario: Redis timeout
- **WHEN** Redis query times out
- **THEN** system logs error and falls back to database query
