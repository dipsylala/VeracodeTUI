# Annotations Service

Service layer for interacting with the Veracode Annotations API (`/appsec/v2/applications/{application_guid}/annotations`).

## Overview

This service provides a clean Go interface to the Veracode Annotations API, enabling you to create annotations (mitigations) for findings in your applications.

## Features

- ✅ Create annotations for one or more findings
- ✅ Support for multiple mitigation actions
- ✅ Sandbox context support
- ✅ Fluent builder pattern for easy annotation creation
- ✅ Full type safety with Go structs
- ✅ Comprehensive action type constants

## Usage

### Initialize the Service

```go
import (
    "github.com/dipsylala/veracode-tui/config"
    "github.com/dipsylala/veracode-tui/services/annotations"
    "github.com/dipsylala/veracode-tui/veracode"
)

// Load configuration
cfg, _ := config.LoadConfig()
keyID, keySecret := cfg.GetAPICredentials()

// Create client and service
client := veracode.NewClient(keyID, keySecret)
service := annotations.NewService(client)
```

### Create Annotation (Direct Method)

```go
annotation := &annotations.AnnotationData{
    IssueList: "123,456,789",
    Comment:   "This is a false positive - input is validated upstream",
    Action:    string(annotations.ActionFalsePositive),
}

opts := &annotations.CreateAnnotationOptions{
    Context: "sandbox-guid-here", // Optional: for sandbox findings
}

response, err := service.CreateAnnotation("app-guid-here", annotation, opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Annotation created. Findings URL: %s\n", response.Findings)
```

### Create Annotation

```go
annotation := &annotations.AnnotationData{
    IssueList: "123,456,789",
    Comment:   "This is a false positive",
    Action:    string(annotations.ActionFalsePositive),
}

opts := &annotations.CreateAnnotationOptions{
    Context: "sandbox-guid", // Optional: for sandbox context
}

response, err := service.CreateAnnotation("app-guid-here", annotation, opts)
if err != nil {
    log.Fatal(err)
}
```

### Mitigation Actions

The service provides constants for all supported mitigation actions:

```go
// Mark as false positive
annotation := &annotations.AnnotationData{
    IssueList: "123",
    Comment:   "Not a real vulnerability",
    Action:    string(annotations.ActionFalsePositive),
}
service.CreateAnnotation(appGUID, annotation, nil)

// Accept risk
annotation := &annotations.AnnotationData{
    IssueList: "456",
    Comment:   "Risk accepted by security team",
    Action:    string(annotations.ActionAcceptRisk),
}
service.CreateAnnotation(appGUID, annotation, nil)

// Mitigated by app design
annotation := &annotations.AnnotationData{
    IssueList: "789",
    Comment:   "Input validation prevents exploitation",
    Action:    string(annotations.ActionAppDesign),
}
service.CreateAnnotation(appGUID, annotation, nil)
```

### Multiple Findings

Create an annotation for multiple findings at once:

```go
annotation := &annotations.AnnotationData{
    IssueList: "123,456,789,101,112",  // Comma-separated list
    Comment:   "All of these are false positives",
    Action:    string(annotations.ActionFalsePositive),
}

response, err := service.CreateAnnotation(appGUID, annotation, nil)
```

### Sandbox Annotations

Annotate findings in a sandbox:

```go
opts := &annotations.CreateAnnotationOptions{
    Context: "sandbox-guid-here",
}

annotation := &annotations.AnnotationData{
    IssueList: "123",
    Comment:   "False positive in sandbox",
    Action:    string(annotations.ActionFalsePositive),
}

response, err := service.CreateAnnotation(appGUID, annotation, opts)
```

## API Endpoint

| Method | Endpoint | Description |
|--------|----------|-------------|
| `CreateAnnotation` | `POST /appsec/v2/applications/{guid}/annotations` | Create annotation for findings |

## Annotation Actions

