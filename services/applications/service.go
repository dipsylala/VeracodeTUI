package applications

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	applicationsBasePath = "/appsec/v1/applications"
)

// Service provides methods to interact with the Veracode Applications API
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

// GetApplicationsOptions contains optional parameters for GetApplications
type GetApplicationsOptions struct {
	BusinessUnit                 string
	CustomFieldNames             []string
	CustomFieldValues            []string
	LegacyID                     int
	ModifiedAfter                string // Format: yyyy-MM-dd
	Name                         string
	Page                         int
	Policy                       string
	PolicyCompliance             string
	PolicyComplianceCheckedAfter string // Format: yyyy-MM-dd
	PolicyGUID                   string
	ScanStatus                   []string
	ScanType                     string
	Size                         int
	SortByCustomFieldName        string
	Tag                          string
	Team                         string
}

// GetApplications retrieves a list of applications with optional filtering
func (s *Service) GetApplications(opts *GetApplicationsOptions) (*PagedResourceOfApplication, error) {
	params := buildApplicationQueryParams(opts)

	body, err := s.client.DoRequestWithQueryParams("GET", applicationsBasePath, params)
	if err != nil {
		return nil, err
	}

	var result PagedResourceOfApplication
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse applications response: %w", err)
	}

	return &result, nil
}

// buildApplicationQueryParams builds URL query parameters from options
//
//nolint:gocyclo // Parameter building with many optional fields
func buildApplicationQueryParams(opts *GetApplicationsOptions) url.Values {
	params := url.Values{}

	if opts == nil {
		return params
	}

	if opts.BusinessUnit != "" {
		params.Add("business_unit", opts.BusinessUnit)
	}
	for _, name := range opts.CustomFieldNames {
		params.Add("custom_field_names", name)
	}
	for _, value := range opts.CustomFieldValues {
		params.Add("custom_field_values", value)
	}
	if opts.LegacyID > 0 {
		params.Add("legacy_id", strconv.Itoa(opts.LegacyID))
	}
	if opts.ModifiedAfter != "" {
		params.Add("modified_after", opts.ModifiedAfter)
	}
	if opts.Name != "" {
		params.Add("name", opts.Name)
	}
	if opts.Page > 0 {
		params.Add("page", strconv.Itoa(opts.Page))
	}
	if opts.Policy != "" {
		params.Add("policy", opts.Policy)
	}
	if opts.PolicyCompliance != "" {
		params.Add("policy_compliance", opts.PolicyCompliance)
	}
	if opts.PolicyComplianceCheckedAfter != "" {
		params.Add("policy_compliance_checked_after", opts.PolicyComplianceCheckedAfter)
	}
	if opts.PolicyGUID != "" {
		params.Add("policy_guid", opts.PolicyGUID)
	}
	for _, status := range opts.ScanStatus {
		params.Add("scan_status", status)
	}
	if opts.ScanType != "" {
		params.Add("scan_type", opts.ScanType)
	}
	if opts.Size > 0 {
		params.Add("size", strconv.Itoa(opts.Size))
	}
	if opts.SortByCustomFieldName != "" {
		params.Add("sort_by_custom_field_name", opts.SortByCustomFieldName)
	}
	if opts.Tag != "" {
		params.Add("tag", opts.Tag)
	}
	if opts.Team != "" {
		params.Add("team", opts.Team)
	}

	return params
}

// GetApplication retrieves a single application by GUID
func (s *Service) GetApplication(applicationGUID string) (*Application, error) {
	if applicationGUID == "" {
		return nil, fmt.Errorf("applicationGUID is required")
	}

	urlPath := fmt.Sprintf("%s/%s", applicationsBasePath, applicationGUID)
	body, err := s.client.DoRequestWithQueryParams("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}

	var result Application
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse application response: %w", err)
	}

	return &result, nil
}

// GetSandboxesOptions contains optional parameters for GetSandboxes
type GetSandboxesOptions struct {
	Page int
	Size int
}

// GetSandboxes retrieves sandboxes for a specific application
func (s *Service) GetSandboxes(applicationGUID string, opts *GetSandboxesOptions) (*PagedResourceOfSandbox, error) {
	if applicationGUID == "" {
		return nil, fmt.Errorf("applicationGUID is required")
	}

	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
		if opts.Size > 0 {
			params.Add("size", strconv.Itoa(opts.Size))
		}
	}

	urlPath := fmt.Sprintf("%s/%s/sandboxes", applicationsBasePath, applicationGUID)
	body, err := s.client.DoRequestWithQueryParams("GET", urlPath, params)
	if err != nil {
		return nil, err
	}

	var result PagedResourceOfSandbox
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse sandboxes response: %w", err)
	}

	return &result, nil
}

// GetSandbox retrieves a single sandbox by application GUID and sandbox GUID
func (s *Service) GetSandbox(applicationGUID, sandboxGUID string) (*Sandbox, error) {
	if applicationGUID == "" {
		return nil, fmt.Errorf("applicationGUID is required")
	}
	if sandboxGUID == "" {
		return nil, fmt.Errorf("sandboxGUID is required")
	}

	urlPath := fmt.Sprintf("%s/%s/sandboxes/%s", applicationsBasePath, applicationGUID, sandboxGUID)
	body, err := s.client.DoRequestWithQueryParams("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}

	var result Sandbox
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse sandbox response: %w", err)
	}

	return &result, nil
}
