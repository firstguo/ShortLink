## ADDED Requirements

### Requirement: Create short link from URL
The system SHALL accept a long URL and return a shortened URL with a unique 6-character code.

#### Scenario: Successful short link creation
- **WHEN** client sends POST /api/v1/links with valid URL
- **THEN** system returns 201 status with short code, original URL, short URL, and creation timestamp

#### Scenario: Invalid URL format
- **WHEN** client sends POST /api/v1/links with invalid URL format
- **THEN** system returns 400 status with error message "Invalid URL format"

#### Scenario: URL exceeds maximum length
- **WHEN** client sends POST /api/v1/links with URL longer than 2048 characters
- **THEN** system returns 400 status with error message

### Requirement: Validate URL format
The system SHALL validate that the provided URL conforms to RFC 3986 standard and uses HTTP or HTTPS scheme.

#### Scenario: Valid HTTP URL
- **WHEN** client provides "http://example.com/path"
- **THEN** validation passes

#### Scenario: Valid HTTPS URL
- **WHEN** client provides "https://example.com/path?query=value"
- **THEN** validation passes

#### Scenario: Invalid scheme
- **WHEN** client provides "ftp://example.com/file"
- **THEN** validation fails

#### Scenario: Empty URL
- **WHEN** client provides empty string
- **THEN** validation fails

### Requirement: Generate unique short code
The system SHALL generate a unique 6-character code using Snowflake algorithm + Base62 encoding for each short link.

#### Scenario: Generate unique codes
- **WHEN** system generates 10,000 short codes
- **THEN** all codes are unique and exactly 6 characters long

#### Scenario: Code character set
- **WHEN** system generates a short code
- **THEN** code contains only characters from [0-9a-zA-Z]

### Requirement: Store short link in database
The system SHALL persist the short link mapping to MySQL database with code, original URL, and status.

#### Scenario: Store new short link
- **WHEN** short code is generated and URL is validated
- **THEN** system inserts record into short_links table with status=1

#### Scenario: Database write failure
- **WHEN** database insertion fails
- **THEN** system returns 500 status with error message

### Requirement: Write-through cache update
The system SHALL write the short link to Redis cache simultaneously with database write to ensure consistency.

#### Scenario: Cache write on creation
- **WHEN** short link is created successfully
- **THEN** system sets Redis key "shortlink:{code}" with original URL and status, TTL 24 hours
