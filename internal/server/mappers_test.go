// server_test.go
package server

import (
	"reflect"
	"testing"

	apipb "github.com/TimRutte/api/proto/api/gen"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func TestToGetReportResponse(t *testing.T) {
	// Test data
	name := "TestReport"
	date := &timestamp.Timestamp{Seconds: 1622519110, Nanos: 0}
	files := []*apipb.File{
		{Id: 1, Size: "1024", Url: "https://example.com/file1", Hash: "abc123"},
		{Id: 2, Size: "2048", Url: "https://example.com/file2", Hash: "def456"},
	}

	// Call function under test
	resp := ToGetReportResponse(name, date, files)

	// Check name
	if resp.Name != name {
		t.Errorf("Expected Name %s, got %s", name, resp.Name)
	}

	// Check date
	if !reflect.DeepEqual(resp.Date, date) {
		t.Errorf("Expected Date %v, got %v", date, resp.Date)
	}

	// Check files length
	if len(resp.Files) != len(files) {
		t.Fatalf("Expected %d files, got %d", len(files), len(resp.Files))
	}

	// Check individual files
	for i, file := range resp.Files {
		expectedFile := files[i]
		if file.Id != expectedFile.Id || file.Size != expectedFile.Size || file.Url != expectedFile.Url || file.Hash != expectedFile.Hash {
			t.Errorf("Expected File %+v, got %+v", expectedFile, file)
		}
	}
}
