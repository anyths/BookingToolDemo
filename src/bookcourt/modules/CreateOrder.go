package modules

import (
	"bookcourt/objs"
	"bookcourt/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

type RcMsg struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
	Data  string `json:"data"`
	Total int
}

func CreateOrder(l *utils.Logs, c *objs.Conf, lock bool) (string, error) {
	l.Add("开始创建订单...")
	ctx, cancel := context.WithCancel(context.Background())
	// 创建线程通道
	rc := make(chan RcMsg)
	go Proc(ctx, *c, rc)
	totalSent := 0
	total404 := 0
	totalTimeOut := 0
	totalDone := 0

	msg := RcMsg{}

	l.Add("开启进程监控...")
	hasResult := false
	for {
		msg = <-rc
		if msg.Total == 1 {
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
		// if msg.Code == 888 {
		// 	totalDone++
		// }
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
		fmt.Printf("\r>> 总请求: %8d  已处理: %8d  线程数: %3d  超时: %8d  404: %8d", totalSent, totalDone, runtime.NumGoroutine(), totalTimeOut, total404)
		if hasResult {
			fmt.Println("")
			l.Add("获得结果, 进行判断...")
			break
		}
	}
	l.Add("清理多余线程...")
	cancel()
	if msg.Code == 1 {
		l.Add("成功生成订单! ")
		return msg.Data, nil
	}

	l.Add("位置状态: 锁定")
	return msg.Data, errors.New("需要进行二次判断")
}

func Proc(pctx context.Context, c objs.Conf, rc chan RcMsg) {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	var msg RcMsg
	msg.Total = 1
	msg.Code = 888
	for {
		go Req(ctx, c, rc)
		rc <- msg
		time.Sleep(time.Duration(c.Freq) * time.Millisecond)
	}
}

func Req(ctx context.Context, c objs.Conf, rc chan RcMsg) {
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
		for k, v := range c.Body {
			body.Add(k, v)
		}

		// 创建请求头
		request, _ = http.NewRequest("POST", c.Api, strings.NewReader(body.Encode()))

		Header(request, c)

		client = &http.Client{
			Timeout: time.Second * 16,
		}
		response, err := client.Do(request)
		if err != nil {
			msg.Code = 999
			msg.Total = 0
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

func Header(r *http.Request, c objs.Conf) {
	r.Header.Add("Origin", "http://gym.dazuiwl.cn")
	r.Header.Add("Host", "gym.dazuiwl.cn")
	r.Header.Add("Referer", "http://gym.dazuiwl.cn/h5/")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Accept-Language", "en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	r.Header.Add("Cookie", c.Cookie)
	r.Header.Add("Pragma", "no-cache")
	r.Header.Add("Proxy-Connection", "keep-alive")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Mobile Safari/537.36 Edg/101.0.1210.32")
	r.Header.Add("X-Requested-With", "XMLHttpRequest")
	r.Header.Add("token", c.Token)
}
