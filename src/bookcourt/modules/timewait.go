package modules

import (
	"bookcourt/objs"
	"bookcourt/utils"
	"time"
)

func TimeWait(l *utils.Logs, c *objs.Conf) {
	for {
		if time.Now().Hour() == 0 {
			l.Add("开始配置更新...")
			c.Date = utils.GetTomorrowDate()
			c.Body["scene"] = "[{\"day\":\"" + utils.GetTomorrowDate() + "\",\"fields\":{\"" + c.Venue_id + "\":[" + c.Period + "]}}]"
			l.Add("配置更新成功!")
			time.Sleep(time.Minute * 30)
			continue
		}
		if time.Now().Hour() != 6 {
			l.Add("** 休眠阶段 ** ")
			time.Sleep(time.Minute * 30)
			continue
		}
		if time.Now().Minute() != 59 {
			l.Add("** 预热阶段 ** ")
			time.Sleep(time.Second * 30)
			continue
		}
		if time.Now().Second() != 59 {
			l.Add("** 加速阶段 **")
			time.Sleep(time.Millisecond * 200)
			continue
		}
		break
	}
}
