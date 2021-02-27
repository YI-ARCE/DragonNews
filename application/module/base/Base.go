package base

type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data ...interface{}) Reply {
	reply := Reply{}
	reply.Code = 20000
	reply.Msg = "success"
	if len(data) > 0 {
		reply.Data = data[0]
		return reply
	}
	reply.Data = nil
	return reply
}

func Error(msg string, data ...interface{}) Reply {
	reply := Reply{}
	reply.Msg = msg
	reply.Code = 50000
	if len(data) > 0 {
		reply.Data = data[0]
		return reply
	} else {
		return reply
	}
}
