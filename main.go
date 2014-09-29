package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"os"

	"github.com/gosexy/redis"
)

var (
	endpoint = flag.String("e", "unix://var/run/docker.sock", "Dockerd endpoint")
	redisurlstr = flag.String("r", os.Getenv("REDIS_URL"), "Redis url to connect to")
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
	
	redisurl, err := url.Parse(*redisurlstr)
	if err != nil {
		log.Fatal(err)
	}
	redisConInfo := strings.Split(redisurl.Host, ":")
	redisport, err := strconv.Atoi(redisConInfo[1])
	if err != nil {
		log.Fatal(err)
	}
	if err = con.Connect(redisConInfo[0], uint(redisport)); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/createcontext", CreateContext)
	
	if err := http.ListenAndServe(":" + *addr, nil); err != nil {
		log.Fatal(err)
	}
}
