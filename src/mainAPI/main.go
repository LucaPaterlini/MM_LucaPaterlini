package main

import (
	"flag"
	cors "github.com/AdhityaRamadhanus/fasthttpcors"
	gocache "github.com/pmylund/go-cache"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

var (
	addr = flag.String("addr",IPADDR+":"+strconv.Itoa(PORT),"TCP address to listen to")
	compress = flag.Bool("compress",false, "Whether to enable transparent response compression ")
	cache = gocache.New(1*time.Minute, 3*time.Minute)
)


func main(){

	// Corse Handeler
	withCors := cors.DefaultHandler()

	flag.Parse()
	h := withCors.CorsMiddleware(logPanics(routingHandler))
	if *compress {
		h = fasthttp.CompressHandler(h)
	}
	if err := fasthttp.ListenAndServe(*addr,h);err != nil{
		log.Fatalf("Errror in ListenAndServer: %s",err)
	}

}
