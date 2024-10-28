package main

import (
	"fmt"

	"github.com/TimRutte/api/internal/version"
)

func Version() {
	fmt.Printf("Build date: %s, Version: %s\n", version.BuildDate, version.Version)
}
