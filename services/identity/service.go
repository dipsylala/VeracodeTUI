package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dipsylala/veracode-tui/veracode"
)

// Service provides access to the Veracode Identity API
type Service struct {
	client *veracode.Client
}

func NewService(client *veracode.Client) *Service {
	return &Service{
		client: client,
	}
}

// GetPrincipal retrieves the current API user's principal information
func (s *Service) GetPrincipal(ctx context.Context) (*Principal, error) {
	body, err := s.client.DoRequestWithQueryParams("GET", "/api/authn/v2/principal", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	var principal Principal
	if err := json.Unmarshal(body, &principal); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &principal, nil
}

// GetAPICredentials retrieves the current user's API credentials (without the secret)
func (s *Service) GetAPICredentials(ctx context.Context) (*APICredentials, error) {
	body, err := s.client.DoRequestWithQueryParams("GET", "/api/authn/v2/api_credentials", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	var creds APICredentials
	if err := json.Unmarshal(body, &creds); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &creds, nil
}
