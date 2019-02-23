package main

import (
	"fmt"
	gocache "github.com/pmylund/go-cache"
	"github.com/valyala/fasthttp"
	"log"
)

func logPanics(function fasthttp.RequestHandler)fasthttp.RequestHandler{
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if x := recover(); x!=nil{
				log.Printf("[%v] caught panic: %v",ctx.RemoteAddr(),x)
			}
		}()
		function(ctx)
	}
}


func middlewareEndpoint(ctx *fasthttp.RequestCtx,f func()string){
	key := string(ctx.Path()) + ctx.QueryArgs().String()
	_ = cache.Add(key, f(), gocache.DefaultExpiration)
	response,_ := cache.Get(key)
	_,_ = fmt.Fprint(ctx, fmt.Sprintf("%v",response))
}



