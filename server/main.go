package main

import (
	"flag"

	"github.com/EmeraldLS/kv-db/server/repository"
)

func main() {
	port := flag.String("p", "3031", "port server is binding to")
	flag.Parse()

	repository.Listen(*port)
}
