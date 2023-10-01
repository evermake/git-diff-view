package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/evermake/git-diff-view/internal/app"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	repoPath := flag.String("repo", wd, "path to the git repository")
	port := flag.String("port", "7777", "port to listen on")
	flag.Parse()

	a, err := app.New(
		app.WithRepoPath(*repoPath),
		app.WithAddr(net.JoinHostPort("", *port)),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
