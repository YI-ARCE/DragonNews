package application

import (
	Index "yiarce/application/packages/index"
	"yiarce/dragonnews/route"
)

func Route() {
	route.Get("/index", Index.Index)
}
