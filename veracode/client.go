package veracode

import (
	"bytes"
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
	BaseWebURL = "https://analysiscenter.veracode.com/"
	BaseAPIURL = "https://api.veracode.com"
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

// DoRequestWithQueryParams performs an authenticated HTTP request with query parameters
// This is used by the service layer for the new REST APIs
func (c *Client) DoRequestWithQueryParams(method, urlPath string, params url.Values) ([]byte, error) {
	fullURL := BaseAPIURL + urlPath
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
	fullURL := BaseAPIURL + urlPath
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

// HealthCheck verifies that authentication services are operational
// Returns nil if successful (200 OK), error otherwise
func (c *Client) HealthCheck() error {
	fullURL := BaseAPIURL + "/healthcheck/status"
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
