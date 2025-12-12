package findings

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	findingsBasePath = "/appsec/v2/applications"
)

// Service provides methods to interact with the Veracode Findings API
type Service struct {
	client HTTPClient
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	DoRequestWithQueryParams(method, urlPath string, params url.Values) ([]byte, error)
}

func NewService(client HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetFindingsOptions contains optional parameters for GetFindings
type GetFindingsOptions struct {
	Context            string   // Context: empty for APPLICATION, sandbox GUID for SANDBOX
	ScanType           []string // Type of scan: STATIC, DYNAMIC, MANUAL, SCA
	Severity           int      // Severity value (0-5)
	SeverityGTE        int      // Severity greater than or equal to (0-5)
	ViolatesPolicy     *bool    // Filter by policy violation
	IncludeAnnotations bool     // Include annotations in the response (not valid for SCA)
	Size               int      // Page size
	Page               int      // Page number
}

// GetFindings retrieves findings for an application
func (s *Service) GetFindings(applicationGUID string, opts *GetFindingsOptions) (*PagedResourceOfFinding, error) {
	if applicationGUID == "" {
		return nil, fmt.Errorf("applicationGUID is required")
	}

	params := url.Values{}

	if opts != nil {
		if opts.Context != "" {
			params.Add("context", opts.Context)
		}
		for _, scanType := range opts.ScanType {
			params.Add("scan_type", scanType)
		}
		if opts.Severity > 0 {
			params.Add("severity", strconv.Itoa(opts.Severity))
		}
		if opts.SeverityGTE > 0 {
			params.Add("severity_gte", strconv.Itoa(opts.SeverityGTE))
		}
		if opts.ViolatesPolicy != nil {
			params.Add("violates_policy", strconv.FormatBool(*opts.ViolatesPolicy))
		}
		if opts.IncludeAnnotations {
			params.Add("include_annot", "true")
		}
		if opts.Size > 0 {
			params.Add("size", strconv.Itoa(opts.Size))
		}
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
	}

	urlPath := fmt.Sprintf("%s/%s/findings", findingsBasePath, applicationGUID)
	body, err := s.client.DoRequestWithQueryParams("GET", urlPath, params)
	if err != nil {
		return nil, err
	}

	var result PagedResourceOfFinding
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse findings response: %w", err)
	}

	return &result, nil
}

// GetStaticFlawInfo retrieves detailed data path information for a static flaw
func (s *Service) GetStaticFlawInfo(applicationGUID string, issueID int64, context string) (*StaticFlawInfo, error) {
	if applicationGUID == "" {
		return nil, fmt.Errorf("applicationGUID is required")
	}
	if issueID == 0 {
		return nil, fmt.Errorf("issueID is required")
	}

	params := url.Values{}
	if context != "" {
		params.Add("context", context)
	}

	urlPath := fmt.Sprintf("%s/%s/findings/%d/static_flaw_info", findingsBasePath, applicationGUID, issueID)
	body, err := s.client.DoRequestWithQueryParams("GET", urlPath, params)
	if err != nil {
		return nil, err
	}

	var result StaticFlawInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse static flaw info response: %w", err)
	}

	return &result, nil
}
