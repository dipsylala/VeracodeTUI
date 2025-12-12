package applications_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/dipsylala/veracode-tui/config"
	"github.com/dipsylala/veracode-tui/services/applications"
	"github.com/dipsylala/veracode-tui/veracode"
)

// Integration tests for Applications Service
// These tests require valid credentials in ~/.veracode/veracode.yml

func setupService(t *testing.T) *applications.Service {
	t.Helper()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return nil
	}

	// Get API credentials
	keyID, keySecret := cfg.GetAPICredentials()

	// Create Veracode API client
	client := veracode.NewClient(keyID, keySecret)

	// Create Applications service
	return applications.NewService(client)
}

func TestGetApplications(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := setupService(t)
	if service == nil {
		return
	}

	t.Run("GetAllApplications", func(t *testing.T) {
		result, err := service.GetApplications(nil)
		if err != nil {
			t.Fatalf("GetApplications failed: %v", err)
		}

		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		t.Logf("Successfully retrieved applications")
		if result.Page != nil {
			t.Logf("Page: %d, Size: %d, Total Elements: %d, Total Pages: %d",
				result.Page.Number, result.Page.Size, result.Page.TotalElements, result.Page.TotalPages)
		}

		if result.Embedded != nil && len(result.Embedded.Applications) > 0 {
			t.Logf("Found %d applications", len(result.Embedded.Applications))
			app := result.Embedded.Applications[0]
			t.Logf("First application: GUID=%s, Name=%s", app.GUID, app.Profile.Name)
		}
	})

	t.Run("GetApplicationsWithPagination", func(t *testing.T) {
		opts := &applications.GetApplicationsOptions{
			Page: 0,
			Size: 10,
		}

		result, err := service.GetApplications(opts)
		if err != nil {
			t.Fatalf("GetApplications with pagination failed: %v", err)
		}

		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		t.Logf("Successfully retrieved applications with pagination (page=0, size=10)")
		if result.Embedded != nil {
			t.Logf("Returned %d applications", len(result.Embedded.Applications))
		}
	})

	t.Run("GetApplicationsByName", func(t *testing.T) {
		// First get all applications to find one we can search for
		allApps, err := service.GetApplications(&applications.GetApplicationsOptions{Size: 1})
		if err != nil {
			t.Fatalf("Failed to get applications for search test: %v", err)
		}

		if allApps.Embedded == nil || len(allApps.Embedded.Applications) == 0 {
			t.Skip("No applications available to test name search")
			return
		}

		appName := allApps.Embedded.Applications[0].Profile.Name
		opts := &applications.GetApplicationsOptions{
			Name: appName,
		}

		result, err := service.GetApplications(opts)
		if err != nil {
			t.Fatalf("GetApplications by name failed: %v", err)
		}

		if result.Embedded == nil || len(result.Embedded.Applications) == 0 {
			t.Errorf("Expected to find application with name '%s'", appName)
		} else {
			t.Logf("Successfully found application by name: %s", appName)
		}
	})
}

