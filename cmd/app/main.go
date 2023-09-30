package main

import (
	"log"

	"github.com/evermake/git-diff-view/internal/app"
)

func main() {
	a := app.New()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
