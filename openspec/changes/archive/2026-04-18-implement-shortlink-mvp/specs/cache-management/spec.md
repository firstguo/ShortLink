## ADDED Requirements

### Requirement: Write-through cache on create
The system SHALL write to Redis cache simultaneously with database write when creating a short link.

#### Scenario: Cache write on successful creation
- **WHEN** short link is created and saved to database
- **THEN** system sets Redis key "shortlink:{code}" with JSON containing original_url and status

#### Scenario: Set cache TTL
- **WHEN** cache entry is created
- **THEN** TTL is set to 24 hours with ±10% random offset

### Requirement: Cache-Aside pattern for reads
The system SHALL check cache first and fall back to database on cache miss when retrieving short links.

#### Scenario: Cache hit returns data
- **WHEN** Redis contains valid data for short code
- **THEN** system returns data from cache without database query

#### Scenario: Cache miss queries database
- **WHEN** Redis does not contain data for short code
- **THEN** system queries database and populates cache

### Requirement: Null value caching for non-existent codes
The system SHALL cache a special "NULL" value when a short code is not found in database to prevent cache penetration.

#### Scenario: Cache NULL for missing code
- **WHEN** database query returns no record
- **THEN** system sets Redis key to "NULL" string with TTL 5 minutes

#### Scenario: Skip database for cached NULL
- **WHEN** cache contains "NULL" for short code
- **THEN** system returns not found error immediately

### Requirement: Handle Redis errors gracefully
The system SHALL continue operating with database fallback when Redis encounters errors.

#### Scenario: Redis connection error on read
- **WHEN** Redis connection fails during cache lookup
- **THEN** system logs error and queries database directly

#### Scenario: Redis connection error on write
- **WHEN** Redis connection fails during cache write
- **THEN** system logs error but continues (database is source of truth)

#### Scenario: Redis timeout
- **WHEN** Redis operation times out
- **THEN** system logs timeout and proceeds with database fallback

### Requirement: Serialize cache data as JSON
The system SHALL serialize short link data as JSON for Redis storage to support multiple fields.

#### Scenario: Store structured data
- **WHEN** system writes to cache
- **THEN** data is JSON string with fields: original_url, status

#### Scenario: Parse cached JSON
- **WHEN** system reads from cache
- **THEN** system deserializes JSON and reconstructs ShortLink object