func TestGetApplication(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := setupService(t)
	if service == nil {
		return
	}

	t.Run("GetApplicationByGUID", func(t *testing.T) {
		// First get an application GUID
		apps, err := service.GetApplications(&applications.GetApplicationsOptions{Size: 1})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available to test GetApplication")
			return
		}

		appGUID := apps.Embedded.Applications[0].GUID

		// Get the specific application
		app, err := service.GetApplication(appGUID)
		if err != nil {
			t.Fatalf("GetApplication failed: %v", err)
		}

		if app == nil {
			t.Fatal("Expected non-nil application")
		}

		if app.GUID != appGUID {
			t.Errorf("Expected GUID %s, got %s", appGUID, app.GUID)
		}

		t.Logf("Successfully retrieved application: GUID=%s, Name=%s", app.GUID, app.Profile.Name)
		if app.Profile.BusinessUnit != nil {
			t.Logf("Business Unit: %s", app.Profile.BusinessUnit.Name)
		}
		if len(app.Profile.Policies) > 0 {
			t.Logf("Policies: %d", len(app.Profile.Policies))
		}

		// Log scan information
		if len(app.Scans) > 0 {
			t.Logf("Scans found: %d", len(app.Scans))
			for i, scan := range app.Scans {
				t.Logf("Scan %d:", i+1)
				t.Logf("  Type: %s", scan.ScanType)
				t.Logf("  Status: %s", scan.Status)
				if scan.InternalStatus != "" {
					t.Logf("  Internal Status: %s", scan.InternalStatus)
				}
				if scan.ModifiedDate != nil {
					t.Logf("  Modified: %s", scan.ModifiedDate.Format("2006-01-02 15:04:05"))
				}
				if scan.ScanURL != "" {
					t.Logf("  URL: %s", scan.ScanURL)
				}
			}
		} else {
			t.Log("No scans found for this application")
		}
	})

	t.Run("GetApplicationInvalidGUID", func(t *testing.T) {
		_, err := service.GetApplication("invalid-guid-12345")
		if err == nil {
			t.Error("Expected error for invalid GUID, got nil")
		} else {
			t.Logf("Got expected error for invalid GUID: %v", err)
		}
	})
}

func TestGetSandboxes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := setupService(t)
	if service == nil {
		return
	}

	t.Run("GetSandboxesForApplication", func(t *testing.T) {
		// First get an application GUID
		apps, err := service.GetApplications(&applications.GetApplicationsOptions{Size: 1})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available to test GetSandboxes")
			return
		}

		appGUID := apps.Embedded.Applications[0].GUID

		// Get sandboxes for the application
		result, err := service.GetSandboxes(appGUID, nil)
		if err != nil {
			t.Fatalf("GetSandboxes failed: %v", err)
		}

		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		t.Logf("Successfully retrieved sandboxes for application %s", appGUID)
		if result.Embedded != nil && len(result.Embedded.Sandboxes) > 0 {
			t.Logf("Found %d sandboxes", len(result.Embedded.Sandboxes))
			sandbox := result.Embedded.Sandboxes[0]
			t.Logf("First sandbox: GUID=%s, Name=%s", sandbox.GUID, sandbox.Name)
		} else {
			t.Log("No sandboxes found for this application")
		}
	})
}

func TestGetSandbox(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := setupService(t)
	if service == nil {
		return
	}

	t.Run("GetSandboxByGUID", func(t *testing.T) {
		// First get an application with sandboxes
		apps, err := service.GetApplications(&applications.GetApplicationsOptions{Size: 10})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available to test GetSandbox")
			return
		}

		// Find an application with sandboxes
		var appGUID, sandboxGUID string
		for _, app := range apps.Embedded.Applications {
			sandboxes, err := service.GetSandboxes(app.GUID, nil)
			if err != nil {
				continue
			}
			if sandboxes.Embedded != nil && len(sandboxes.Embedded.Sandboxes) > 0 {
				appGUID = app.GUID
				sandboxGUID = sandboxes.Embedded.Sandboxes[0].GUID
				break
			}
		}

		if sandboxGUID == "" {
			t.Skip("No sandboxes found to test GetSandbox")
			return
		}

		// Get the specific sandbox
		sandbox, err := service.GetSandbox(appGUID, sandboxGUID)
		if err != nil {
			t.Fatalf("GetSandbox failed: %v", err)
		}

		if sandbox == nil {
			t.Fatal("Expected non-nil sandbox")
		}

		if sandbox.GUID != sandboxGUID {
			t.Errorf("Expected sandbox GUID %s, got %s", sandboxGUID, sandbox.GUID)
		}

		t.Logf("Successfully retrieved sandbox: GUID=%s, Name=%s", sandbox.GUID, sandbox.Name)
	})
}

