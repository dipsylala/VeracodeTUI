package annotations

import (
	"net/url"
	"testing"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	DoRequestWithBodyFunc        func(method, urlPath string, body []byte, params url.Values) ([]byte, error)
	DoRequestWithQueryParamsFunc func(method, urlPath string, params url.Values) ([]byte, error)
}

func (m *MockHTTPClient) DoRequestWithQueryParams(method, urlPath string, params url.Values) ([]byte, error) {
	if m.DoRequestWithQueryParamsFunc != nil {
		return m.DoRequestWithQueryParamsFunc(method, urlPath, params)
	}
	return []byte("{}"), nil
}

func (m *MockHTTPClient) DoRequestWithBody(method, urlPath string, body []byte, params url.Values) ([]byte, error) {
	if m.DoRequestWithBodyFunc != nil {
		return m.DoRequestWithBodyFunc(method, urlPath, body, params)
	}
	return []byte(`{"findings":"https://api.veracode.com/application/app-guid/findings"}`), nil
}

func TestNewService(t *testing.T) {
	client := &MockHTTPClient{}
	service := NewService(client)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.client == nil {
		t.Fatal("Expected client to be set, got nil")
	}
}

func TestCreateAnnotation_Success(t *testing.T) {
	client := &MockHTTPClient{
		DoRequestWithBodyFunc: func(method, urlPath string, body []byte, params url.Values) ([]byte, error) {
			// Verify request parameters
			if method != "POST" {
				t.Errorf("Expected POST method, got %s", method)
			}
			if urlPath != "/appsec/v2/applications/test-app-guid/annotations" {
				t.Errorf("Expected correct URL path, got %s", urlPath)
			}
			return []byte(`{"findings":"https://api.veracode.com/application/test-app-guid/findings"}`), nil
		},
	}

	service := NewService(client)
	annotation := &AnnotationData{
		IssueList: "123,456",
		Comment:   "Test annotation",
		Action:    string(ActionFalsePositive),
	}

	response, err := service.CreateAnnotation("test-app-guid", annotation, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response to be returned, got nil")
	}

	if response.Findings != "https://api.veracode.com/application/test-app-guid/findings" {
		t.Errorf("Expected findings URL, got %s", response.Findings)
	}
}

func TestCreateAnnotation_WithContext(t *testing.T) {
	client := &MockHTTPClient{
		DoRequestWithBodyFunc: func(method, urlPath string, body []byte, params url.Values) ([]byte, error) {
			// Verify context parameter is passed
			if params.Get("context") != "sandbox-guid" {
				t.Errorf("Expected context parameter to be sandbox-guid, got %s", params.Get("context"))
			}
			return []byte(`{"findings":"https://api.veracode.com/application/app-guid/findings"}`), nil
		},
	}

	service := NewService(client)
	annotation := &AnnotationData{
		IssueList: "123",
		Comment:   "Test",
		Action:    string(ActionComment),
	}
	opts := &CreateAnnotationOptions{
		Context: "sandbox-guid",
	}

	_, err := service.CreateAnnotation("app-guid", annotation, opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateAnnotation_MissingApplicationGUID(t *testing.T) {
	client := &MockHTTPClient{}
	service := NewService(client)

	annotation := &AnnotationData{
		IssueList: "123",
	}

	_, err := service.CreateAnnotation("", annotation, nil)
	if err == nil {
		t.Fatal("Expected error for missing applicationGUID, got nil")
	}
}

func TestCreateAnnotation_MissingAnnotation(t *testing.T) {
	client := &MockHTTPClient{}
	service := NewService(client)

	_, err := service.CreateAnnotation("app-guid", nil, nil)
	if err == nil {
		t.Fatal("Expected error for missing annotation, got nil")
	}
}

func TestCreateAnnotation_MissingIssueList(t *testing.T) {
	client := &MockHTTPClient{}
	service := NewService(client)

	annotation := &AnnotationData{
		Comment: "Test",
		Action:  string(ActionComment),
	}

	_, err := service.CreateAnnotation("app-guid", annotation, nil)
	if err == nil {
		t.Fatal("Expected error for missing issue_list, got nil")
	}
}

func TestAnnotationDataConstruction(t *testing.T) {
	annotation := &AnnotationData{
		IssueList: "123,456,789",
		Comment:   "Test comment",
		Action:    string(ActionFalsePositive),
	}

	opts := &CreateAnnotationOptions{
		Context: "sandbox-guid",
	}

	if annotation.IssueList != "123,456,789" {
		t.Errorf("Expected IssueList to be '123,456,789', got %s", annotation.IssueList)
	}

	if annotation.Comment != "Test comment" {
		t.Errorf("Expected Comment to be 'Test comment', got %s", annotation.Comment)
	}

	if annotation.Action != string(ActionFalsePositive) {
		t.Errorf("Expected Action to be FP, got %s", annotation.Action)
	}

	if opts.Context != "sandbox-guid" {
		t.Errorf("Expected Context to be 'sandbox-guid', got %s", opts.Context)
	}
}

func TestAnnotationDirectCreation(t *testing.T) {
	client := &MockHTTPClient{
		DoRequestWithBodyFunc: func(method, urlPath string, body []byte, params url.Values) ([]byte, error) {
			return []byte(`{"findings":"https://api.veracode.com/application/app-guid/findings"}`), nil
		},
	}
	service := NewService(client)

	annotation := &AnnotationData{
		IssueList: "123",
		Comment:   "Direct creation test",
		Action:    string(ActionAcceptRisk),
	}

	response, err := service.CreateAnnotation("app-guid", annotation, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestAnnotationActions(t *testing.T) {
	actions := []AnnotationAction{
		ActionComment,
		ActionFalsePositive,
		ActionAppDesign,
		ActionOSEnv,
		ActionNetEnv,
		ActionRejected,
		ActionAccepted,
		ActionLibrary,
		ActionAcceptRisk,
	}

	expectedValues := []string{
		"COMMENT",
		"FP",
		"APPDESIGN",
		"OSENV",
		"NETENV",
		"REJECTED",
		"ACCEPTED",
		"LIBRARY",
		"ACCEPTRISK",
	}

	for i, action := range actions {
		if string(action) != expectedValues[i] {
			t.Errorf("Expected action %s to have value %s, got %s", action, expectedValues[i], string(action))
		}
	}
}
