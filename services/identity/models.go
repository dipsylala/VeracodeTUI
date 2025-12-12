package identity

import "time"

// Principal represents the current API user's information
type Principal struct {
	Email            string   `json:"email"`
	Features         []string `json:"features"`
	LearnGroupID     int      `json:"learnGroupId"`
	LearnTrackID     int      `json:"learnTrackId"`
	OrganizationID   int      `json:"organizationId"`
	OrganizationName string   `json:"organizationName"`
	OrganizationUUID string   `json:"organizationUuid"`
	Permissions      []string `json:"permissions"`
	PinRequired      bool     `json:"pinRequired"`
	Roles            []string `json:"roles"`
	SAMLUser         bool     `json:"samlUser"`
	SandboxEnabled   bool     `json:"sandboxEnabled"`
	UserFirstName    string   `json:"userFirstName"`
	UserID           int      `json:"userId"`
	UserLastName     string   `json:"userLastName"`
	UserUUID         string   `json:"userUuid"`
	Username         string   `json:"username"`
}

// APICredentials represents API credential information
type APICredentials struct {
	APIID          string    `json:"api_id"`
	APISecret      string    `json:"api_secret,omitempty"`
	ExpirationTS   time.Time `json:"expiration_ts"`
	OrgID          string    `json:"org_id,omitempty"`
	RevocationTS   time.Time `json:"revocation_ts,omitempty"`
	RevocationUser string    `json:"revocation_user,omitempty"`
	UserID         string    `json:"user_id,omitempty"`
}
