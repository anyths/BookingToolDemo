package modules

import (
	"bookcourt/objs"
	"bookcourt/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Res struct {
	Code int
	Data ResData
	Msg  string
}
type ResData struct {
	Limit string
	Total int
	List  []OrderDetail
}
type OrderDetail struct {
	Amount           int64
	Config           DetailConfig
	Createtime       int
	Id               int
	Ip               string
	Memo             string
	Orderid          string
	Ordertype        string
	Pay_id           int
	Paytype          string
	Sport_events_id  int
	Sportevents_name string
	Status           string
	Updatetime       int64
	User_id          int64
	Useragent        string
	Venue_id         int
	Venue_name       string
}
type DetailConfig struct {
	Pay_id     int
	Scene_list []SceneItem
}
type SceneItem struct {
	Fields int
}

func FindOrder(l *utils.Logs, c objs.Conf) (string, error) {
	n := 0
	success := false
	orderid := ""
	for {
		res, err := Find(c)
		if err != nil {
			l.Err(err.Error())
			l.Add("服务器超时, 重试...请等待")
			continue
		}
		if res.Data.Total >= 1 {
			for _, r := range res.Data.List {
				vid, _ := strconv.Atoi(c.Venue_id)
				if r.Config.Scene_list[0].Fields == vid {
					success = true
					// needCreated = true
					orderid = r.Orderid
					l.Add("查询到相应订单! ")
					break
				}
				continue
			}
			break
		}
		if res.Data.Total == 0 {
			l.Add("订单查询: 无成功订单")
			l.Err("位置已被抢")
			break
		}
		if res.Code == 401 {
			l.Err("权限错误, 验证失效, 请重新配置token和cookie。")
			break
		}
		n++
		fmt.Printf(">>> 验证次数：%d    @响应码:%4d | @订单：%4d \n", n, res.Code, res.Data.Total)
		time.Sleep(time.Millisecond * 700)
	}
	if !success {
		return orderid, errors.New("未查询到有效订单")
	}
	return orderid, nil
}

func Find(c objs.Conf) (Res, error) {
	var (
		client  *http.Client
		res     Res
		request *http.Request // 请求对象
	)
	body := url.Values{}
	for k, v := range c.Body {
		body.Add(k, v)
	}
	// 创建请求头
	request, _ = http.NewRequest("GET", c.OApi, strings.NewReader(body.Encode()))
	// 添加头内容
	Header(request, c)

	client = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(request)
	if err != nil {
		return res, errors.New(err.Error())
	}

	resBody, _ := ioutil.ReadAll(response.Body)

	json.Unmarshal([]byte(resBody), &res) // json转map

	defer response.Body.Close()
	return res, nil
}

func NewOrder(oid string, c objs.Conf) (RcMsg, error) {
	var (
		client  *http.Client
		request *http.Request // 请求对象
		res     RcMsg
		reqBody map[string]string
	)

	// 生成url编码格式的body，用来生成application/x-www-form-urlencoded格式数据
	reqBody = map[string]string{
		"orderid":         oid,
		"card_id":         "",
		"sport_events_id": "",
		"money":           "",
		"ordertype":       "",
		"paytype":         "bitpay",
		"scene":           "",
		"openid":          "",
	}

	body := url.Values{}
	for k, v := range reqBody {
		body.Add(k, v)
	}

	api := "http://gym.dazuiwl.cn/api/order/submit"

	request, _ = http.NewRequest("POST", api, strings.NewReader(body.Encode()))
	Header(request, c)
	client = &http.Client{
		Timeout: time.Second * 20,
	}
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	resBody, _ := ioutil.ReadAll(response.Body)

	json.Unmarshal([]byte(resBody), &res)
	defer response.Body.Close()
	return res, nil
}
