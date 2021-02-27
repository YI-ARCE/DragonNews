package application

import (
	Index "yiarce/application/packages/index"
	"yiarce/dragonnews/route"
)

func Route() {
	route.Get("/StartOrderCancelRollBack", Index.StartFunc)
	route.Get("/StopOrderCancelRollBack", Index.StopFunc)
	route.Get("/Dec", Index.WxUniTest)
	//route.Get("/LogTest",Index.LogTest)
}
