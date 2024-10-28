package server

import (
	apipb "github.com/TimRutte/api/proto/api/gen"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func ToGetReportResponse(name string, date *timestamp.Timestamp, files []*apipb.File) *apipb.GetReportResponse {
	var filesResponse []*apipb.File
	for _, c := range files {
		filesResponse = append(filesResponse, &apipb.File{
			Id:   c.Id,
			Size: c.Size,
			Url:  c.Url,
			Hash: c.Hash,
		})
	}
	return &apipb.GetReportResponse{
		Name:  name,
		Date:  date,
		Files: filesResponse,
	}
}
