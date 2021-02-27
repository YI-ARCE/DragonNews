package config

import (
	"io/ioutil"
	"os"
	Yaml "yiarce/dragonnews/config/driver"
	"yiarce/dragonnews/http"
	Db "yiarce/dragonnews/orm"
)

func Init() {
	conf := Yaml.GetConfig()
	Db.Init(Db.Sql(conf.Sql))
	http.Start(conf.DragonNews)
}

//获取根目录下的config目录内的yaml配置文件
//  只需输入名称即可,无需后缀,注意区分大小写
func Get(name string) (map[string]string, error) {
	p, _ := os.Getwd()
	conf, err := ioutil.ReadFile(p + "/config/" + name + ".yaml")
	if err != nil {
		return nil, err
	}

	File := Yaml.File{
		conf,
	}

	return File.Get()
}
