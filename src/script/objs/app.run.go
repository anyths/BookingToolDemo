package objs

import (
	"bookcourt/utils"
	"context"
	"errors"
	"time"
)

func (a *App) RunAuto(mode int) (string, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if mode == 0 {
		// 1模式定时
		success := a.TimeWait()
		a.Lock()
		if !success {
			return "", errors.New("用户终止")
		}
	}
	lock := false
	if mode == 3 {
		// 3模式锁单
		lock = true
	}
	// 6点59分59秒放行
	// 并发创建订单
	form, err := a.OrderCreate(lock, ctx)
	if err != nil {
		if err.Error() == "用户中断" {
			return "", err
		}
		a.LogErr(err.Error())
		a.LogAdd("开始二次订单判定")

		// 查询订单列表验证同订单
		orderid, err := a.OrderFind()
		if err != nil {
			a.LogErr(err.Error())
			return "", err
		}
		// 生成支付表单
		a.LogAdd("重新生成支付申请表单")

		// 写for循环，失败重试
		for {
			rst, err := a.OrderNewForm(orderid)
			if err != nil {
				a.LogErr(err.Error())
				continue
			}
			if rst.Code == 0 {
				a.LogErr(rst.Msg)
				a.LogAdd("建议使用: 2 - 锁单模式 二次抢单")
				return "", errors.New("订单已过支付期: 尝试[锁单]二次抢单")
			}
			form = rst.Data
			break
		}
	}
	return form, nil
}

func (a *App) RunPay() (string, error) {
	payLink := ""
	if a.Form == "" {
		return "", errors.New("空表单, 请[下单]! ")
	}
	payForm, err := FormParse(a.Form)
	if err != nil {
		return payLink, err
	}
	a.LogAdd("解析支付表单成功!")
	a.LogAdd("开始解析 BitPay 表单...")
	bitPayForm, err := FormBitPay(payForm)
	if err != nil {
		return payLink, err
	}
	a.LogAdd("BitPay 支付表单解析成功!")
	a.LogAdd("开始解析微信支付链接...")
	payLink, err = GetWxPay(bitPayForm)
	if err != nil {
		return payLink, err
	}
	return payLink, nil
}

func (a *App) TimeWait() bool {
	a.Unlock()
	uptime := false
	for {
		select {
		case <-a.Bc:
			return false
		default:
			if time.Now().Hour() == 0 {
				if !uptime {
					a.LogAdd("[auto]开始配置更新...")
					a.Date = utils.GetTomorrowDate()
					a.Body["scene"] = "[{\"day\":\"" + utils.GetTomorrowDate() + "\",\"fields\":{\"" + a.St.Venue_id + "\":[" + a.St.Period + "]}}]"
					a.LogAdd("[auto]配置更新成功!")
					uptime = true
				}
				a.LogAdd("[auto]待机阶段")
				time.Sleep(time.Second)
				continue
			}
			if time.Now().Hour() != 6 {
				a.LogAdd("[auto]待机阶段")
				time.Sleep(time.Second)
				continue
			}
			if time.Now().Minute() != 59 {
				a.LogAdd("[auto]预热阶段")
				time.Sleep(time.Second)
				continue
			}
			if time.Now().Second() != 59 {
				a.LogAdd("[auto]加速阶段")
				time.Sleep(time.Millisecond * 200)
				continue
			}
			return true
		}
	}

}
