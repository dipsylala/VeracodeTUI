package ui

// Theme defines the color scheme for the TUI
type Theme struct {
	// Text colors
	DefaultText   string
	SecondaryText string
	DimmedText    string

	// Label and header colors
	Label        string
	ColumnHeader string
	Separator    string

	// Status and severity colors
	Error   string
	Warning string
	Info    string
	Success string
	InfoAlt string

	// Interactive element colors
	New      string
	Approved string
	Rejected string
	Pending  string

	// UI component colors
	Border                     string
	BorderFocused              string
	SelectionBackground        string
	SelectionForeground        string
	DropDownBackground         string
	DropDownText               string
	DropDownSelectedBackground string
	DropDownSelectedForeground string

	// Severity level colors
	SeverityVeryHigh string
	SeverityHigh     string
	SeverityMedium   string
	SeverityLow      string
	SeverityVeryLow  string
	SeverityDefault  string

	// Policy compliance colors
	PolicyPass    string
	PolicyFail    string
	PolicyNeutral string
}

//nolint:dupl // Theme functions have structural duplication - each theme defines all color fields
func DefaultTheme() *Theme {
	return &Theme{
		// Text colors
		DefaultText:   "#FFFFFF",
		SecondaryText: "#888888",
		DimmedText:    "#666666",

		// Label and header colors
		Label:        "#AAAAAA",
		ColumnHeader: "#3B78FF",
		Separator:    "#444444",

		// Status and severity colors
		Error:   "#FF0000",
		Warning: "#FF6600",
		Info:    "#FFFF00",
		Success: "#00FF00",
		InfoAlt: "#00FFFF",

		// Interactive element colors
		New:      "#FFFF00",
		Approved: "#00FF00",
		Rejected: "#FF0000",
		Pending:  "#FFFF00",

		// UI component colors
		Border:                     "#666666",
		BorderFocused:              "#3B78FF",
		SelectionBackground:        "#3B78FF",
		SelectionForeground:        "#FFFFFF",
		DropDownBackground:         "#3B78FF",
		DropDownText:               "#FFFF00",
		DropDownSelectedBackground: "#FFFF00",
		DropDownSelectedForeground: "#000000",

		// Severity level colors
		SeverityVeryHigh: "#FF0000",
		SeverityHigh:     "#FF6600",
		SeverityMedium:   "#FFFF00",
		SeverityLow:      "#00FFFF",
		SeverityVeryLow:  "#888888",
		SeverityDefault:  "#FFFFFF",

		// Policy compliance colors
		PolicyPass:    "#00FF00",
		PolicyFail:    "#FF0000",
		PolicyNeutral: "#FFFFFF",
	}
}

func MonochromeTheme() *Theme {
	// Use simple grayscale hex values
	white := "#FFFFFF"
	gray := "#808080"
	darkGray := "#404040"

	return &Theme{
		// Text colors - use different shades of gray
		DefaultText:   white,
		SecondaryText: gray,
		DimmedText:    darkGray,

		// Label and header colors
		Label:        gray,
		ColumnHeader: white,
		Separator:    darkGray,

		// Status and severity colors - all white, rely on symbols
		Error:   white,
		Warning: white,
		Info:    white,
		Success: white,
		InfoAlt: gray,

		// Interactive element colors
		New:      white,
		Approved: white,
		Rejected: white,
		Pending:  gray,

		// UI component colors
		Border:                     gray,
		BorderFocused:              white,
		SelectionBackground:        gray,
		SelectionForeground:        white,
		DropDownBackground:         darkGray,
		DropDownText:               white,
		DropDownSelectedBackground: white,
		DropDownSelectedForeground: darkGray,

		// Severity level colors - use shades to differentiate
		SeverityVeryHigh: white,
		SeverityHigh:     white,
		SeverityMedium:   gray,
		SeverityLow:      gray,
		SeverityVeryLow:  darkGray,
		SeverityDefault:  white,

		// Policy compliance colors
		PolicyPass:    white,
		PolicyFail:    white,
		PolicyNeutral: gray,
	}
}

