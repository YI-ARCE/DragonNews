package Index

import (
	Db "yiarce/dragonnews/orm"
	"yiarce/dragonnews/reply"
)

func Index(r *reply.Reply) {

	//s1 := md5.Sum([]byte(`asodhakjsdhkahd`))
	//ss := string(hex.EncodeToString([]byte{0:s1[5],1:s1[6],2:s1[7],3:s1[8],4:s1[9],5:s1[10],6:s1[11],7:s1[12]}))
	//fmt.Println(ss)
	Db.Table(`123`).Where(`bb = cc`).Field(`(ss1)`).Find()
	r.Rs(200, "Hello World!")
}
