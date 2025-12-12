// Package annotations provides access to the Veracode Annotations API.
//
// The Annotations API allows you to create annotations (mitigations) for findings
// in your Veracode applications. Annotations can include comments, mitigation actions,
// and apply to one or more findings.
//
// Example usage:
//
//	service := annotations.NewService(client)
//	annotation := &annotations.AnnotationData{
//		IssueList: "123,456,789",
//		Comment:   "This is a false positive",
//		Action:    string(annotations.ActionFalsePositive),
//	}
//	response, err := service.CreateAnnotation(appGUID, annotation, nil)
package annotations
