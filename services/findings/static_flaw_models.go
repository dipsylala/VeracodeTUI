package findings

// StaticFlawInfo represents the detailed data path information for a static flaw
type StaticFlawInfo struct {
	IssueSummary *IssueSummary `json:"issue_summary,omitempty"`
	DataPaths    []DataPath    `json:"data_paths,omitempty"`
}

// IssueSummary contains basic information about the flaw
type IssueSummary struct {
	AppGUID string `json:"app_guid,omitempty"`
	Name    string `json:"name,omitempty"`
	BuildID int    `json:"build_id,omitempty"`
	IssueID int    `json:"issue_id,omitempty"`
	Context string `json:"context,omitempty"`
}

// DataPath represents a call stack for the static flaw
type DataPath struct {
	ModuleName   string `json:"module_name,omitempty"`
	Steps        int    `json:"steps,omitempty"`
	LocalPath    string `json:"local_path,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	LineNumber   int    `json:"line_number,omitempty"`
	Calls        []Call `json:"calls,omitempty"`
}

// Call represents a single call in the data path
type Call struct {
	DataPath     int    `json:"data_path,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	LineNumber   int    `json:"line_number,omitempty"`
}
