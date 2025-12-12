package annotations

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	annotationsBasePath = "/appsec/v2/applications"
)

// Service provides methods to interact with the Veracode Annotations API
type Service struct {
	client HTTPClient
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	DoRequestWithQueryParams(method, urlPath string, params url.Values) ([]byte, error)
	DoRequestWithBody(method, urlPath string, body []byte, params url.Values) ([]byte, error)
}

func NewService(client HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// CreateAnnotationOptions contains optional parameters for CreateAnnotation
type CreateAnnotationOptions struct {
	Context string // GUID of the specified development sandbox
}

// CreateAnnotation creates an annotation for findings in an application
func (s *Service) CreateAnnotation(applicationGUID string, annotation *AnnotationData, opts *CreateAnnotationOptions) (*AnnotationResponse, error) {
	if applicationGUID == "" {
		return nil, fmt.Errorf("applicationGUID is required")
	}
	if annotation == nil {
		return nil, fmt.Errorf("annotation is required")
	}
	if annotation.IssueList == "" {
		return nil, fmt.Errorf("issue_list is required")
	}

	params := url.Values{}
	if opts != nil && opts.Context != "" {
		params.Add("context", opts.Context)
	}

	jsonBody, err := json.Marshal(annotation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal annotation: %w", err)
	}

	urlPath := fmt.Sprintf("%s/%s/annotations", annotationsBasePath, applicationGUID)
	body, err := s.client.DoRequestWithBody("POST", urlPath, jsonBody, params)
	if err != nil {
		return nil, err
	}

	var result AnnotationResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse annotation response: %w", err)
		}
	}

	return &result, nil
}
