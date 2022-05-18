package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type StatusList map[string]Status

type Status struct {
	Key             string `json:"key"`
	Id              string `json:"id"`
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
	Money           string `json:"money"`
	Venue_id        string `json:"venue_id"`
	Period          string `json:"period"`
	Freq            int    `json:"freq"`
}

type WsCmd struct {
	Cmd  string `json:"cmd"`
	Data Conf   `json:"data"`
}

// 脚本端的链接数据库
type WsDb map[string]Ws

// 脚本和用户端的ws连接模型
type Ws struct {
	Client    *websocket.Conn
	SentChan  chan []byte
	ReadChan  chan []byte
	HeartChan chan []byte
}

// 脚本返回的数据格式
type WsData struct {
	Type string  `json:"type"`
	Sta  Status  `json:"sta"`
	Msg  Massage `json:"msg"`
}

type Massage struct {
	Ok      bool   `json:"ok"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

// 新建一个ws
func (db WsDb) GetWs(symbol string) error {
	if _, ok := db[symbol]; ok {
		return nil
	}
	return errors.New("脚本未连接")
}

func (w *WsDb) Enable(symbol string, client *websocket.Conn) {
	var rst Ws
	rst.Client = client
	rst.SentChan = make(chan []byte)
	rst.ReadChan = make(chan []byte)
	rst.HeartChan = make(chan []byte)
	(*w)[symbol] = rst
}

func (w *WsDb) Disable(symbol string) {
	delete((*w), symbol)
}

func main() {
	r := gin.Default()

	// ws储存

	staDb := map[string]StatusList{}
	wsDb := WsDb{}
	uwsDb := UserWsDb{}
	stas := r.Group("/sta")
	{
		stas.GET("/:key", func(ctx *gin.Context) {
			key := ctx.Param("key")
			if rst, ok := staDb[key]; ok {
				ctx.JSON(200, gin.H{
					"statusCode": 200,
					"data":       rst,
				})
				return
			}
			ctx.JSON(204, gin.H{
				"statusCode": 204,
				"data":       StatusList{},
			})
		})
		stas.GET("/:key/:id", func(ctx *gin.Context) {
			key := ctx.Param("key")
			id := ctx.Param("id")

			if rst, ok := staDb[key][id]; ok {
				ctx.JSON(200, gin.H{
					"statusCode": 200,
					"data":       rst,
				})
				return
			}
			ctx.JSON(204, gin.H{
				"statusCode": 204,
				"data":       Status{},
			})
		})
	}
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // use default options
	r.GET("/ws/:key/:id", func(ctx *gin.Context) {
		key := ctx.Param("key")
		id := ctx.Param("id")
		symbol := key + "#" + id
		err := wsDb.GetWs(symbol)
		if err == nil {
			return
		}
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		wsDb.Enable(symbol, conn)    // 创建一个脚本端ws连接
		uwsDb.EnableUserPool(symbol) // 为用户端创建连接池
		defer uwsDb.DisableUserPool(symbol)
		defer wsDb.Disable(symbol)
		defer delete(staDb[key], id)
		defer conn.Close()

		ErrChan := make(chan int)
		defer close(ErrChan)

		go func() {
			for {
				_, message, err := wsDb[symbol].Client.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					ErrChan <- 0
					return
				}
				wsDb[symbol].ReadChan <- message
				// fmt.Println("消息发送了")
			}
		}()
		// 计时器，用来验证心跳
		ticker := time.NewTicker(time.Second * 12)
		heartBeat := true
		defer ticker.Stop()
		for {
			select {
			case cmd := <-wsDb[symbol].SentChan:
				err = wsDb[symbol].Client.WriteMessage(websocket.TextMessage, cmd)
				if err != nil {
					log.Println(err)
					return
				}
			case <-ticker.C:
				// 超时检测，12s检测一次
				if !heartBeat {
					log.Println("心跳验证失败")
					// 如果心跳没跳起来，说明超时了。对方离线了，就退出
					return

				}
				heartBeat = false // 拨回心跳状态
			case msg := <-wsDb[symbol].ReadChan:
				// 过滤心跳
				if len(msg) == 0 {
					heartBeat = true
					err = wsDb[symbol].Client.WriteMessage(websocket.TextMessage, []byte{})
					if err != nil {
						log.Println(err.Error())
						return
					}
					continue
				}
				// 转发给用户端
				for _, v := range uwsDb[symbol] {
					v.SentChan <- msg
				}
				// 解析数据sta
				data := ParseMsg(msg)
				if data.Type == "sta" {
					if _, ok := staDb[key]; !ok {
						staDb[key] = StatusList{}
					}
					staDb[key][id] = data.Sta
				}
			case <-ErrChan:
				return
			}
		}
	})
	r.GET("/uws/:key/:id", func(ctx *gin.Context) {
		key := ctx.Param("key")
		id := ctx.Param("id")
		symbol := key + "#" + id
		if _, ok := uwsDb[symbol]; !ok {
			return
		}
		// 生成唯一uid
		uid := uuid.New().String()
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		uwsDb.Enable(symbol, uid, conn)
		defer uwsDb.Disable(symbol, uid)
		defer conn.Close()

		ErrChan := make(chan int)
		defer close(ErrChan)

		go func() {
			for {
				_, userCmd, err := uwsDb[symbol][uid].Client.ReadMessage()
				if err != nil {
					log.Println("read uws:", err)
					ErrChan <- 0
					return
				}
				// 收到用户的请求发到ReadChan
				uwsDb[symbol][uid].ReadChan <- userCmd
			}
		}()

		// 为用户添加计时器
		ticker := time.NewTicker(time.Second * 12)
		heartBeat := true
		defer ticker.Stop()

		// 消费处理各个通道信息
		for {
			select {
			case <-ticker.C:
				// 心跳判断
				if !heartBeat {
					fmt.Println("心跳错误")
					return // 心跳没有在规定时间内重新激活
				}
				heartBeat = false // 重置heartbeat
			case cmd := <-uwsDb[symbol][uid].ReadChan:
				// 过滤心跳
				if len(cmd) == 0 {
					heartBeat = true
					err = uwsDb[symbol][uid].Client.WriteMessage(websocket.TextMessage, []byte(""))
					if err != nil {
						log.Println(err.Error())
						return
					}
					continue
				}
				// 转发给脚本端
				wsDb[symbol].SentChan <- cmd
			case data := <-uwsDb[symbol][uid].SentChan:
				err = uwsDb[symbol][uid].Client.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println(err)
					return
				}
			case <-ErrChan:
				return
			}

		}
	})
	r.Run(":9999")
}

func ParseMsg(data []byte) WsData {
	var msg WsData
	json.Unmarshal(data, &msg)
	return msg
}
