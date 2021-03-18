package Index

import "yiarce/dragonnews/reply"

func Index(r *reply.Reply) {
	r.Rs(200, "Hello World!")
}
