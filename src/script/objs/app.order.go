package objs

import (
	"bookcourt/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 请求通道信息
type RcMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
	Data string `json:"data"`
}

// api返回结果格式
type OrderListRes struct {
	Code int
	Data OrderListData
	Msg  string
}
type OrderListData struct {
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

// 创建订单流程
func (a *App) OrderCreate(lock bool, pctx context.Context) (string, error) {
	ctx, cancel := context.WithCancel(pctx)
	a.LogAdd("开始创建订单...")

	// 创建线程通道
	rc := make(chan RcMsg)
	go a.OrderProc(ctx, rc)
	totalSent := 0
	total404 := 0
	totalTimeOut := 0
	totalDone := 0
	total502 := 0
	msg := RcMsg{}

	a.LogAdd("开启进程监控...")
	hasResult := false

	// 计数器，每10个请求同步一次msg给用户端
	n10 := 0
	for {
		msg = <-rc
		if msg.Code == 888 {
			totalSent++
		}
		if msg.Code == 404 {
			total404++
			totalDone++
		}
		if msg.Code == 999 {
			totalTimeOut++
			totalDone++
		}
		if msg.Code == 502 {
			total502++
			totalDone++
		}
		if msg.Code == 789 {
			a.Lock()
			cancel()
			return "", errors.New("用户中断")
		}
		if msg.Code == 0 {
			totalDone++
			if !lock {
				hasResult = true
			}
		}
		if msg.Code == 1 {
			totalDone++
			hasResult = true
		}
		n10++
		log.Println(msg.Msg)
		// fmt.Printf("\rReq: %5d  Res: %5d  线程: %3d  超时: %8d  404: %8d  502: %8d  Msg:%s", totalSent, totalDone, runtime.NumGoroutine(), totalTimeOut, total404, total502, msg.Msg)
		if n10 == 10 {
			a.LogAdd("Req:" + strconv.Itoa(totalSent) + " Res:" + strconv.Itoa(totalDone) + " Rt:" + strconv.Itoa(runtime.NumGoroutine()) + " 408:" + strconv.Itoa(totalTimeOut) + " 404:" + strconv.Itoa(total404) + "  502:" + strconv.Itoa(total502))
			// 重置计时器
			n10 = 0
		}
		if hasResult {
			break
		}
	}
	a.LogAdd("获得结果, 进行判断...")
	a.St.Lock = true
	cancel()
	a.LogAdd("清理多余线程...")

	if msg.Code == 1 {
		a.LogAdd("成功生成订单! ")
		a.St.Oct = msg.Time
		return msg.Data, nil
	}

	a.LogErr("位置状态: 锁定")
	a.LogAdd("服务器信息:")
	a.LogAdd(msg.Msg)
	return msg.Data, errors.New("需要进行二次判断")
}

// 循环提交订单进程
func (a *App) OrderProc(pctx context.Context, rc chan RcMsg) {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	// 把body放在这里生成
	a.Body["scene"] = "[{\"day\":\"" + utils.GetTomorrowDate() + "\",\"fields\":{\"" + a.St.Venue_id + "\":[" + a.St.Period + "]}}]"

	var msg RcMsg
	msg.Code = 888
	a.St.Lock = false
	for {
		select {
		case <-a.Bc:
			msg.Code = 789
			rc <- msg
			a.St.Lock = true
			return
		default:
			go a.OrderSubmit(ctx, rc)
			rc <- msg
			time.Sleep(time.Duration(a.St.Freq) * time.Millisecond)
		}
	}
}

// 订单提交
func (a *App) OrderSubmit(ctx context.Context, rc chan RcMsg) {
	select {
	case <-ctx.Done():
		return
	default:
		var (
			client  *http.Client
			request *http.Request
			msg     RcMsg
		)
		body := url.Values{}
		for k, v := range a.Body {
			body.Add(k, v)
		}
		fmt.Println(body)
		// 创建请求头
		request, _ = http.NewRequest("POST", a.Api, strings.NewReader(body.Encode()))

		a.OrderHeader(request)

		client = &http.Client{
			Timeout: time.Second * 20,
		}
		response, err := client.Do(request)
		if err != nil {
			msg.Code = 999
			msg.Msg = err.Error()
			rc <- msg
			return
		}
		if response.StatusCode != 200 && response.StatusCode != 201 {
			msg.Code = response.StatusCode
			msg.Msg = response.Status
			rc <- msg
			return
		}
		resBody, _ := ioutil.ReadAll(response.Body)

		json.Unmarshal([]byte(resBody), &msg) // json转map
		defer response.Body.Close()
		rc <- msg
		return
	}
}

func (a *App) OrderList() (OrderListRes, error) {
	var (
		client  *http.Client
		res     OrderListRes
		request *http.Request // 请求对象
	)
	body := url.Values{}
	for k, v := range a.Body {
		body.Add(k, v)
	}
	// 创建请求头
	request, _ = http.NewRequest("GET", a.OApi, strings.NewReader(body.Encode()))
	// 添加头内容
	a.OrderHeader(request)

	client = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		return res, errors.New(response.Status)
	}

	resBody, _ := ioutil.ReadAll(response.Body)

	json.Unmarshal([]byte(resBody), &res) // json转map

	defer response.Body.Close()
	return res, nil
}

