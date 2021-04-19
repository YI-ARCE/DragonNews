package application

import (
	Index "yiarce/application/packages/index"
	"yiarce/dragonnews/route"
)

func Route() {
	route.Post("/index", Index.Index)
}
