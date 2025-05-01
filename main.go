package main

import (
	"fmt"

	"github.com/kaibling/cerodev/bootstrap/app"
	"github.com/kaibling/cerodev/config"
)

var (
	version      string //nolint:gochecknoglobals,nolintlint
	buildTime    string //nolint:gochecknoglobals
	architecture string //nolint:gochecknoglobals
)

func main() {
	config.Version = version
	config.BuildTime = buildTime
	config.Architecture = architecture

	if err := app.New(); err != nil {
		fmt.Println(err.Error()) //nolint: forbidigo
	}
}
