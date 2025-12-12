# Applications Service

Service layer for interacting with the Veracode Applications API (`/appsec/v1/applications`).

## Overview

This service provides a clean Go interface to the Veracode Applications API, implementing all GET endpoints as defined in the Swagger specification.

## Features

- ✅ Get all applications with filtering and pagination
- ✅ Get single application by GUID
- ✅ Get sandboxes for an application
- ✅ Get single sandbox by GUID
- ✅ Full type safety with Go structs
- ✅ Integration tests using `~/.veracode/veracode.yml` credentials

## Usage

### Initialize the Service

```go
import (
    "github.com/dipsylala/veracode-tui/config"
    "github.com/dipsylala/veracode-tui/services/applications"
    "github.com/dipsylala/veracode-tui/veracode"
)

// Load configuration
cfg, _ := config.LoadConfig()
keyID, keySecret := cfg.GetAPICredentials()

// Create client and service
client := veracode.NewClient(keyID, keySecret)
service := applications.NewService(client)
```

### Get All Applications

```go
// Get all applications (with defaults)
result, err := service.GetApplications(nil)
if err != nil {
    log.Fatal(err)
}

for _, app := range result.Embedded.Applications {
    fmt.Printf("App: %s (GUID: %s)\n", app.Profile.Name, app.GUID)
}
```

### Get Applications with Filtering

```go
opts := &applications.GetApplicationsOptions{
    Name: "MyApp",
    Page: 0,
    Size: 50,
    BusinessUnit: "Engineering",
    ScanType: "STATIC",
    PolicyCompliance: "PASSED",
}

result, err := service.GetApplications(opts)
```

### Get Single Application

```go
app, err := service.GetApplication("app-guid-here")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Application: %s\n", app.Profile.Name)
fmt.Printf("Business Unit: %s\n", app.Profile.BusinessUnit.Name)
fmt.Printf("Policies: %d\n", len(app.Profile.Policies))
```

### Get Sandboxes

```go
// Get all sandboxes for an application
sandboxes, err := service.GetSandboxes("app-guid", nil)
if err != nil {
    log.Fatal(err)
}

for _, sandbox := range sandboxes.Embedded.Sandboxes {
    fmt.Printf("Sandbox: %s (GUID: %s)\n", sandbox.Name, sandbox.GUID)
}
```

### Get Single Sandbox

```go
sandbox, err := service.GetSandbox("app-guid", "sandbox-guid")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Sandbox: %s\n", sandbox.Name)
fmt.Printf("Owner: %s\n", sandbox.OwnerUsername)
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GetApplications` | `GET /appsec/v1/applications` | List applications with optional filtering |
| `GetApplication` | `GET /appsec/v1/applications/{guid}` | Get single application details |
| `GetSandboxes` | `GET /appsec/v1/applications/{guid}/sandboxes` | List sandboxes for an application |
| `GetSandbox` | `GET /appsec/v1/applications/{guid}/sandboxes/{sandboxGuid}` | Get single sandbox details |

## Filtering Options

The `GetApplications` method supports extensive filtering:

- `BusinessUnit` - Filter by business unit name
- `CustomFieldNames` - Array of custom field names
- `CustomFieldValues` - Array of custom field values
- `LegacyID` - Filter by legacy application ID
- `ModifiedAfter` - Only apps modified after date (yyyy-MM-dd)
- `Name` - Filter by application name
- `Page` - Page number (defaults to 0)
- `Size` - Page size, up to 500 (default 50)
- `Policy` - Filter by policy name
- `PolicyCompliance` - Filter by compliance status (DETERMINING, NOT_ASSESSED, DID_NOT_PASS, CONDITIONAL_PASS, PASSED, VENDOR_REVIEW)
- `PolicyComplianceCheckedAfter` - Filter by policy compliance check date
- `PolicyGUID` - Filter by policy GUID
- `ScanStatus` - Array of scan statuses
- `ScanType` - Filter by scan type (STATIC, DYNAMIC, MANUAL)
- `Tag` - Filter by tag
- `Team` - Filter by team name
- `SortByCustomFieldName` - Custom field to sort by

## Integration Tests

Run integration tests (requires valid `~/.veracode/veracode.yml`):

```powershell
# Run all tests
go test -v ./services/applications/...

# Run specific test
go test -v ./services/applications/... -run TestGetApplications

# Skip integration tests (short mode)
go test -short ./services/applications/...
```

### Test Coverage

- ✅ Get all applications
- ✅ Get applications with pagination
- ✅ Get applications filtered by name
- ✅ Get single application by GUID
- ✅ Get application with invalid GUID (error handling)
- ✅ Get sandboxes for application
- ✅ Get single sandbox by GUID

## Data Models

### Application

Main application object with full profile information:

```go
type Application struct {
    GUID                   string
    ID                     int
    Profile                *ApplicationProfile
    Scans                  []ApplicationScan
    Created                *time.Time
    Modified               *time.Time
    LastCompletedScanDate  *time.Time
    // ... more fields
}
```

### ApplicationProfile

```go
type ApplicationProfile struct {
    Name               string
    Description        string
    BusinessCriticality string
    BusinessUnit       *BusinessUnit
    BusinessOwners     []BusinessOwner
    Policies           []AppPolicy
    Teams              []AppTeam
    Tags               string
    CustomFields       []CustomNameValue
    Settings           *ApplicationSettings
    // ... more fields
}
```

### Sandbox

```go
type Sandbox struct {
    GUID            string
    ID              int
    Name            string
    ApplicationGUID string
    OwnerUsername   string
    AutoRecreate    bool
    CustomFields    []CustomNameValue
    Created         *time.Time
    Modified        *time.Time
    // ... more fields
}
```

## Error Handling

All methods return errors for:
- Authentication failures (401)
- Authorization failures (403)
- Not found errors (404)
- Rate limiting (429)
- Invalid requests (400)
- Server errors (500)

```go
app, err := service.GetApplication("invalid-guid")
if err != nil {
    // Handle error - includes HTTP status and response body
    fmt.Printf("Error: %v\n", err)
}
```

## Pagination

Results are paginated with metadata:

```go
result, _ := service.GetApplications(&applications.GetApplicationsOptions{
    Size: 100,
    Page: 0,
})

fmt.Printf("Page: %d/%d\n", result.Page.Number, result.Page.TotalPages)
fmt.Printf("Total: %d applications\n", result.Page.TotalElements)
```

## Future Enhancements

The service currently implements only GET operations. Future additions will include:

- POST `/appsec/v1/applications` - Create application
- PUT `/appsec/v1/applications/{guid}` - Update application
- DELETE `/appsec/v1/applications/{guid}` - Delete application
- POST `/appsec/v1/applications/{guid}/sandboxes` - Create sandbox
- Additional filtering and sorting options

## Architecture

```
services/
└── applications/
    ├── models.go        # Data structures (Application, Sandbox, etc.)
    ├── service.go       # Service implementation
    └── service_test.go  # Integration tests
```

The service uses the `veracode.Client` HTTP client for authenticated requests and handles:
- URL construction
- Query parameter encoding
- JSON unmarshaling
- Error handling

## See Also

- [Veracode Applications API Documentation](https://docs.veracode.com/r/c_applications_intro)
- [Main README](../../README.md)
- [Veracode HMAC Authentication](../../veracode/auth.go)
