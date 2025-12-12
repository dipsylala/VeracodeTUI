package findings

// ScanType represents the type of security scan
type ScanType string

// Scan types
const (
	ScanTypeStatic  ScanType = "STATIC"
	ScanTypeDynamic ScanType = "DYNAMIC"
	ScanTypeSCA     ScanType = "SCA"
	ScanTypeManual  ScanType = "MANUAL"
)

// Status represents the status of a finding
type Status string

// Finding statuses
const (
	StatusOpen     Status = "OPEN"
	StatusClosed   Status = "CLOSED"
	StatusReopened Status = "REOPENED"
)

// ResolutionStatus represents the resolution status of a finding
type ResolutionStatus string

// Resolution statuses
const (
	ResolutionNone       ResolutionStatus = "NONE"
	ResolutionApproved   ResolutionStatus = "APPROVED"
	ResolutionProposed   ResolutionStatus = "PROPOSED"
	ResolutionRejected   ResolutionStatus = "REJECTED"
	ResolutionPending    ResolutionStatus = "PENDING"
	ResolutionCanceled   ResolutionStatus = "CANCELED"
	ResolutionAccepted   ResolutionStatus = "ACCEPTED"
	ResolutionUnresolved ResolutionStatus = "UNRESOLVED"
)

// ContextType represents the context of a finding (application or sandbox)
type ContextType string

// Context types
const (
	ContextTypeApplication ContextType = "APPLICATION"
	ContextTypeSandbox     ContextType = "SANDBOX"
)

// ScanFilterType represents the scan type filter for UI
type ScanFilterType string

// Scan filter types
const (
	ScanFilterStatic  ScanFilterType = "STATIC"
	ScanFilterDynamic ScanFilterType = "DYNAMIC"
	ScanFilterSCA     ScanFilterType = "SCA"
)

// PolicyFilterType represents the policy filter for UI
type PolicyFilterType string

// Policy filter options
const (
	PolicyFilterAll           PolicyFilterType = "All"
	PolicyFilterViolations    PolicyFilterType = "Violations"
	PolicyFilterNonViolations PolicyFilterType = "Non-Violations"
)

// Severity levels
