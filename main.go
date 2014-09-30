package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gosexy/redis"
)

var (
	endpoint    = flag.String("e", "unix://var/run/docker.sock", "Dockerd endpoint")
	redisurlstr = flag.String("r", os.Getenv("REDIS_URL"), "Redis url to connect to")
	addr        = flag.String("p", os.Getenv("PORT"), "Address and port to serve dockerui")
	assets      = flag.String("a", ".", "Path to the assets")
)

var (
	VMsKey = "VMs"
)

var con *redis.Client = redis.New()

func CreateContext(w http.ResponseWriter, r *http.Request) {
	contextbytes := make([]byte, 9)
	if _, err := rand.Read(contextbytes); err != nil {
		log.Panic(err)
	}
	contextname := base64.URLEncoding.EncodeToString(contextbytes)
	con.SAdd(VMsKey, contextname)
	w.Write([]byte(contextname))
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	pathSplit := strings.Split(r.URL.Path, "/")
	if len(pathSplit) != 4 {
		http.NotFound(w, r)
		return
	}
	context := pathSplit[2]
	file := pathSplit[3]

	if exists, err := con.SIsMember(VMsKey, context); err != nil {
		log.Panic(err)
	} else if !exists {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	if _, err := con.HSet(context, file, string(body)); err != nil {
		log.Panic(err)
	}

	w.Write([]byte("Success!"))
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
	http.HandleFunc("/uploadfile/", UploadFile)

	if err := http.ListenAndServe(":"+*addr, nil); err != nil {
		log.Fatal(err)
	}
}
