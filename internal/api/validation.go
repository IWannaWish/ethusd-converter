package api

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func CheckViolations(violations []*errdetails.BadRequest_FieldViolation) error {
	if len(violations) > 0 {
		st := status.Newf(codes.InvalidArgument, "data validation failed")
		st, _ = st.WithDetails(&errdetails.BadRequest{FieldViolations: violations})
		return st.Err()
	}
	return nil
}

func ValidateExample(name string, violations []*errdetails.BadRequest_FieldViolation) []*errdetails.BadRequest_FieldViolation {
	if !strings.HasPrefix(name, "Report") {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: "name should start with 'Report'",
		})
	}
	return violations
}
