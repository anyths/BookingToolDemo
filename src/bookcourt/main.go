package main

import (
	"bookcourt/modes"
	"bookcourt/modules"
	"bookcourt/objs"
	"bookcourt/utils"
	"fmt"
	"time"
)

func main() {
	var (
		conf objs.Conf
		logs utils.Logs
	)

	// 初始化配置
	err := conf.Init(&logs)
	if err != nil {
		logs.Err("初始化配置失败! 请检查配置文件.ini")
		logs.Write()
		return
	}

	fmt.Println(">>>>>>>>>>>> 当前配置 <<<<<<<<<<<<")
	fmt.Println("预约日期: ", conf.Date)
	fmt.Println("请求间隔: ", conf.Freq, "毫秒")
	fmt.Println("项目编号: ", conf.Sport_events_id)
	fmt.Println("场地编号: ", conf.Venue_id)
	fmt.Println("预约时间段: ", conf.Period)
	fmt.Println("价格: ", conf.Money)

	fmt.Println("--------------------------------------")
	fmt.Println("(!务必配置正确的价格，否则会出错)")
	fmt.Println("======================================")

	// 设置模式
	input := modules.SelectMode(&logs)

	if input == "0" || input == "1" || input == "2" {
		if input == "0" {
			logs.Add("启动: 0 - 自动模式")
		}
		if input == "1" {
			logs.Add("启动: 1 - 立即执行")
		}
		if input == "2" {
			logs.Add("启动: 2 - 锁单模式")
		}

		result := modes.Auto(&logs, &conf, input)
		if result {
			logs.Add("执行全部完成! ")
		} else {
			logs.Err("执行未全部完成!")
		}
		fmt.Println("保存日志中...")
		logs.Write()
		fmt.Println("保存", utils.GetTomorrowDate()+".txt", "成功!")
	}
	defer fmt.Println("脚本已停止, 10s后自动退出...")
	defer time.Sleep(time.Second * 10)

}
