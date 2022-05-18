package main

import "github.com/gorilla/websocket"

// 用户uws模型

type UserWsDb map[string]UserPool

type UserPool map[string]Ws

func (w *UserWsDb) Enable(symbol string, uuid string, client *websocket.Conn) {
	var rst Ws
	rst.Client = client
	rst.SentChan = make(chan []byte)
	rst.ReadChan = make(chan []byte)
	rst.HeartChan = make(chan []byte)
	(*w)[symbol][uuid] = rst
}
func (w *UserWsDb) EnableUserPool(symbol string) {
	(*w)[symbol] = UserPool{}
}

// 删除一个用户连接
func (w *UserWsDb) Disable(symbol string, uuid string) {
	delete((*w)[symbol], uuid)
}

// 删除所有用户连接
func (w *UserWsDb) DisableUserPool(symbol string) {
	delete((*w), symbol)
}
