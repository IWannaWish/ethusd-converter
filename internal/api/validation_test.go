// server_test.go
package api

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestValidateExample(t *testing.T) {
	// Test case where the name does NOT start with "Report"
	name := "Example"
	var violations []*errdetails.BadRequest_FieldViolation
	violations = ValidateExample(name, violations)

	if len(violations) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(violations))
	}

	// Check if the violation is populated correctly
	expectedField := "name"
	expectedDescription := "name should start with 'Report'"

	if violations[0].Field != expectedField {
		t.Errorf("Expected Field %s, got %s", expectedField, violations[0].Field)
	}

	if violations[0].Description != expectedDescription {
		t.Errorf("Expected Description %s, got %s", expectedDescription, violations[0].Description)
	}

	// Test case where the name starts with "Report" and no violation should occur
	name = "ReportExample"
	violations = []*errdetails.BadRequest_FieldViolation{}
	violations = ValidateExample(name, violations)

	if len(violations) != 0 {
		t.Errorf("Expected 0 violations, got %d", len(violations))
	}
}

func TestCheckViolations(t *testing.T) {
	// Test case with one violation
	violations := []*errdetails.BadRequest_FieldViolation{
		{
			Field:       "name",
			Description: "name should start with 'Report'",
		},
	}
	err := CheckViolations(violations)

	// Check if a gRPC status error is returned
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Expected gRPC status error, got %v", err)
	}

	if st.Code() != codes.InvalidArgument {
		t.Errorf("Expected code %v, got %v", codes.InvalidArgument, st.Code())
	}

	// Check if the details were correctly added
	details := st.Details()
	if len(details) != 1 {
		t.Fatalf("Expected 1 detail, got %d", len(details))
	}

	badRequest, ok := details[0].(*errdetails.BadRequest)
	if !ok {
		t.Fatalf("Expected BadRequest detail, got %T", details[0])
	}

	if len(badRequest.FieldViolations) != 1 {
		t.Fatalf("Expected 1 field violation, got %d", len(badRequest.FieldViolations))
	}

	if badRequest.FieldViolations[0].Field != "name" {
		t.Errorf("Expected Field 'name', got %s", badRequest.FieldViolations[0].Field)
	}

	if badRequest.FieldViolations[0].Description != "name should start with 'Report'" {
		t.Errorf("Expected Description 'name should start with 'Report'', got %s", badRequest.FieldViolations[0].Description)
	}

	// Test case with no violations
	violations = []*errdetails.BadRequest_FieldViolation{}
	err = CheckViolations(violations)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
