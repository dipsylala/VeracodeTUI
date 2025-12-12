package ui

import (
	"github.com/dipsylala/veracode-tui/services/annotations"
	"github.com/dipsylala/veracode-tui/services/applications"
	"github.com/dipsylala/veracode-tui/services/findings"
	"github.com/dipsylala/veracode-tui/services/identity"
	"github.com/rivo/tview"
)

// UI represents the TUI application
type UI struct {
	app                *tview.Application
	pages              *tview.Pages
	appService         *applications.Service
	findingsService    *findings.Service
	identityService    *identity.Service
	annotationsService *annotations.Service
	theme              *Theme

	// Data
	applications           []applications.Application
	filteredApps           []applications.Application
	currentPage            int
	totalPages             int
	totalApps              int
	pageSize               int
	searchQuery            string
	selectedApp            *applications.Application
	sandboxes              []applications.Sandbox
	selectionIndex         int // -1 for policy, 0+ for sandbox index
	findings               []findings.Finding
	findingsScanFilter     findings.ScanFilterType
	findingsSeverityFilter int // 0-5, 0 means no filter
	findingsPolicyFilter   findings.PolicyFilterType
	selectedFinding        *findings.Finding
	staticCount            int64
	dynamicCount           int64
	scaCount               int64
	scaExpandedComponents  map[string]bool // Tracks which SCA components are expanded

	// Data path navigation
	currentStaticFlawInfo *findings.StaticFlawInfo
	currentDataPathIndex  int
	currentDataPathsView  *tview.TextView

	// Views - Applications List
	applicationsTable *tview.Table
	statusBar         *tview.TextView
	searchInput       *tview.InputField

	// Views - Application Detail
	detailFlex      *tview.Flex
	appInfoView     *tview.TextView
	complianceView  *tview.TextView
	recentScansView *tview.TextView
	contextsTable   *tview.Table

	// Views - Findings
	findingsTable                  *tview.Table
	findingsFilter                 *tview.DropDown
	findingsSeverityFilterDropdown *tview.DropDown
	findingsPolicyFilterDropdown   *tview.DropDown
	findingsCountsLabel            *tview.TextView
	findingsTitleView              *tview.TextView
	findingsFlex                   *tview.Flex
	findingDetailView              tview.Primitive
	findingAnnotationsView         *tview.TextView // Annotations view in finding detail
}

func NewUI(appService *applications.Service, findingsService *findings.Service, identityService *identity.Service, annotationsService *annotations.Service, theme *Theme) *UI {
	if theme == nil {
		theme = DefaultTheme()
	}

	ui := &UI{
		app:                    tview.NewApplication(),
		pages:                  tview.NewPages(),
		appService:             appService,
		findingsService:        findingsService,
		identityService:        identityService,
		annotationsService:     annotationsService,
		theme:                  theme,
		findingsScanFilter:     "STATIC",
		findingsSeverityFilter: 0,
		findingsPolicyFilter:   findings.PolicyFilterAll,
		currentPage:            0,
		pageSize:               100,
		scaExpandedComponents:  make(map[string]bool),
	}

	ui.setupApplicationsView()

	return ui
}

func (ui *UI) Run() error {
	// Enable mouse support for scrolling and focus
	ui.app.EnableMouse(true)

	// Load initial data
	go ui.loadApplications()

	// Set root and run
	ui.app.SetRoot(ui.pages, true)
	return ui.app.Run()
}
