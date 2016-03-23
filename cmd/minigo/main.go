package main

import (
	"flag"
	"log"
	"net/http"

	"nelhage.com/minigo/web"
)

func main() {
	root := flag.String("root", "public", "Path to the http public file root")
	bind := flag.String("bind", "127.0.0.1:4040", "listen address")
	flag.Parse()
	srv := &web.Server{
		Public: *root,
	}
	srv.Bind(http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
