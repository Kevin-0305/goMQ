package goClient

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type onMessage func(string)
type onClose func(string)
type onError func(string)

type websocketClient struct {
	conn        *websocket.Conn
	addr        *string
	path        string
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int
	onMessage   *onMessage
	onClose     *onClose
	onError     *onError
}

// 构造函数
func NewMQClient(addrIp string, addrPort string, timeout int) *websocketClient {
	addrString := addrIp + ":" + addrPort
	var sendChan = make(chan string, 10)
	var recvChan = make(chan string, 10)
	var conn *websocket.Conn
	path := "/mq/channelRegister/"
	return &websocketClient{
		addr:        &addrString,
		path:        path,
		conn:        conn,
		sendMsgChan: sendChan,
		recvMsgChan: recvChan,
		isAlive:     false,
		timeout:     timeout,
	}
}

// 链接服务端
func (wsc *websocketClient) dail() {
	var err error
	u := url.URL{Scheme: "ws", Host: *wsc.addr, Path: wsc.path}
	log.Printf("connecting to %s", u.String())
	wsc.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return

	}
	wsc.isAlive = true
	log.Printf("connecting to %s 链接成功！！！", u.String())

}

// 发送消息
func (wsc *websocketClient) sendMsg() {
	go func() {
		for {
			msg := <-wsc.sendMsgChan
			err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("write:", err)
				continue
			}
		}
	}()
}

// 读取消息
func (wsc *websocketClient) readMsg() {
	go func() {
		for {
			if wsc.conn != nil {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					wsc.isAlive = false
					// 出现错误，退出读取，尝试重连
					break
				}
				log.Printf("recv: %s", message)
				// 需要读取数据，不然会阻塞
				wsc.recvMsgChan <- string(message)
			}

		}
	}()
}

// 开启服务并重连
func (wsc *websocketClient) Start() {
	for {
		if wsc.isAlive == false {
			wsc.dail()
			wsc.sendMsg()
			wsc.readMsg()
		}
		time.Sleep(time.Second * time.Duration(wsc.timeout))
	}
}
