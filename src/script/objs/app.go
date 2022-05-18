package objs

import (
	"bookcourt/utils"
	"log"
	"os"
	"strconv"

	"gopkg.in/ini.v1"
)

type App struct {
	St Status
	// 接收通道
	Rc chan []byte
	// 状态更新通道
	Sc chan int
	// 中断程序通道
	Bc chan int
	// 命令获取通道
	Cc chan WsRec
	// 消息更新通道
	Lc chan WsData
	// 下单api
	Api     string
	OApi    string
	Cookie  string
	Token   string
	Date    string
	Host    string
	HttpPre string
	WsPre   string
	Body    map[string]string

	Form string
	// 测试api
	Tapi string
}

func (a *App) LogAdd(str string) {
	log.Println(str)
	var msg Message
	msg.Ok = true
	msg.Time = utils.GetNowTime()
	msg.Content = str
	var data WsData
	data.Type = "msg"
	data.Msg = msg
	a.Lc <- data
}
func (a *App) LogErr(str string) {
	log.Println(str)
	var msg Message
	msg.Ok = false
	msg.Time = utils.GetNowTime()
	msg.Content = str
	var data WsData
	data.Type = "msg"
	data.Msg = msg
	a.Lc <- data
}

type Status struct {
	Key             string `json:"key"`
	Id              string `json:"id"`
	Oid             string `json:"oid"`
	Status          bool   `json:"status"`
	Cmd             string `json:"cmd"`
	Form            bool   `json:"form"`
	Pay             string `json:"pay"`
	Pct             string `json:"pct"`
	Oct             string `json:"oct"`
	Sport_events_id string `json:"sport_events_id"`
	Money           string `json:"money"`
	Venue_id        string `json:"venue_id"`
	Period          string `json:"period"`
	Freq            int    `json:"freq"`
	Lock            bool   `json:"lock"`
}

type Conf struct {
	Sport_events_id string `json:"sport_events_id"`
	Venue_id        string `json:"venue_id"`
	Period          string `json:"period"`
	Money           string `json:"money"`
	Freq            int    `json:"freq"`
}
type WsRec struct {
	Cmd  string `json:"cmd"`
	Data Conf   `json:"data"`
}

type WsData struct {
	Type string  `json:"type"`
	Sta  Status  `json:"sta"`
	Msg  Message `json:"msg"`
}

// 解析格式 33_146_328099,238100_1000
func (a *App) Lock() {
	a.St.Lock = true
	a.Sc <- 0
}
func (a *App) Unlock() {
	a.St.Lock = false
	a.Sc <- 0
}

func (a *App) Init() error {
	osGwd, _ := os.Getwd()
	path := osGwd + "/.ini"
	log.Println("加载配置文件中...")
	cfg, err := ini.Load(path)
	if err != nil {
		a.LogErr(err.Error())
		return err
	}
	log.Println("加载成功!")
	log.Println("初始化参数...")
	cfgDate := cfg.Section("main").Key("date").String()
	if len([]rune(string(cfgDate))) == 10 {
		a.Date = cfgDate
	} else {
		a.Date = utils.GetTomorrowDate()
	}
	a.St.Key = cfg.Section("").Key("key").String()
	a.St.Id = cfg.Section("").Key("id").String()
	a.St.Period = cfg.Section("main").Key("period").String()
	a.St.Sport_events_id = cfg.Section("main").Key("sport_events_id").String()
	a.St.Money = cfg.Section("main").Key("money").String()
	a.St.Venue_id = cfg.Section("main").Key("venue_id").String()
	a.St.Freq, _ = strconv.Atoi(cfg.Section("main").Key("freq").String())

	a.Api = cfg.Section("main").Key("api").String()
	a.OApi = cfg.Section("main").Key("oapi").String()
	a.Cookie = cfg.Section("main").Key("cookie").String()
	a.Token = cfg.Section("main").Key("token").String()
	a.Body = map[string]string{
		"orderid":         "",
		"card_id":         "",
		"sport_events_id": a.St.Sport_events_id,
		"money":           a.St.Money,
		"ordertype":       "makeappointment",
		"paytype":         "bitpay",
		"openid":          "",
	}
	// a.Body["scene"] = "[{\"day\":\"" + utils.GetTomorrowDate() + "\",\"fields\":{\"" + a.St.Venue_id + "\":[" + a.St.Period + "]}}]"

	return nil
}

func (a *App) Clear() {
	a.St.Form = false
	a.St.Oct = ""
	a.St.Oid = ""
	a.St.Pct = ""
	a.St.Pay = ""
	a.Form = ""
}
