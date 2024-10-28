package main_test

import (
	"io"
	"os"
	"testing"

	versionWrapper "github.com/TimRutte/api/cmd/api"
	"github.com/TimRutte/api/internal/version"
	"github.com/stretchr/testify/assert"
)

func TestVersionPrintsTheCurrentProjectVersion(t *testing.T) {
	testVersion := "V1"
	version.Version = testVersion

	rescueStdout := os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer

	versionWrapper.Version()

	err := writer.Close()
	if err != nil {
		t.Error("Failed to close io writer")
	}
	output, _ := io.ReadAll(reader)
	os.Stdout = rescueStdout

	assert.Equal(t, "Build date: , Version: "+testVersion+"\n", string(output))
}
