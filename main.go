package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	endpoint = flag.String("e", "unix://var/run/docker.sock", "Dockerd endpoint")
	addr = flag.String("p", ":"+os.Getenv("PORT"), "Address and port to serve dockerui")
	assets = flag.String("a", ".", "Path to the assets")
)

func main() {
	flag.Parse()
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
