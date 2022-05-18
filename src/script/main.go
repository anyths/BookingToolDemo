package main

import (
	"bookcourt/objs"
	"time"
)

func main() {
	var (
		app objs.App
		// 状态通道
	)

	app.Host = "wx.gt0.cn"
	// app.Host = "localhost:888"
	app.WsPre = "wss"
	app.HttpPre = "https"
	app.Tapi = app.HttpPre + "://" + app.Host + "/api"

	// 初始化配置
	err := app.Init()
	if err != nil {
		app.LogErr("初始化配置失败! 请检查配置文件.ini")
		time.Sleep(time.Second * 10)
		return
	}

	// ws接收数据
	app.Rc = make(chan []byte)
	// 状态通道
	app.Sc = make(chan int)
	// 命令通道
	app.Cc = make(chan objs.WsRec)
	// 终止程序通道
	app.Bc = make(chan int)
	// 日志通道
	app.Lc = make(chan objs.WsData)

	app.St.Status = true
	app.St.Form = false
	app.St.Pay = ""
	app.St.Pct = ""
	app.St.Oct = ""
	app.St.Cmd = ""

	// app.SyncConfig()

	go app.ProcWs()

	app.Lock()
	app.Sc <- 0

	// 异步同步状态变化，所有状态变化都在sc通道中
	// 过滤后的cmd消费者
	for {
		// 从通道获取命令
		cmdrec := <-app.Cc
		// 同步命令给sta
		app.St.Cmd = cmdrec.Cmd

		app.LogAdd("执行: " + app.St.Cmd)
		if !app.St.Status && app.St.Cmd != "break" {
			app.LogErr("请先停止任务!")
			continue
		}
		app.Sc <- 0
		if app.St.Cmd == "auto" {
			app.Doing()
			// 自动模式
			form, err := app.RunAuto(0)
			if err != nil {
				app.Err(err.Error())
				continue
			}
			app.St.Form = true
			app.Form = form
			app.Done()
		}
		if app.St.Cmd == "now" {
			app.Doing()

			// 立即模式
			form, err := app.RunAuto(1)
			if err != nil {
				app.Err(err.Error())
				continue
			}
			app.St.Form = true
			app.Form = form
			app.Done()
		}
		if app.St.Cmd == "lock" {
			app.Doing()
			// 锁单模式
			form, err := app.RunAuto(2)
			if err != nil {
				app.Err(err.Error())
				continue
			}
			app.St.Form = true
			app.Form = form
			app.Done()
		}
		if app.St.Cmd == "pay" {
			app.Doing()
			// 生成支付
			pay, err := app.RunPay()
			if err != nil {
				app.Err(err.Error())
				continue
			}
			app.St.Pay = pay
			app.Done()
		}
		if app.St.Cmd == "conf" {
			app.Doing()
			err := app.GetConfig(cmdrec.Data)
			if err != nil {
				app.Err(err.Error())
				continue
			}
			app.Done()
		}
		if app.St.Cmd == "clear" {
			app.Doing()
			app.Clear()
			app.Done()
		}
	}
}
