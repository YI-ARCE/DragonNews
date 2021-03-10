package application

import (
	Index "yiarce/application/packages/index"
	"yiarce/dragonnews/route"
)

func Route() {
	route.Get("/StartOrderCancelRollBack", Index.StartFunc)
	route.Get("/StopOrderCancelRollBack", Index.StopFunc)
	route.Get("/ceshi", Index.WxUniTest)
	route.Get("/key", Index.Key)
	//route.Get("/LogTest",Index.LogTest)
}
