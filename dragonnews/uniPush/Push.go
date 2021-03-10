package uniPush

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"yiarce/dragonnews/curl"
	"yiarce/dragonnews/general"
)

type Uni struct {
	AppId        string
	AppKey       string
	AppSecret    string
	Token        string
	EndTimestamp string
}

//详情参考 https://docs.getui.com/getui/server/rest_v2/push/
type Single struct {
	//任务ID
	RequestId string `json:"request_id"`
	//ttl设置时间
	Settings struct {
		Ttl int `json:"ttl"`
	} `json:"settings"`
	//用户ID
	Audience struct {
		//内容写Cid
		Alias [1]string `json:"alias"`
	} `json:"audience"`
	//推送
	PushMessage struct {
		//内容
		Notification struct {
			//标题
			Title string `json:"title"`
			//内容
			Body string `json:"body"`
			//点击类型
			ClickType string `json:"click_type"`
			//地址
			Url string `json:"url"`
		} `json:"notification"`
	} `json:"push_message"`
}

type singleReturn struct {
	Code   int
	Msg    string
	TaskId string
	Cid    string
	Status string
}

func Init(appId string, appKey string, appSecret string) (*Uni, error) {
	h := sha256.New()
	times := general.Date().Timestamp("ms")
	h.Write([]byte(appKey + times + appSecret))
	sign := hex.EncodeToString(h.Sum(nil))
	bodys := map[string]string{
		"sign":      sign,
		"timestamp": times,
		"appkey":    appKey,
	}
	body, _ := json.Marshal(bodys)
	replys, err := curl.Post("https://restapi.getui.com/v2/"+appId+"/auth", curl.Json, map[string]string{}, string(body))
	if err != nil {
		return nil, err
	}
	if replys.Data["code"].(float64) != 0 {
		return nil, errors.New(replys.Data["msg"].(string))
	}
	return &Uni{
		AppId:        appId,
		AppKey:       appKey,
		AppSecret:    appSecret,
		Token:        (replys.Data["data"].(map[string]interface{}))["token"].(string),
		EndTimestamp: (replys.Data["data"].(map[string]interface{}))["expire_time"].(string),
	}, nil
}

//对某个别名推送内容
func (u *Uni) SingleAliasPush(s Single) (singleReturn, error) {
	body, _ := json.Marshal(s)
	replys, err := curl.Post("https://restapi.getui.com/v2/"+u.AppId+"/push/single/alias", curl.Json, map[string]string{"Token": u.Token}, string(body))
	if err != nil {
		return singleReturn{}, err
	}
	if replys.Data["code"].(float64) != 0 {
		return singleReturn{}, errors.New(replys.Data["msg"].(string))
	}
	fmt.Println(replys.Data)
	task := reflect.ValueOf(replys.Data["data"]).MapRange()
	task.Next()
	taskId := task.Key().String()
	c := task.Value().Elem().MapRange()
	c.Next()
	return singleReturn{int(replys.Data["code"].(float64)), replys.Data["msg"].(string), taskId, c.Key().String(), c.Value().Elem().String()}, nil
}

//对部分用户推送同一条内容
func (u *Uni) PartPush() {

}

//对部分用户,每个用户都为不同的数据
func (u *Uni) PartDiffPush() {

}

//对全体用户推送内容
func (u *Uni) AllPush() {

}
