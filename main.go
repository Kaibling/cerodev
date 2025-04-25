package main

import (
	"fmt"

	"github.com/kaibling/cerodev/bootstrap/app"
)

func main() {
	if err := app.New(); err != nil {
		fmt.Println(err.Error()) //nolint: forbidigo
	}
}
