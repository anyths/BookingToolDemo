package objs

func (a *App) ParseCmd(msg WsRec) {
	if msg.Cmd == "break" {
		if !a.St.Lock {
			a.Bc <- 0
		}
		return
	}
	if msg.Cmd == "sta" {
		a.Sc <- 0
		return
	}
	if !a.St.Status {
		a.LogErr("请先中断任务! ")
		a.Sc <- 0
		return
	}
	a.Cc <- msg
}