func TestApplicationScans(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := setupService(t)
	if service == nil {
		return
	}

	t.Run("InspectApplicationScans", func(t *testing.T) {
		// Get multiple applications to find one with scans
		apps, err := service.GetApplications(&applications.GetApplicationsOptions{Size: 20})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available")
			return
		}

		// Find an application with scans
		var appWithScans *applications.Application
		for i := range apps.Embedded.Applications {
			app, err := service.GetApplication(apps.Embedded.Applications[i].GUID)
			if err != nil {
				continue
			}
			if len(app.Scans) > 0 {
				appWithScans = app
				break
			}
		}

		if appWithScans == nil {
			t.Skip("No applications with scans found")
			return
		}

		t.Logf("=== Application with Scans ===")
		t.Logf("Name: %s", appWithScans.Profile.Name)
		t.Logf("GUID: %s", appWithScans.GUID)
		t.Logf("Total Scans: %d", len(appWithScans.Scans))
		t.Logf("")

		// Display detailed information about each scan
		for i, scan := range appWithScans.Scans {
			t.Logf("=== Scan %d ===", i+1)
			t.Logf("  ScanType:       %q", scan.ScanType)
			t.Logf("  Status:         %q", scan.Status)
			t.Logf("  InternalStatus: %q", scan.InternalStatus)
			if scan.ModifiedDate != nil {
				t.Logf("  ModifiedDate:   %s", scan.ModifiedDate.Format("2006-01-02 15:04:05 MST"))
			} else {
				t.Logf("  ModifiedDate:   <nil>")
			}
			t.Logf("  ScanURL:        %q", scan.ScanURL)
			t.Logf("")
		}

		// Group scans by type
		scansByType := make(map[string]int)
		for _, scan := range appWithScans.Scans {
			scansByType[scan.ScanType]++
		}

		t.Log("=== Scans by Type ===")
		for scanType, count := range scansByType {
			t.Logf("  %s: %d", scanType, count)
		}
	})
}

func TestApplicationScansRawJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	keyID, keySecret := cfg.GetAPICredentials()
	client := veracode.NewClient(keyID, keySecret)

	t.Run("ShowRawJSONResponse", func(t *testing.T) {
		// Get multiple applications to find one with scans
		service := applications.NewService(client)
		apps, err := service.GetApplications(&applications.GetApplicationsOptions{Size: 20})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available")
			return
		}

		// Find an application with scans
		var appGUID, appName string
		for i := range apps.Embedded.Applications {
			app, err := service.GetApplication(apps.Embedded.Applications[i].GUID)
			if err != nil {
				continue
			}
			if len(app.Scans) > 0 {
				appGUID = app.GUID
				appName = app.Profile.Name
				t.Logf("Found application with scans: %s (GUID: %s)", appName, appGUID)
				break
			}
		}

		if appGUID == "" {
			t.Skip("No applications with scans found")
			return
		}

		// Make raw API call to get JSON response
		urlPath := fmt.Sprintf("/appsec/v1/applications/%s", appGUID)
		bodyBytes, err := client.DoRequestWithQueryParams("GET", urlPath, nil)
		if err != nil {
			t.Fatalf("Failed to make API request: %v", err)
		}

		t.Logf("\n=== RAW JSON RESPONSE for %s ===\n", appName)

		// Pretty print the JSON
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &prettyJSON); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		prettyBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
		if err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}

		t.Logf("%s\n", string(prettyBytes))
	})
}

// Example of running integration tests manually
func ExampleService() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Get API credentials
	keyID, keySecret := cfg.GetAPICredentials()

	// Create client and service
	client := veracode.NewClient(keyID, keySecret)
	service := applications.NewService(client)

	// Get applications
	apps, err := service.GetApplications(&applications.GetApplicationsOptions{
		Size: 10,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d applications\n", len(apps.Embedded.Applications))

	// Get first application details
	if len(apps.Embedded.Applications) > 0 {
		app, err := service.GetApplication(apps.Embedded.Applications[0].GUID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Application: %s\n", app.Profile.Name)
	}
}