//nolint:dupl // Theme functions have structural duplication - each theme defines all color fields
func HotdogTheme() *Theme {
	return &Theme{
		// Text colors - yellow for main text
		DefaultText:   "#FFFF00",
		SecondaryText: "#FFD700",
		DimmedText:    "#FFA500",

		// Label and header colors - red accents
		Label:        "#FF6347",
		ColumnHeader: "#FF4500",
		Separator:    "#8B0000",

		// Status and severity colors - hot colors
		Error:   "#FF0000",
		Warning: "#FF6347",
		Info:    "#FFD700",
		Success: "#32CD32",
		InfoAlt: "#FFA500",

		// Interactive element colors
		New:      "#FFFF00",
		Approved: "#32CD32",
		Rejected: "#DC143C",
		Pending:  "#FFA500",

		// UI component colors - red borders, yellow highlights
		Border:                     "#DC143C",
		BorderFocused:              "#FF0000",
		SelectionBackground:        "#DC143C",
		SelectionForeground:        "#FFFF00",
		DropDownBackground:         "#8B0000",
		DropDownText:               "#FFFF00",
		DropDownSelectedBackground: "#FFFF00",
		DropDownSelectedForeground: "#8B0000",

		// Severity level colors - graduated from red to yellow
		SeverityVeryHigh: "#FF0000",
		SeverityHigh:     "#FF4500",
		SeverityMedium:   "#FF6347",
		SeverityLow:      "#FFA500",
		SeverityVeryLow:  "#FFD700",
		SeverityDefault:  "#FFFF00",

		// Policy compliance colors
		PolicyPass:    "#32CD32",
		PolicyFail:    "#FF0000",
		PolicyNeutral: "#FFFF00",
	}
}

//nolint:dupl // Theme functions have structural duplication - each theme defines all color fields
func MatrixTheme() *Theme {
	return &Theme{
		// Text colors - various shades of Matrix green
		DefaultText:   "#00FF00", // Bright green
		SecondaryText: "#00AA00", // Medium green
		DimmedText:    "#006600", // Dark green

		// Label and header colors - cyan/bright green accents
		Label:        "#00FF00",
		ColumnHeader: "#00FFFF", // Cyan for headers
		Separator:    "#003300",

		// Status and severity colors - Matrix palette
		Error:   "#FF0000", // Red for errors (anomaly in the Matrix)
		Warning: "#FFFF00", // Yellow warnings
		Info:    "#00FFFF", // Cyan info
		Success: "#00FF00", // Bright green success
		InfoAlt: "#00FF7F", // Spring green

		// Interactive element colors
		New:      "#00FFFF", // Cyan for new items
		Approved: "#00FF00", // Green for approved
		Rejected: "#FF0000", // Red for rejected
		Pending:  "#ADFF2F", // Yellow-green for pending

		// UI component colors - dark green borders with bright green focus
		Border:                     "#004400",
		BorderFocused:              "#00FF00",
		SelectionBackground:        "#003300",
		SelectionForeground:        "#00FF00",
		DropDownBackground:         "#002200",
		DropDownText:               "#00FF00",
		DropDownSelectedBackground: "#00FF00",
		DropDownSelectedForeground: "#000000",

		// Severity level colors - green gradient with red/yellow for high severity
		SeverityVeryHigh: "#FF0000", // Red (break from Matrix for urgency)
		SeverityHigh:     "#FF8800", // Orange
		SeverityMedium:   "#FFFF00", // Yellow
		SeverityLow:      "#7FFF00", // Chartreuse
		SeverityVeryLow:  "#00AA00", // Medium green
		SeverityDefault:  "#00FF00", // Bright green

		// Policy compliance colors
		PolicyPass:    "#00FF00", // Green for pass
		PolicyFail:    "#FF0000", // Red for fail
		PolicyNeutral: "#00FFFF", // Cyan for neutral
	}
}
