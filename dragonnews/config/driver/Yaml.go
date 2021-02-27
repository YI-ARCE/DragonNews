package driver

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	Sql        Sql        `yaml:"sql"`
	DragonNews DragonNews `yaml:"dragonnews"`
}

type Sql struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DragonNews struct {
	Request  Request  `yaml:"request"`
	Response Response `yaml:"response"`
	Server   Server   `yaml:"server"`
}

type Request struct {
	Enable     string `yaml:"enable"`
	SecretMode string `yaml:"secretMode"`
	Public     string `yaml:"public"`
	Private    string `yaml:"private"`
}

type Response struct {
	ErrorNotice bool   `yaml:"errorNotice"`
	ReturnType  string `yaml:"returnType"`
	StatusCode  int    `yaml:"statusCode"`
	ErrorData   string `yaml:"errorData"`
}

type Server struct {
	Host string
	Port string
}

type File struct {
	File []byte
}

/*
	获取参数配置文件
*/
func GetConfig() *config {
	var Path, _ = os.Getwd()
	c := config{}
	conf, err := ioutil.ReadFile(Path + "/config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(conf, &c)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	return &c
}

//读取文件配置转化成map
func (f File) Get() (map[string]string, error) {
	c := map[string]string{}
	err := yaml.Unmarshal(f.File, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
