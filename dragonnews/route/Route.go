package route

import "yiarce/dragonnews/reply"

var Route = map[string]func(reply *reply.Reply){}

func Get(url string, method func(reply *reply.Reply)) {
	Route[url+"_GET"] = method
}

func Post(url string, method func(reply *reply.Reply)) {
	Route[url+"_POST"] = method
}

func Rule(url string, method func(reply *reply.Reply)) {
	Route[url] = method
}
