package main

import (
	"log"

	"github.com/evermake/git-diff-view/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
