package veracode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Veracode API endpoint URLs
const (
	BaseWebURL        = "https://analysiscenter.veracode.com/"
	BaseAPIURL        = "https://analysiscenter.veracode.com/api"
	HealthCheckAPIURL = "https://api.veracode.com"
	AppSecAPIURL      = "https://api.veracode.com"
)

// HTTPError represents an HTTP error response from the Veracode API
type HTTPError struct {
	StatusCode int    // HTTP status code (e.g., 400, 404, 500)
	Status     string // HTTP status text (e.g., "Bad Request")
	Body       []byte // Raw response body
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, string(e.Body))
}

// Client represents a Veracode API client
type Client struct {
	apiKeyID     string
	apiKeySecret string
	httpClient   *http.Client
	debugLogger  *log.Logger
	debugFile    *os.File
}

func NewClient(apiKeyID, apiKeySecret string) *Client {
	return &Client{
		apiKeyID:     apiKeyID,
		apiKeySecret: apiKeySecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Application represents a Veracode application
type Application struct {
	ID           int    `json:"id"`
	GUID         string `json:"guid"`
	Name         string `json:"profile_name"`
	BusinessUnit string `json:"business_unit"`
	Policy       string `json:"policy"`
	Teams        string `json:"teams"`
}

// ApplicationList represents the response from the applications API
type ApplicationList struct {
	Applications []Application `json:"_embedded,omitempty"`
}

// doRequest performs an authenticated HTTP request
func (c *Client) doRequest(method, urlPath string) ([]byte, error) {
	return c.doRequestWithBaseURL(method, BaseAPIURL+urlPath)
}

// DoRequestWithQueryParams performs an authenticated HTTP request with query parameters
// This is used by the service layer for the new REST APIs
func (c *Client) DoRequestWithQueryParams(method, urlPath string, params url.Values) ([]byte, error) {
	fullURL := AppSecAPIURL + urlPath
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	body, err := c.doRequestWithBaseURL(method, fullURL)
	if err != nil {
		// Add URL details to error for debugging
		return nil, fmt.Errorf("%w (URL: %s)", err, fullURL)
	}
	return body, nil
}

// DoRequestWithBody performs an authenticated HTTP request with a JSON body and query parameters
// This is used for POST/PUT/PATCH requests that need to send data
func (c *Client) DoRequestWithBody(method, urlPath string, body []byte, params url.Values) ([]byte, error) {
	fullURL := AppSecAPIURL + urlPath
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	respBody, err := c.doRequestWithBodyAndBaseURL(method, fullURL, body)
	if err != nil {
		// Add URL details to error for debugging
		return nil, fmt.Errorf("%w (URL: %s)", err, fullURL)
	}
	return respBody, nil
}

// doRequestWithBaseURL performs an authenticated HTTP request with a full URL
func (c *Client) doRequestWithBaseURL(method, fullURL string) ([]byte, error) {
	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Generate authentication header
	authHeader, err := GenerateAuthHeader(c.apiKeyID, c.apiKeySecret, method, fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth header: %w", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")

	// Log request if debug logging is enabled
	if c.debugLogger != nil {
		c.debugLogger.Printf("\n>>> REQUEST: %s %s\n", method, fullURL)
		c.debugLogger.Printf(">>> Headers: %v\n", req.Header)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && c.debugLogger != nil {
			c.debugLogger.Printf("Warning: failed to close response body: %v", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response if debug logging is enabled
	if c.debugLogger != nil {
		c.debugLogger.Printf("<<< RESPONSE: Status %d\n", resp.StatusCode)
		c.debugLogger.Printf("<<< Headers: %v\n", resp.Header)
		c.debugLogger.Printf("<<< Body: %s\n", string(body))
		c.debugLogger.Println("---")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	return body, nil
}

// doRequestWithBodyAndBaseURL performs an authenticated HTTP request with a full URL and request body
func (c *Client) doRequestWithBodyAndBaseURL(method, fullURL string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, fullURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Generate authentication header
	authHeader, err := GenerateAuthHeader(c.apiKeyID, c.apiKeySecret, method, fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth header: %w", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Log request if debug logging is enabled
	if c.debugLogger != nil {
		c.debugLogger.Printf("\n>>> REQUEST: %s %s\n", method, fullURL)
		c.debugLogger.Printf(">>> Headers: %v\n", req.Header)
		c.debugLogger.Printf(">>> Body: %s\n", string(body))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && c.debugLogger != nil {
			c.debugLogger.Printf("Warning: failed to close response body: %v", closeErr)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response if debug logging is enabled
	if c.debugLogger != nil {
		c.debugLogger.Printf("<<< RESPONSE: Status %d\n", resp.StatusCode)
		c.debugLogger.Printf("<<< Headers: %v\n", resp.Header)
		c.debugLogger.Printf("<<< Body: %s\n", string(respBody))
		c.debugLogger.Println("---")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       respBody,
		}
	}

	return respBody, nil
}

// GetApplications retrieves all applications from Veracode
func (c *Client) GetApplications() ([]Application, error) {
	body, err := c.doRequest("GET", "/5.0/getapplist.do")
	if err != nil {
		return nil, err
	}

	// Parse XML response (Veracode v5 API uses XML)
	// For now, return empty slice - we'll need to parse XML
	// or switch to newer JSON-based APIs
	var apps []Application

	// Try to parse as JSON first (for newer API endpoints)
	var appList ApplicationList
	if err := json.Unmarshal(body, &appList); err == nil {
		return appList.Applications, nil
	}

	// If JSON parsing fails, we have XML which needs different handling
	return apps, nil
}

// GetApplicationByID retrieves a specific application by ID
func (c *Client) GetApplicationByID(appID int) (*Application, error) {
	urlPath := fmt.Sprintf("/5.0/getappinfo.do?app_id=%d", appID)
	body, err := c.doRequest("GET", urlPath)
	if err != nil {
		return nil, err
	}

	var app Application
	if err := json.Unmarshal(body, &app); err != nil {
		return nil, fmt.Errorf("failed to parse application info: %w", err)
	}

	return &app, nil
}

// GetSandboxes retrieves sandboxes for an application
func (c *Client) GetSandboxes(appID int) ([]byte, error) {
	urlPath := fmt.Sprintf("/5.0/getsandboxlist.do?app_id=%d", appID)
	return c.doRequest("GET", urlPath)
}

// GetBuilds retrieves builds for an application
func (c *Client) GetBuilds(appID int) ([]byte, error) {
	urlPath := fmt.Sprintf("/5.0/getbuildlist.do?app_id=%d", appID)
	return c.doRequest("GET", urlPath)
}

// HealthCheck verifies that authentication services are operational
// Returns nil if successful (200 OK), error otherwise
func (c *Client) HealthCheck() error {
	fullURL := HealthCheckAPIURL + "/healthcheck/status"
	_, err := c.doRequestWithBaseURL("GET", fullURL)
	return err
}

// EnableDebugLog enables logging of all REST requests and responses to the specified file
func (c *Client) EnableDebugLog(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open debug log file: %w", err)
	}
	c.debugFile = f
	c.debugLogger = log.New(f, "", log.LstdFlags)
	c.debugLogger.Println("=== Debug logging started ===")
	return nil
}

// Close closes the debug log file if open
func (c *Client) Close() error {
	if c.debugFile != nil {
		return c.debugFile.Close()
	}
	return nil
}
