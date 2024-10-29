package api

import (
	"context"
	apipb "github.com/TimRutte/api/proto/api/gen"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"strconv"
)

type Api struct {
	apipb.UnimplementedApiServer
}

func (api Api) GetReport(ctx context.Context, request *apipb.GetReportRequest) (*apipb.GetReportResponse, error) {
	date := request.GetDate()
	name := request.GetName()

	// Validate input if needed
	violations := make([]*errdetails.BadRequest_FieldViolation, 0, 10)
	violations = ValidateExample(name, violations)

	err := CheckViolations(violations)
	if err != nil {
		return nil, err
	}

	// To some fancy stuff
	log.Info().Msg("We can do some fancy stuff here...")

	// Simulate fetching files
	files := []*apipb.File{
		{
			Url:  "path/to/the/file1",
			Size: strconv.Itoa(100),
			Hash: "hash1",
		},
		{
			Url:  "path/to/the/file2",
			Size: strconv.Itoa(200),
			Hash: "hash2",
		},
		{
			Url:  "path/to/the/file3",
			Size: strconv.Itoa(300),
			Hash: "hash3",
		},
	}

	// Return the response in right format
	return ToGetReportResponse(name, date, files), nil
}
