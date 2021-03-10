package http

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	Yaml "yiarce/dragonnews/config/driver"
	"yiarce/dragonnews/reply"
	"yiarce/dragonnews/route"
)

var config Yaml.Response

func Start(server Yaml.DragonNews) {
	config = server.Response
	http.HandleFunc("/", request)
	fmt.Print("----DragonNews----\n")
	fmt.Print("----Start   OK----\n")
	fmt.Print("----Http    OK----\n")
	fmt.Print("----Listening-----\n")
	fmt.Print("----Port " + server.Server.Port + "-----\n")
	_ = http.ListenAndServe(server.Server.Host+":"+server.Server.Port, nil)
}

func request(w http.ResponseWriter, r *http.Request) {
	var Request strings.Builder
	Uri := r.RequestURI
	UriIndex := strings.Index(Uri, "?")
	if UriIndex != -1 {
		Uri = Uri[0:UriIndex]
	}
	if r.Method == "OPTIONS" {
		_, _ = w.Write([]byte{})
	}
	var returnType string
	switch config.ReturnType {
	case "json":
		returnType = "application/json"
		break
	case "html":
		returnType = "text/html"
		break
	case "xml":
		returnType = "application/xml"
		break
	default:
		returnType = "text/plain"
		break
	}
	Request.WriteString(Uri)
	Request.WriteString("_")
	Request.WriteString(r.Method)
	w.Header().Set("Content-type", returnType)
	Reply := reply.Start(w, r)
	defer func() {
		if r := recover(); r != nil {
			//*\t"+r+"\n
			fmt.Println(r)
			msg := "Positioning error:\n"
			for i := 1; i <= 5; i++ {
				_, file, line, _ := runtime.Caller(i)
				msg += " *\t" + file + "(Line:" + strconv.Itoa(line) + ")"
				if i != 5 {
					msg += "\n"
				}
			}
			Reply.Log.Insert(msg)
			if config.ErrorNotice {
				str, flag := r.(string)
				if flag {
					msg = strings.ReplaceAll(config.ErrorData, "[errorMsg]", str)
				}
			}
			Reply.Return(config.StatusCode, msg)
		}
		if Reply.Log.Judge() {
			err := Reply.Log.Out()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	if function, ok := route.Route[Request.String()]; ok {
		function(&Reply)
	} else if function, ok = route.Route[Uri]; ok {
		function(&Reply)
	} else {
		w.WriteHeader(config.StatusCode)
		_, _ = w.Write([]byte(`Nonexistent address`))
	}
}
