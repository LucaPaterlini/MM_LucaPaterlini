package main

import (
	"../coreCli"
	"github.com/valyala/fasthttp"

)

// routing
func routingHandler (ctx *fasthttp.RequestCtx){
	ctx.SetContentType("text/json; charset=utf-8")

	switch string(ctx.Path()) {

	//Call n°0
	case "/hello":
		middlewareEndpoint(ctx,coreCli.Create100NewIds)

	default:
		ctx.Error(ERRPATH,fasthttp.StatusNotFound)

	}
}



