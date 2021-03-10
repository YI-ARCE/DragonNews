package Index

import "yiarce/dragonnews/reply"

func Index(r *reply.Reply) {
	r.Return(200, "Hello World!")
	return
}
