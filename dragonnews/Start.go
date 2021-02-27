package dragonnews

import (
	. "yiarce/application"
	Config "yiarce/dragonnews/config"
)

func Start() {
	Route()
	Config.Init()
}
