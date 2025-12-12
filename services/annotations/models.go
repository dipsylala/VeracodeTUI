package annotations

// AnnotationResponse represents the response from creating an annotation
type AnnotationResponse struct {
	Findings string `json:"findings,omitempty"`
}

// AnnotationErrorResponse represents an error response from the annotations API
type AnnotationErrorResponse struct {
	Embedded struct {
		APIErrors []APIError `json:"api_errors,omitempty"`
	} `json:"_embedded,omitempty"`
}

// APIError represents a single API error
type APIError struct {
	ID     string `json:"id,omitempty"`
	Code   string `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Status string `json:"status,omitempty"`
}

// AnnotationData represents the data required to create an annotation
type AnnotationData struct {
	IssueList string `json:"issue_list,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Action    string `json:"action,omitempty"`
}

// AnnotationAction represents valid annotation actions
type AnnotationAction string

const (
	// ActionComment adds a comment to the finding
	ActionComment AnnotationAction = "COMMENT"
	// ActionFalsePositive marks the finding as a false positive
	ActionFalsePositive AnnotationAction = "FP"
	// ActionAppDesign marks the finding as mitigated by app design
	ActionAppDesign AnnotationAction = "APPDESIGN"
	// ActionOSEnv marks the finding as mitigated by OS environment
	ActionOSEnv AnnotationAction = "OSENV"
	// ActionNetEnv marks the finding as mitigated by network environment
	ActionNetEnv AnnotationAction = "NETENV"
	// ActionRejected marks the mitigation proposal as rejected
	ActionRejected AnnotationAction = "REJECTED"
	// ActionAccepted marks the mitigation proposal as accepted
	ActionAccepted AnnotationAction = "ACCEPTED"
	// ActionLibrary marks the finding as existing in a library
	ActionLibrary AnnotationAction = "LIBRARY"
	// ActionAcceptRisk marks the finding risk as accepted
	ActionAcceptRisk AnnotationAction = "ACCEPTRISK"
)
