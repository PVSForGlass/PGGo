package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gosexy/redis"
)

var (
	endpoint = flag.String("e", "unix://var/run/docker.sock", "Dockerd endpoint")
	redisip = flag.String("r", os.Getenv("REDISIP"), "Redis ip to connect to")
	redisport = flag.Int("s", 6379, "Redis port to connect to")
	addr = flag.String("p", os.Getenv("PORT"), "Address and port to serve dockerui")
	assets = flag.String("a", ".", "Path to the assets")
)

var con *redis.Client = redis.New()

func CreateContext(w http.ResponseWriter, r *http.Request) {
	contextbytes := make([]byte, 9)
	if _, err := rand.Read(contextbytes); err != nil {
		log.Panic(err)
	}
	contextname := base64.URLEncoding.EncodeToString(contextbytes)
	con.SAdd("VMs", contextname)
	w.Write([]byte(contextname))
}

func main() {
	flag.Parse()
	
	if err := con.Connect(*redisip, uint(*redisport)); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/createcontext", CreateContext)
	
	if err := http.ListenAndServe(":" + *addr, nil); err != nil {
		log.Fatal(err)
	}
}
