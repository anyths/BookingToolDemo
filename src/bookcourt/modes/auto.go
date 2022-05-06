package modes

import (
	"bookcourt/modules"
	"bookcourt/objs"
	"bookcourt/utils"
)

func Auto(l *utils.Logs, c *objs.Conf, mode string) bool {
	lock := false
	// 计时
	if mode == "0" {
		modules.TimeWait(l, c)
	}
	if mode == "2" {
		lock = true
	}
	var form string
	l.Add("启动执行...")
	form, err := modules.CreateOrder(l, c, lock)
	if err != nil {
		l.Add(err.Error())
		l.Add("进行订单判定...")

		orderid, err := modules.FindOrder(l, *c)
		if err != nil {
			l.Err(err.Error())
			l.Err("抢单失败!")
			return false
		}
		l.Add("查询到订单...")
		l.Add("开始重新生成支付...")

		for {
			rst, err := modules.NewOrder(orderid, *c)
			if err != nil {
				l.Err(err.Error())
				l.Err("服务器响应超时: 开始重试...")
				continue
			}
			if rst.Code == 0 {
				l.Err(rst.Msg)
				l.Add("建议使用: 2 - 锁定模式 二次抢单")
				return false
			}
			form = rst.Data
			break
		}
	}
	l.Add("订单已生成...")
	l.Add("进入支付流程...")

	payForm, err := modules.ParseForm(form)
	if err != nil {
		l.Err(err.Error())
		l.Err("解析表单时出错...")
		l.Add("请尝试手动前往操作付款")
		return false
	}
	l.Add("解析支付表单成功! ")
	l.Add("解析 BitPay 支付表单")
	prepForm, err := modules.GetBitPay(payForm)
	if err != nil {
		l.Err(err.Error())
		l.Err("解析 北理工支付系统 表单时出错...")
		l.Add("请尝试手动前往支付")
		return false
	}
	l.Add("BitPay 支付表单解析成功...")
	l.Add("开始生成微信支付链接...")
	payLink, err := modules.GetWxPay(prepForm)
	if err != nil {
		l.Err(err.Error())
		l.Err("获取支付链接失败...")
		l.Add("请尝试手动付款")
		return false
	}
	l.Add("解析微信支付成功!")
	l.Add("支付链接已生成\n\n支付链接: \n" + payLink + "\n\n使用方式: \n1. 手机浏览器打开 wx.gt0.cn\n2. 将此链接粘贴\n3. 点击“调用微信支付”\n4. 支付成功后，点击“验证订单状态”\n")
	return true
}
