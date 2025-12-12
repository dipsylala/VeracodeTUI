package applications

import "time"

// PagedResourceOfApplication represents a paginated list of applications
type PagedResourceOfApplication struct {
	Embedded *EmbeddedApplication `json:"_embedded,omitempty"`
	Links    *Link                `json:"_links,omitempty"`
	Page     *PageMetadata        `json:"page,omitempty"`
}

// EmbeddedApplication contains the applications array
type EmbeddedApplication struct {
	Applications []Application `json:"applications,omitempty"`
}

// Application represents a Veracode application
type Application struct {
	GUID                  string              `json:"guid,omitempty"`
	ID                    int                 `json:"id,omitempty"`
	LegacyID              int                 `json:"legacy_id,omitempty"`
	AppProfileURL         string              `json:"app_profile_url,omitempty"`
	Created               *time.Time          `json:"created,omitempty"`
	Modified              *time.Time          `json:"modified,omitempty"`
	LastCompletedScanDate *time.Time          `json:"last_completed_scan_date,omitempty"`
	OID                   int                 `json:"oid,omitempty"`
	OrganizationID        int                 `json:"organization_id,omitempty"`
	Profile               *ApplicationProfile `json:"profile,omitempty"`
	ResultsURL            string              `json:"results_url,omitempty"`
	Scans                 []ApplicationScan   `json:"scans,omitempty"`
}

// ApplicationProfile contains application profile details
type ApplicationProfile struct {
	Name                string                `json:"name,omitempty"`
	Description         string                `json:"description,omitempty"`
	BusinessCriticality string                `json:"business_criticality,omitempty"`
	BusinessUnit        *BusinessUnit         `json:"business_unit,omitempty"`
	BusinessOwners      []BusinessOwner       `json:"business_owners,omitempty"`
	Policies            []AppPolicy           `json:"policies,omitempty"`
	Teams               []AppTeam             `json:"teams,omitempty"`
	Tags                string                `json:"tags,omitempty"`
	CustomFields        []CustomNameValue     `json:"custom_fields,omitempty"`
	CustomFieldValues   []AppCustomFieldValue `json:"custom_field_values,omitempty"`
	Settings            *ApplicationSettings  `json:"settings,omitempty"`
	ArcherAppName       string                `json:"archer_app_name,omitempty"`
	GitRepoURL          string                `json:"git_repo_url,omitempty"`
	CustomKMSAlias      string                `json:"custom_kms_alias,omitempty"`
}

// ApplicationScan represents scan information
type ApplicationScan struct {
	// Original fields
	ScanType       string     `json:"scan_type,omitempty"`
	Status         string     `json:"status,omitempty"`
	InternalStatus string     `json:"internal_status,omitempty"`
	ModifiedDate   *time.Time `json:"modified_date,omitempty"`
	ScanURL        string     `json:"scan_url,omitempty"`

	// Additional fields from /scans endpoint
	GUID            string                 `json:"guid,omitempty"`
	ID              int64                  `json:"id,omitempty"`
	AnalysisID      int64                  `json:"analysis_id,omitempty"`
	AppVerID        int64                  `json:"app_ver_id,omitempty"`
	ApplicationGUID string                 `json:"application_guid,omitempty"`
	ApplicationID   int                    `json:"application_id,omitempty"`
	Context         string                 `json:"context,omitempty"`
	SandboxGUID     string                 `json:"sandbox_guid,omitempty"`
	SandboxID       int                    `json:"sandbox_id,omitempty"`
	ScanGUID        string                 `json:"scan_guid,omitempty"`
	ScanID          string                 `json:"scan_id,omitempty"`
	PublishedDate   *time.Time             `json:"published_date,omitempty"`
	Deleted         bool                   `json:"deleted,omitempty"`
	DisplayStatus   map[string]interface{} `json:"display_status,omitempty"`
}

// ApplicationSettings contains application settings
type ApplicationSettings struct {
	ScaEnabled                     bool `json:"sca_enabled,omitempty"`
	DynamicScanApprovalNotRequired bool `json:"dynamic_scan_approval_not_required,omitempty"`
	NextdayConsultationAllowed     bool `json:"nextday_consultation_allowed,omitempty"`
	StaticScanDependenciesAllowed  bool `json:"static_scan_dependencies_allowed,omitempty"`
}