| Constant | Value | Description |
|----------|-------|-------------|
| `ActionComment` | `COMMENT` | Add a comment to the finding |
| `ActionFalsePositive` | `FP` | Mark as false positive |
| `ActionAppDesign` | `APPDESIGN` | Mitigated by application design |
| `ActionOSEnv` | `OSENV` | Mitigated by OS environment |
| `ActionNetEnv` | `NETENV` | Mitigated by network environment |
| `ActionRejected` | `REJECTED` | Mitigation proposal rejected |
| `ActionAccepted` | `ACCEPTED` | Mitigation proposal accepted |
| `ActionLibrary` | `LIBRARY` | Finding exists in a library |
| `ActionAcceptRisk` | `ACCEPTRISK` | Risk accepted |

## Data Models

### AnnotationData

Request payload for creating an annotation:

```go
type AnnotationData struct {
    IssueList string  // Comma-separated list of flaw IDs (required)
    Comment   string  // Annotation comment
    Action    string  // Mitigation action (use AnnotationAction constants)
}
```

### AnnotationResponse

Response from creating an annotation:

```go
type AnnotationResponse struct {
    Findings string  // Link to the findings URL
}
```

### CreateAnnotationOptions

Optional parameters:

```go
type CreateAnnotationOptions struct {
    Context string  // GUID of sandbox (empty for application context)
}
```

## Builder Pattern

The `CreateAnnotationBuilder` provides a fluent interface:

```go
builder := annotations.NewAnnotationBuilder()

// Chain methods
builder.
    WithIssueList("123,456").
    WithComment("Your comment").
    WithAction(annotations.ActionFalsePositive).
    WithContext("sandbox-guid")

// Get the annotation data
annotation, opts := builder.Build()

// Or create directly
response, err := builder.CreateWithService(service, appGUID)
```

## Error Handling

All methods return errors for:
- Authentication failures (401)
- Authorization failures (403)
- Application not found (404)
- Rate limiting (429)
- Invalid requests (400)
- Server errors (500)

```go
annotation := &annotations.AnnotationData{
    IssueList: "123",
    // Missing required fields or invalid action
}

_, err := service.CreateAnnotation(appGUID, annotation, nil)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    // Error: issue_list is required
}
```

## Response Codes

| Code | Description |
|------|-------------|
| 200 | Request successfully submitted |
| 202 | Annotation applied |
| 401 | Not authorized |
| 403 | Access denied |
| 404 | Application not found |
| 429 | Request limit exceeded |
| 500 | Server error |

## Example: Bulk Annotation Workflow

```go
// Get findings that need annotation
findings, _ := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
    ViolatesPolicy: boolPtr(true),
    Severity:       5,
})

// Collect issue IDs
var issueIDs []string
for _, finding := range findings.Embedded.Findings {
    if shouldAnnotate(finding) {
        issueIDs = append(issueIDs, fmt.Sprintf("%d", finding.IssueID))
    }
}

// Create bulk annotation
if len(issueIDs) > 0 {
    response, err := annotations.NewAnnotationBuilder().
        WithIssueList(strings.Join(issueIDs, ",")).
        WithComment("Reviewed and marked as false positives").
        WithAction(annotations.ActionFalsePositive).
        CreateWithService(annotationService, appGUID)
}
```

## Architecture

```
services/
└── annotations/
    ├── doc.go        # Package documentation
    ├── models.go     # Data structures (AnnotationData, AnnotationResponse, etc.)
    ├── service.go    # Service implementation with builder pattern
    └── README.md     # This file
```

The service uses the `veracode.Client` HTTP client for authenticated requests and handles:
- URL construction
- Query parameter encoding
- JSON marshaling/unmarshaling
- Error handling

## See Also

- [Veracode Annotations API Documentation](https://docs.veracode.com/r/c_annotations_intro)
- [Applications Service](../applications/README.md)
- [Findings Service](../findings/README.md)
- [Main README](../../README.md)
