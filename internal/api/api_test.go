package api

import (
	"context"
	"testing"
	"time"

	apipb "github.com/TimRutte/api/proto/api/gen"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestApi_GetReport_ValidRequest(t *testing.T) {
	// Initialize the API server
	apiServer := Api{}

	// Create a valid request
	date := timestamppb.New(time.Now())
	request := &apipb.GetReportRequest{
		Name: "Report 2024",
		Date: date,
	}

	// Call the GetReport method
	response, err := apiServer.GetReport(context.Background(), request)

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response
	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Name != request.Name {
		t.Errorf("Expected name %s, got %s", request.Name, response.Name)
	}

	if len(response.Files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(response.Files))
	}
}

func TestApi_GetReport_InvalidRequest(t *testing.T) {
	// Initialize the API server
	apiServer := Api{}

	// Create an invalid request (name doesn't start with "Report")
	request := &apipb.GetReportRequest{
		Name: "InvalidName",
		Date: timestamppb.New(time.Now()),
	}

	// Call the GetReport method
	response, err := apiServer.GetReport(context.Background(), request)

	// Check that an error is returned
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	// Check that the response is nil
	if response != nil {
		t.Fatal("Expected response to be nil")
	}

	// Check that the error is of the expected type and has the correct details
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.InvalidArgument {
		t.Fatalf("Expected InvalidArgument error, got %v", err)
	}

	// Verify the details of the error
	if len(st.Proto().GetDetails()) == 0 {
		t.Fatal("Expected error details, got none")
	}

	// Get the first detail and assert its type
	details := st.Proto().GetDetails()[0]

	// Assert that the details are of type BadRequest
	var badRequest errdetails.BadRequest
	if err := proto.Unmarshal(details.GetValue(), &badRequest); err != nil {
		t.Fatalf("Failed to unmarshal BadRequest details: %v", err)
	}

	// Verify the field violations
	if len(badRequest.FieldViolations) == 0 {
		t.Fatal("Expected field violations, got none")
	}

	if badRequest.FieldViolations[0].Field != "name" {
		t.Errorf("Expected field violation for 'name', got %s", badRequest.FieldViolations[0].Field)
	}
}
