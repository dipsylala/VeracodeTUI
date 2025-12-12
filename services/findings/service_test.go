package findings_test

import (
	"fmt"
	"testing"

	"github.com/dipsylala/veracode-tui/config"
	"github.com/dipsylala/veracode-tui/services/applications"
	"github.com/dipsylala/veracode-tui/services/findings"
	"github.com/dipsylala/veracode-tui/veracode"
)

// Integration tests for Findings Service
// These tests require valid credentials in ~/.veracode/veracode.yml

func setupServices(t *testing.T) (*applications.Service, *findings.Service) {
	t.Helper()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return nil, nil
	}

	// Get API credentials
	keyID, keySecret := cfg.GetAPICredentials()

	// Create Veracode API client
	client := veracode.NewClient(keyID, keySecret)

	// Create services
	return applications.NewService(client), findings.NewService(client)
}

func TestGetFindings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	appService, findingsService := setupServices(t)
	if appService == nil || findingsService == nil {
		return
	}

	t.Run("InspectFindingsRequest", func(t *testing.T) {
		// Get an application to test with
		apps, err := appService.GetApplications(&applications.GetApplicationsOptions{Size: 1})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available")
			return
		}

		appGUID := apps.Embedded.Applications[0].GUID
		appName := apps.Embedded.Applications[0].Profile.Name

		t.Logf("Testing with application: %s (GUID: %s)", appName, appGUID)

		// Test 1: No context parameter (should return policy findings)
		t.Log("\n=== Test 1: No context parameter (policy) ===")
		result1, err1 := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
			Size: 10,
		})
		if err1 != nil {
			t.Logf("Error (no context): %v", err1)
		} else if result1 != nil && result1.Embedded != nil {
			t.Logf("Success: Found %d findings", len(result1.Embedded.Findings))
		}

		// Test 2: Empty context (explicitly policy)
		t.Log("\n=== Test 2: Empty context string ===")
		result2, err2 := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
			Context: "",
			Size:    10,
		})
		if err2 != nil {
			t.Logf("Error (empty context): %v", err2)
		} else if result2 != nil && result2.Embedded != nil {
			t.Logf("Success: Found %d findings", len(result2.Embedded.Findings))
		}
	})

	t.Run("GetPolicyFindings", func(t *testing.T) {
		// Get an application
		apps, err := appService.GetApplications(&applications.GetApplicationsOptions{Size: 1})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available")
			return
		}

		appGUID := apps.Embedded.Applications[0].GUID
		appName := apps.Embedded.Applications[0].Profile.Name

		// Get findings for the application (policy context)
		result, err := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
			Context: "",
			Size:    100,
		})
		if err != nil {
			t.Fatalf("GetFindings failed: %v", err)
		}

		t.Logf("=== Policy Findings for %s ===", appName)
		t.Logf("Application GUID: %s", appGUID)

		if result.Embedded != nil && len(result.Embedded.Findings) > 0 {
			t.Logf("Total Findings: %d", len(result.Embedded.Findings))

			// Show first few findings
			maxShow := 5
			if len(result.Embedded.Findings) < maxShow {
				maxShow = len(result.Embedded.Findings)
			}

			for i := 0; i < maxShow; i++ {
				finding := result.Embedded.Findings[i]
				t.Logf("\nFinding %d:", i+1)
				t.Logf("  Issue ID: %d", finding.IssueID)
				t.Logf("  Scan Type: %s", finding.ScanType)
				t.Logf("  Context Type: %s", finding.ContextType)
				t.Logf("  Violates Policy: %t", finding.ViolatesPolicy)
				if finding.FindingStatus != nil {
					t.Logf("  Status: %s", finding.FindingStatus.Status)
				}
				if len(finding.Description) > 100 {
					t.Logf("  Description: %s...", finding.Description[:100])
				} else {
					t.Logf("  Description: %s", finding.Description)
				}
			}

			// Group by scan type
			scanTypes := make(map[string]int)
			for _, finding := range result.Embedded.Findings {
				scanTypes[string(finding.ScanType)]++
			}

			t.Log("\n=== Findings by Scan Type ===")
			for scanType, count := range scanTypes {
				t.Logf("  %s: %d", scanType, count)
			}
		} else {
			t.Log("No findings found for this application")
		}
	})

	t.Run("GetSandboxFindings", func(t *testing.T) {
		// Get an application with sandboxes
		apps, err := appService.GetApplications(&applications.GetApplicationsOptions{Size: 10})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available")
			return
		}

		// Find an application with sandboxes
		var appGUID, sandboxGUID, sandboxName string
		for _, app := range apps.Embedded.Applications {
			sandboxes, err := appService.GetSandboxes(app.GUID, nil)
			if err != nil {
				continue
			}
			if sandboxes.Embedded != nil && len(sandboxes.Embedded.Sandboxes) > 0 {
				appGUID = app.GUID
				sandboxGUID = sandboxes.Embedded.Sandboxes[0].GUID
				sandboxName = sandboxes.Embedded.Sandboxes[0].Name
				break
			}
		}

		if sandboxGUID == "" {
			t.Skip("No sandboxes found")
			return
		}

		// Get findings for the sandbox context
		result, err := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
			Context: sandboxGUID,
			Size:    100,
		})
		if err != nil {
			t.Fatalf("GetFindings for sandbox failed: %v", err)
		}

		t.Logf("=== Sandbox Findings for %s ===", sandboxName)
		t.Logf("Application GUID: %s", appGUID)
		t.Logf("Sandbox GUID: %s", sandboxGUID)

		if result.Embedded != nil && len(result.Embedded.Findings) > 0 {
			t.Logf("Total Findings: %d", len(result.Embedded.Findings))

			// Group by scan type
			scanTypes := make(map[string]int)
			policyViolations := 0
			for _, finding := range result.Embedded.Findings {
				scanTypes[string(finding.ScanType)]++
				if finding.ViolatesPolicy {
					policyViolations++
				}
			}

			t.Log("\n=== Findings by Scan Type ===")
			for scanType, count := range scanTypes {
				t.Logf("  %s: %d", scanType, count)
			}
			t.Logf("\nPolicy Violations: %d", policyViolations)
		} else {
			t.Log("No findings found for this sandbox")
		}
	})

	t.Run("GetFindingsWithFilters", func(t *testing.T) {
		// Get an application
		apps, err := appService.GetApplications(&applications.GetApplicationsOptions{Size: 1})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available")
			return
		}

		appGUID := apps.Embedded.Applications[0].GUID

		// Test filtering by scan type
		result, err := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
			ScanType: []string{"STATIC"},
			Size:     50,
		})
		if err != nil {
			t.Fatalf("GetFindings with scan type filter failed: %v", err)
		}

		t.Log("=== Findings filtered by STATIC scan type ===")
		if result.Embedded != nil && len(result.Embedded.Findings) > 0 {
			t.Logf("Found %d STATIC findings", len(result.Embedded.Findings))

			// Show first few findings with severity
			maxShow := 5
			if len(result.Embedded.Findings) < maxShow {
				maxShow = len(result.Embedded.Findings)
			}

			for i := 0; i < maxShow; i++ {
				finding := result.Embedded.Findings[i]
				severity := "Unknown"
				if details, ok := finding.FindingDetails.(map[string]interface{}); ok {
					if sev, ok := details["severity"].(float64); ok {
						severity = fmt.Sprintf("%d", int(sev))
					}
				}
				t.Logf("  Finding %d: Severity=%s, IssueID=%d, ScanType=%s",
					i+1, severity, finding.IssueID, finding.ScanType)

				// Verify scan type
				if finding.ScanType != findings.ScanTypeStatic {
					t.Errorf("Expected STATIC scan type, got %s", finding.ScanType)
				}
			}
		} else {
			t.Log("No STATIC findings found")
		}
	})

	t.Run("GetFindingsWithPagination", func(t *testing.T) {
		// Get an application
		apps, err := appService.GetApplications(&applications.GetApplicationsOptions{Size: 1})
		if err != nil {
			t.Fatalf("Failed to get applications: %v", err)
		}

		if apps.Embedded == nil || len(apps.Embedded.Applications) == 0 {
			t.Skip("No applications available")
			return
		}

		appGUID := apps.Embedded.Applications[0].GUID

		// Test with page=0 and size=500 (like the TUI does)
		result, err := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
			ScanType: []string{"STATIC"},
			Page:     0,
			Size:     500,
		})
		if err != nil {
			t.Fatalf("GetFindings with page=0 and size=500 failed: %v", err)
		}

		t.Logf("=== Findings with page=0, size=500 ===")
		if result.Embedded != nil && len(result.Embedded.Findings) > 0 {
			t.Logf("Found %d findings", len(result.Embedded.Findings))
		} else {
			t.Log("No findings found")
		}
	})

	// TestGetSandboxFindingStaticFlawInfo verifies a bug where the static_flaw_info endpoint
	// returns 404 when the 'context' query parameter is provided for sandbox findings.
	//
	// Issue: When fetching data paths for a finding in a sandbox context, the API returns:
	//   HTTP 404: "Build does not have static flaws"
	//
	// Workaround: Omit the 'context' parameter when calling static_flaw_info endpoint.
	// The endpoint returns the correct data without the context filter.
	//
	// Test case uses:
	//   - Application GUID: 304c7929-27f0-4257-90e3-7d9e6cfb4cd3 (MCPApacheSpark)
	//   - Sandbox GUID: 51f65d69-9f52-4976-af23-d24ad709c5b4 (JS Only)
	//   - Finding ID: 134
	t.Run("GetSandboxFindingStaticFlawInfo", func(t *testing.T) {
		// Use the specific application and sandbox GUID reported in the bug
		appGUID := "304c7929-27f0-4257-90e3-7d9e6cfb4cd3"
		sandboxGUID := "51f65d69-9f52-4976-af23-d24ad709c5b4"
		findingID := int64(134)

		t.Logf("=== Testing Static Flaw Info for Sandbox Finding ===")
		t.Logf("Application GUID: %s", appGUID)
		t.Logf("Sandbox GUID: %s", sandboxGUID)
		t.Logf("Finding ID: %d", findingID)

		// First, check if the sandbox exists
		t.Log("\n--- Step 0: Verify sandbox exists ---")
		sandboxes, err := appService.GetSandboxes(appGUID, nil)
		if err != nil {
			t.Fatalf("Failed to get sandboxes: %v", err)
		}

		var sandboxFound bool
		if sandboxes.Embedded != nil {
			for _, sb := range sandboxes.Embedded.Sandboxes {
				t.Logf("Found sandbox: %s (GUID: %s)", sb.Name, sb.GUID)
				if sb.GUID == sandboxGUID {
					sandboxFound = true
					t.Logf("  ✓ Target sandbox found: %s", sb.Name)
				}
			}
		}

		if !sandboxFound {
			t.Logf("WARNING: Target sandbox GUID not found in current sandboxes list")
		}

		// Try to get the finding itself
		t.Log("\n--- Step 1: Get the finding from findings list ---")
		findingsResult, err := findingsService.GetFindings(appGUID, &findings.GetFindingsOptions{
			Context: sandboxGUID,
			Size:    500,
		})
		if err != nil {
			t.Logf("ERROR getting sandbox findings: %v", err)
			t.Log("Trying to get static flaw info anyway to see the exact error...")
		} else if findingsResult != nil && findingsResult.Embedded != nil {
			t.Logf("Got %d findings from sandbox", len(findingsResult.Embedded.Findings))
		}

		// Find the specific finding if we have results
		var targetFinding *findings.Finding
		if findingsResult != nil && findingsResult.Embedded != nil {
			for i := range findingsResult.Embedded.Findings {
				if findingsResult.Embedded.Findings[i].IssueID == findingID {
					targetFinding = &findingsResult.Embedded.Findings[i]
					break
				}
			}
		}

		if targetFinding == nil {
			t.Logf("Could not find issue ID %d in current sandbox findings", findingID)
			t.Log("This might be expected if the finding was removed or the sandbox was cleaned")
		} else {
			t.Logf("Found the target finding:")
			t.Logf("  Issue ID: %d", targetFinding.IssueID)
			t.Logf("  Scan Type: %s", targetFinding.ScanType)
			t.Logf("  Context Type: %s", targetFinding.ContextType)
			t.Logf("  Context GUID: %s", targetFinding.ContextGUID)
			if targetFinding.FindingStatus != nil {
				t.Logf("  Status: %s", targetFinding.FindingStatus.Status)
			}
		}

		// Now try to get the static flaw info
		t.Log("\n--- Step 2: Get static flaw info (data paths) with sandbox context ---")
		staticInfo, err := findingsService.GetStaticFlawInfo(appGUID, findingID, sandboxGUID)
		if err != nil {
			t.Logf("ERROR getting static flaw info with sandbox context: %v", err)
			t.Logf("This confirms the bug - getting static_flaw_info for sandbox findings returns 404")

			// Try without context to see if it's a sandbox-specific issue
			t.Log("\n--- Step 3: Try getting static flaw info WITHOUT context parameter ---")
			staticInfoNoContext, err2 := findingsService.GetStaticFlawInfo(appGUID, findingID, "")
			if err2 != nil {
				t.Logf("ERROR without context: %v", err2)
				t.Logf("The issue affects both policy and sandbox contexts for this finding")
			} else {
				t.Logf("SUCCESS without context!")
				if staticInfoNoContext != nil {
					t.Logf("Got static flaw info: %+v", staticInfoNoContext)
					t.Logf("This suggests the API may not support the 'context' parameter for static_flaw_info endpoint")
				}
			}
		} else {
			t.Logf("SUCCESS: Got static flaw info with sandbox context")
			if staticInfo != nil {
				t.Logf("Static Flaw Info: %+v", staticInfo)
			}
		}
	})
}
