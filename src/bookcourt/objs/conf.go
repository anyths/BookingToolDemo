package objs

import (
	"bookcourt/utils"
	"os"
	"strconv"

	"gopkg.in/ini.v1"
)

type Conf struct {
	// 下单api
	Api             string
	OApi            string
	Cookie          string
	Token           string
	Sport_events_id string
	// 总金额
	Money string
	// 场地编号id
	Venue_id string
	// 时间段编号，多个以,隔开
	Period string
	// 频率
	Freq int
	// 日期
	Date string
	Body map[string]string
}

func (c *Conf) Init(l *utils.Logs) error {
	osGwd, _ := os.Getwd()
	path := osGwd + "/.ini"
	l.Add("加载配置文件中...")
	cfg, err := ini.Load(path)
	if err != nil {
		l.Err(err.Error())
		return err
	}
	l.Add("加载成功!")
	l.Add("初始化参数...")
	cfgDate := cfg.Section("main").Key("date").String()
	if len([]rune(string(cfgDate))) == 10 {
		c.Date = cfgDate
	} else {
		c.Date = utils.GetTomorrowDate()
	}
	c.Api = cfg.Section("main").Key("api").String()
	c.OApi = cfg.Section("main").Key("oapi").String()
	c.Cookie = cfg.Section("main").Key("cookie").String()
	c.Token = cfg.Section("main").Key("token").String()
	c.Period = cfg.Section("main").Key("period").String()
	c.Sport_events_id = cfg.Section("main").Key("sport_events_id").String()
	c.Money = cfg.Section("main").Key("money").String()
	c.Venue_id = cfg.Section("main").Key("venue_id").String()
	c.Freq, _ = strconv.Atoi(cfg.Section("main").Key("freq").String())
	c.Body = map[string]string{
		"orderid":         "",
		"card_id":         "",
		"sport_events_id": c.Sport_events_id,
		"money":           c.Money,
		"ordertype":       "makeappointment",
		"paytype":         "bitpay",
		"openid":          "",
	}
	c.Body["scene"] = "[{\"day\":\"" + utils.GetTomorrowDate() + "\",\"fields\":{\"" + c.Venue_id + "\":[" + c.Period + "]}}]"

	return nil
}
