package findings_test

import (
	"encoding/json"
	"testing"

	"github.com/dipsylala/veracode-tui/config"
	"github.com/dipsylala/veracode-tui/services/applications"
	"github.com/dipsylala/veracode-tui/services/findings"
	"github.com/dipsylala/veracode-tui/veracode"
)

// TestMCPVerademoStaticFlaws verifies the 'affects policy' behavior
// by pulling static flaws from the MCPVerademo application
func TestMCPVerademoStaticFlaws(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	// Get API credentials
	keyID, keySecret := cfg.GetAPICredentials()

	// Create Veracode API client
	client := veracode.NewClient(keyID, keySecret)

	// Create services
	appService := applications.NewService(client)
	findingsService := findings.NewService(client)

	// Find the MCPVerademo application
	t.Log("=== Searching for MCPVerademo Application ===")
	apps, err := appService.GetApplications(&applications.GetApplicationsOptions{
		Size: 500, // Get enough apps to find MCPVerademo
	})
	if err != nil {
		t.Fatalf("Failed to get applications: %v", err)
	}

	var mcpVerademoGUID string
	if apps.Embedded != nil {
		for _, app := range apps.Embedded.Applications {
			if app.Profile != nil && app.Profile.Name == "MCPVerademo" {
				mcpVerademoGUID = app.GUID
				t.Logf("Found MCPVerademo: GUID=%s", mcpVerademoGUID)
				break
			}
		}
	}

	if mcpVerademoGUID == "" {
		t.Fatal("MCPVerademo application not found")
	}

	// Get STATIC findings for MCPVerademo
	t.Log("\n=== Fetching STATIC Findings for MCPVerademo ===")
	result, err := findingsService.GetFindings(mcpVerademoGUID, &findings.GetFindingsOptions{
		ScanType: []string{"STATIC"},
		Size:     500,
	})
	if err != nil {
		t.Fatalf("Failed to get findings: %v", err)
	}

	if result.Embedded == nil || len(result.Embedded.Findings) == 0 {
		t.Fatal("No static findings found for MCPVerademo")
	}

	t.Logf("Total STATIC findings: %d", len(result.Embedded.Findings))

	// Analyze findings by policy violation and mitigation status
	var (
		totalFindings      = len(result.Embedded.Findings)
		violatesPolicy     = 0
		notViolatesPolicy  = 0
		approved           = 0
		approvedViolates   = 0
		approvedNotViolate = 0
		closedViolates     = 0
		closedNotViolate   = 0
		proposed           = 0
		rejected           = 0
	)

	t.Log("\n=== Detailed Finding Analysis ===")
	for i, finding := range result.Embedded.Findings {
		// Track policy violations
		if finding.ViolatesPolicy {
			violatesPolicy++
		} else {
			notViolatesPolicy++
		}

		// Track resolution status
		var resolutionStatus findings.ResolutionStatus
		var mitigationReviewStatus findings.ResolutionStatus
		var status findings.Status

		if finding.FindingStatus != nil {
			resolutionStatus = finding.FindingStatus.ResolutionStatus
			mitigationReviewStatus = finding.FindingStatus.MitigationReviewStatus
			status = finding.FindingStatus.Status

			if resolutionStatus == findings.ResolutionApproved {
				approved++
				if finding.ViolatesPolicy {
					approvedViolates++
				} else {
					approvedNotViolate++
				}
			} else if resolutionStatus == findings.ResolutionProposed {
				proposed++
			} else if resolutionStatus == findings.ResolutionRejected {
				rejected++
			}

			if status == findings.StatusClosed {
				if finding.ViolatesPolicy {
					closedViolates++
				} else {
					closedNotViolate++
				}
			}
		}

		// Log first 10 findings with details
		if i < 10 {
			t.Logf("\nFinding #%d:", i+1)
			t.Logf("  Issue ID: %d", finding.IssueID)
			t.Logf("  Violates Policy: %t", finding.ViolatesPolicy)
			t.Logf("  Status: %s", status)
			t.Logf("  Resolution Status: %s", resolutionStatus)
			t.Logf("  Mitigation Review Status: %s", mitigationReviewStatus)

			// Extract severity and CWE
			if details, ok := finding.FindingDetails.(map[string]interface{}); ok {
				if sev, ok := details["severity"].(float64); ok {
					t.Logf("  Severity: %d", int(sev))
				}
				if cweData, ok := details["cwe"].(map[string]interface{}); ok {
					if cweID, ok := cweData["id"].(float64); ok {
						t.Logf("  CWE: %d", int(cweID))
					}
				}
			}

			// Check annotations
			if len(finding.Annotations) > 0 {
				t.Logf("  Annotations: %d", len(finding.Annotations))
				for j, annotation := range finding.Annotations {
					t.Logf("    Annotation %d:", j+1)
					t.Logf("      Action: %s", annotation.Action)
					t.Logf("      User: %s", annotation.User)
					if annotation.Description != "" {
						t.Logf("      Description: %s", annotation.Description)
					}
					if annotation.Comment != "" {
						commentPreview := annotation.Comment
						if len(commentPreview) > 100 {
							commentPreview = commentPreview[:100] + "..."
						}
						t.Logf("      Comment: %s", commentPreview)
					}
				}
			}
		}

		// For Finding #6 (the APPROVED one), show full JSON response
		if finding.IssueID == 6 {
			t.Log("\n=== Full JSON Response for Finding #6 (APPROVED) ===")
			jsonBytes, err := json.MarshalIndent(finding, "  ", "  ")
			if err != nil {
				t.Logf("  Error marshaling JSON: %v", err)
			} else {
				t.Logf("%s", string(jsonBytes))
			}
		}
	}

	// Summary statistics
	t.Log("\n=== Summary Statistics ===")
	t.Logf("Total Findings: %d", totalFindings)
	t.Logf("Violates Policy: %d (%.1f%%)", violatesPolicy, float64(violatesPolicy)/float64(totalFindings)*100)
	t.Logf("Does NOT Violate Policy: %d (%.1f%%)", notViolatesPolicy, float64(notViolatesPolicy)/float64(totalFindings)*100)

	t.Log("\n=== Mitigation Status ===")
	t.Logf("APPROVED Mitigations: %d", approved)
	t.Logf("  - APPROVED + Violates Policy: %d", approvedViolates)
	t.Logf("  - APPROVED + Does NOT Violate Policy: %d", approvedNotViolate)
	t.Logf("PROPOSED Mitigations: %d", proposed)
	t.Logf("REJECTED Mitigations: %d", rejected)

	t.Log("\n=== Status Analysis ===")
	t.Logf("CLOSED Findings: %d", closedViolates+closedNotViolate)
	t.Logf("  - CLOSED + Violates Policy: %d", closedViolates)
	t.Logf("  - CLOSED + Does NOT Violate Policy: %d", closedNotViolate)

	// Validation checks
	t.Log("\n=== Validation Checks ===")

	// Check 1: APPROVED findings should show as mitigated (✓) even if they violate policy
	if approvedViolates > 0 {
		t.Logf("✓ Found %d APPROVED findings that still violate policy", approvedViolates)
		t.Logf("  These should display with ✓ checkmark in the TUI")
	}

	// Check 2: CLOSED findings that don't violate policy are truly mitigated
	if closedNotViolate > 0 {
		t.Logf("✓ Found %d CLOSED findings that no longer violate policy", closedNotViolate)
		t.Logf("  These should display with ✓ checkmark in the TUI")
	}

	// Check 3: Findings that violate policy without approved mitigation should show ❌
	unapprovedViolations := violatesPolicy - approvedViolates
	if unapprovedViolations > 0 {
		t.Logf("✓ Found %d findings that violate policy without APPROVED mitigation", unapprovedViolations)
		t.Logf("  These should display with ❌ in the TUI")
	}

	// Report expected TUI behavior
	t.Log("\n=== Expected TUI Display Behavior ===")
	t.Logf("Total findings that should show ✓ (mitigated): %d", approved+closedNotViolate)
	t.Logf("  - APPROVED resolution: %d", approved)
	t.Logf("  - CLOSED and no policy violation: %d", closedNotViolate)
	t.Logf("Total findings that should show ❌ (violates policy): %d", unapprovedViolations)
	t.Logf("Total findings that should show ' ' (never violated): %d", notViolatesPolicy-closedNotViolate-approvedNotViolate)
}
