package objs

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// 设置 运行中 状态
func (a *App) Doing() {
	a.St.Status = false

	a.LogAdd("开始执行[" + a.St.Cmd + "]")
	a.Sc <- 0
}
func (a *App) Done() {
	a.St.Status = true

	a.LogAdd("执行[" + a.St.Cmd + "]成功!")
	a.St.Cmd = ""
	a.Sc <- 0
}
func (a *App) Err(err string) {
	a.LogErr(err)
	a.St.Status = true
	a.St.Cmd = ""
	a.Sc <- 0
}

type Message struct {
	Ok      bool   `json:"ok"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

func (a *App) ProcWs() {
	for {
		WsDone := make(chan struct{})
		go a.Ws(WsDone)
		<-WsDone
		log.Println("ws断开连接, 15s后尝试重连...")
		time.Sleep(time.Second * 15)
	}
}

var addr = flag.String("addr", "wx.gt0.cn", "http service address")

// var addr = flag.String("addr", "localhost:888", "http service address")

func (a *App) Ws(WsDone chan struct{}) {
	defer close(WsDone)

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: a.WsPre, Host: *addr, Path: "/api/ws/" + a.St.Key + "/" + a.St.Id}
	log.Println("已接入服务...[key]:" + a.St.Key + " [id]:" + a.St.Id)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer c.Close()

	// 创建通道done，用来识别ws是否关闭
	done := make(chan struct{})

	// 创建协程读取ws连接里的信息
	go func() {
		defer close(done)
		for {
			// for循环读取通道消息，如果有消息就能读取
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			a.Rc <- message
		}
	}()

	// 创建一个计时器用来发送心跳
	ticker := time.NewTicker(time.Second * 8)
	heartCheck := time.NewTicker(time.Second * 12)
	defer ticker.Stop()
	defer heartCheck.Stop()
	heart := false
	for {
		select {
		case <-done:
			// done通道有信息，ws关闭，所以这个进程也要关闭
			return
		case <-a.Sc:
			// 接收到更新状态
			var d WsData
			d.Sta = a.St
			d.Type = "sta"
			bytesData, err := json.Marshal(d)
			if err != nil {
				a.LogErr(err.Error())
				continue
			}
			// 发送消息
			err = c.WriteMessage(websocket.TextMessage, bytesData)
			if err != nil {
				a.LogErr(err.Error())
				return
			}
		case <-ticker.C:
			err = c.WriteMessage(websocket.TextMessage, []byte(""))
			if err != nil {
				a.LogErr(err.Error())
				return
			}
		case <-heartCheck.C:
			if !heart {
				log.Println("心跳错误")
				return
			}
			heart = false
		case byteCmd := <-a.Rc:
			if len(byteCmd) == 0 {
				// 这里是心跳
				heart = true
				continue
			}
			var msg WsRec
			json.Unmarshal(byteCmd, &msg)
			go a.ParseCmd(msg)
		case m := <-a.Lc:
			byteM, _ := json.Marshal(m)
			err = c.WriteMessage(websocket.TextMessage, byteM)
			if err != nil {
				a.LogErr(err.Error())
				return
			}
		case <-interrupt:
			// 接收到系统信号
			log.Println("interrupt")

			// 通过发送一个关闭消息，清楚地关闭连接。同时等待系统关闭连接
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}

// func (a *App) SyncStatusReq() {
// 	api := a.Tapi + "/sta/" + a.St.Key + "/" + a.St.Id

// 	msg := (*a).St

// 	bytesData, err := json.Marshal(msg)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	request, _ := http.NewRequest("POST", api, bytes.NewReader(bytesData))
// 	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
// 	client := http.Client{
// 		Timeout: time.Second * 8,
// 	}
// 	_, err = client.Do(request)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// }

// func (a *App) Remove() {
// 	api := a.Tapi + "/sta/" + a.St.Key + "/" + a.St.Id
// 	request, _ := http.NewRequest("DELETE", api, nil)
// 	client := &http.Client{
// 		Timeout: time.Second * 5,
// 	}
// 	_, err := client.Do(request)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// }

type ConfRes struct {
	StatusCode int  `json:"statusCode"`
	Data       Conf `json:"data"`
}

func (a *App) GetConfig(conf Conf) error {
	// api := a.Tapi + "/conf/" + a.St.Key + "/" + a.St.Id

	// request, _ := http.NewRequest("GET", api, nil)
	// client := &http.Client{}

	// var res ConfRes
	// response, err := client.Do(request)
	// if err != nil {
	// 	return err
	// }
	// if response.StatusCode != 200 {
	// 	return errors.New("空数据 或 服务器问题")
	// }
	// resbody, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	return err
	// }
	// json.Unmarshal([]byte(resbody), &res)

	a.St.Sport_events_id = conf.Sport_events_id
	a.St.Venue_id = conf.Venue_id
	a.St.Period = conf.Period
	a.St.Money = conf.Money
	a.St.Freq = conf.Freq

	return nil
}