// BusinessUnit represents a business unit
type BusinessUnit struct {
	GUID string `json:"guid,omitempty"`
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// BusinessOwner represents a business owner
type BusinessOwner struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// AppPolicy represents an application policy
type AppPolicy struct {
	GUID                   string `json:"guid,omitempty"`
	Name                   string `json:"name,omitempty"`
	IsDefault              bool   `json:"is_default,omitempty"`
	PolicyComplianceStatus string `json:"policy_compliance_status,omitempty"`
}

// AppTeam represents an application team
type AppTeam struct {
	GUID     string `json:"guid,omitempty"`
	TeamID   int    `json:"team_id,omitempty"`
	TeamName string `json:"team_name,omitempty"`
}

// CustomNameValue represents a custom field
type CustomNameValue struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// AppCustomFieldValue represents a custom field value
type AppCustomFieldValue struct {
	ID                 int                 `json:"id,omitempty"`
	FieldNameID        int                 `json:"field_name_id,omitempty"`
	Value              string              `json:"value,omitempty"`
	Created            *time.Time          `json:"created,omitempty"`
	AppCustomFieldName *AppCustomFieldName `json:"app_custom_field_name,omitempty"`
}

// AppCustomFieldName represents a custom field name
type AppCustomFieldName struct {
	ID             int        `json:"id,omitempty"`
	Name           string     `json:"name,omitempty"`
	OrganizationID int        `json:"organization_id,omitempty"`
	SortOrder      int        `json:"sort_order,omitempty"`
	Created        *time.Time `json:"created,omitempty"`
	Modified       *time.Time `json:"modified,omitempty"`
}

// PageMetadata represents pagination metadata
type PageMetadata struct {
	Number        int64 `json:"number,omitempty"`
	Size          int64 `json:"size,omitempty"`
	TotalElements int64 `json:"total_elements,omitempty"`
	TotalPages    int64 `json:"total_pages,omitempty"`
}

// Link represents a hypermedia link
type Link struct {
	Href        string `json:"href,omitempty"`
	Rel         string `json:"rel,omitempty"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	Templated   bool   `json:"templated,omitempty"`
	Deprecation string `json:"deprecation,omitempty"`
	Hreflang    string `json:"hreflang,omitempty"`
	Media       string `json:"media,omitempty"`
}

// PagedResourceOfSandbox represents a paginated list of sandboxes
type PagedResourceOfSandbox struct {
	Embedded *EmbeddedSandbox `json:"_embedded,omitempty"`
	Links    *Link            `json:"_links,omitempty"`
	Page     *PageMetadata    `json:"page,omitempty"`
}

// EmbeddedSandbox contains the sandboxes array
type EmbeddedSandbox struct {
	Sandboxes []Sandbox `json:"sandboxes,omitempty"`
}

// Sandbox represents a Veracode sandbox
type Sandbox struct {
	GUID            string            `json:"guid,omitempty"`
	ID              int               `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	ApplicationGUID string            `json:"application_guid,omitempty"`
	OrganizationID  int               `json:"organization_id,omitempty"`
	OwnerUsername   string            `json:"owner_username,omitempty"`
	AutoRecreate    bool              `json:"auto_recreate,omitempty"`
	CustomFields    []CustomNameValue `json:"custom_fields,omitempty"`
	Created         *time.Time        `json:"created,omitempty"`
	Modified        *time.Time        `json:"modified,omitempty"`
}

// PagedResourceOfScan represents a paginated list of scans
type PagedResourceOfScan struct {
	Embedded *EmbeddedScan `json:"_embedded,omitempty"`
	Links    *Link         `json:"_links,omitempty"`
	Page     *PageMetadata `json:"page,omitempty"`
}

// EmbeddedScan contains the scans array
type EmbeddedScan struct {
	Scans []ApplicationScan `json:"scans,omitempty"`
}
