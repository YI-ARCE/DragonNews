package dhttp

type HttpRequest interface {
	Session() SessionReader
	Host() string
}
