package main

import (
	"dragonNews/router"
	"yiarce/core/dhttp"
	"yiarce/core/frame"
)

func main() {
	frame.SetPackageName(`DragonNews`)
	router.Register()
	dhttp.Server(`127.0.0.1`, 55520).Listen()
}