// 订单查找
// 	return
// 		订单号	{string}
//		错误	{error}
func (a *App) OrderFind() (string, error) {
	success := false
	orderid := ""
	orderListRes := OrderListRes{}
	errTime := 0
	for {
		// 获取订单列表
		res, err := a.OrderList()
		if err != nil {
			a.LogErr(err.Error())
			a.LogAdd("服务器发生错误, 重试...请等待")
			errTime++
			// 超过6次错误就先休息下。
			if errTime >= 5 {
				a.LogErr("错误太多, 休息10s后重新开始验证...")
				errTime = 0
				time.Sleep(time.Second * 5)
			}
			time.Sleep(time.Millisecond * 700)
			// 获取订单列表失败就重复进行请求
			continue
		}
		orderListRes = res
		break
	}
	// 订单列表中存在未支付订单
	if orderListRes.Data.Total >= 1 {
		// 进行遍历和比对
		for _, r := range orderListRes.Data.List {
			vid, _ := strconv.Atoi(a.St.Venue_id)
			log.Println(vid)
			if r.Config.Scene_list[0].Fields == vid {
				success = true
				// needCreated = true
				orderid = r.Orderid
				a.LogAdd("查询到相应订单! ")
				break
			}
			continue
		}
	}
	// 没有订单
	if orderListRes.Data.Total == 0 {
		return orderid, errors.New("位置被抢, 无生成订单！")
	}
	if orderListRes.Code == 401 {
		return orderid, errors.New("权限错误! 401")
	}
	// 查看success标记是否被标为true
	if !success {
		return orderid, errors.New("位置被抢, 但存在其他未支付订单。")
	}
	return orderid, nil
}

func (a *App) OrderNewForm(orderid string) (RcMsg, error) {
	var (
		client  *http.Client
		request *http.Request // 请求对象
		res     RcMsg
		reqBody map[string]string
	)

	// 生成url编码格式的body，用来生成application/x-www-form-urlencoded格式数据
	reqBody = map[string]string{
		"orderid":         orderid,
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

	request, _ = http.NewRequest("POST", a.Api, strings.NewReader(body.Encode()))
	a.OrderHeader(request)
	client = &http.Client{
		Timeout: time.Second * 5,
	}
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		return res, errors.New(response.Status)
	}
	resBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(resBody), &res)
	defer response.Body.Close()
	return res, nil
}

func (a *App) OrderHeader(r *http.Request) {
	r.Header.Add("Origin", "http://gym.dazuiwl.cn")
	r.Header.Add("Host", "gym.dazuiwl.cn")
	r.Header.Add("Referer", "http://gym.dazuiwl.cn/h5/")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Accept-Language", "en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	r.Header.Add("Cookie", a.Cookie)
	r.Header.Add("Pragma", "no-cache")
	r.Header.Add("Proxy-Connection", "keep-alive")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Mobile Safari/537.36 Edg/101.0.1210.32")
	r.Header.Add("X-Requested-With", "XMLHttpRequest")
	r.Header.Add("token", a.Token)
}
