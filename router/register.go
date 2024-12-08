package router

import (
	"dragonNews/app/index"
	"yiarce/core/dhttp"
)

func Register() {
	dhttp.Get(`index/hello`, index.Hello)
}
