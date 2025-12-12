package findings

import "time"

// PagedResourceOfFinding represents a paged response of findings
type PagedResourceOfFinding struct {
	Embedded *EmbeddedFinding `json:"_embedded,omitempty"`
	Page     *PageMetadata    `json:"page,omitempty"`
}

// EmbeddedFinding contains the list of findings
type EmbeddedFinding struct {
	Findings []Finding `json:"findings,omitempty"`
}

// PageMetadata contains pagination information
type PageMetadata struct {
	Number        int64 `json:"number,omitempty"`
	Size          int64 `json:"size,omitempty"`
	TotalElements int64 `json:"total_elements,omitempty"`
	TotalPages    int64 `json:"total_pages,omitempty"`
}

// Finding represents a security finding
type Finding struct {
	IssueID                int64          `json:"issue_id,omitempty"`
	ScanType               ScanType       `json:"scan_type,omitempty"`
	Description            string         `json:"description,omitempty"`
	Count                  int            `json:"count,omitempty"`
	ContextType            ContextType    `json:"context_type,omitempty"`
	ContextGUID            string         `json:"context_guid,omitempty"`
	ViolatesPolicy         bool           `json:"violates_policy,omitempty"`
	FindingStatus          *FindingStatus `json:"finding_status,omitempty"`
	FindingDetails         interface{}    `json:"finding_details,omitempty"`
	Annotations            []Annotation   `json:"annotations,omitempty"`
	GracePeriodExpiresDate *time.Time     `json:"grace_period_expires_date,omitempty"`
}

// FindingStatus represents the status of a finding
type FindingStatus struct {
	FirstFoundDate         *time.Time       `json:"first_found_date,omitempty"`
	LastSeenDate           *time.Time       `json:"last_seen_date,omitempty"`
	Status                 Status           `json:"status,omitempty"`
	Resolution             string           `json:"resolution,omitempty"`
	ResolutionStatus       ResolutionStatus `json:"resolution_status,omitempty"`
	New                    bool             `json:"new,omitempty"`
	MitigationReviewStatus ResolutionStatus `json:"mitigation_review_status,omitempty"`
}

// Annotation represents a mitigation annotation on a finding
type Annotation struct {
	Action      string     `json:"action,omitempty"`
	Comment     string     `json:"comment,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	UserName    string     `json:"user_name,omitempty"`
	Description string     `json:"description,omitempty"`
	// Legacy fields for backward compatibility
	User string     `json:"user,omitempty"`
	Date *time.Time `json:"date,omitempty"`
}
