package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Connection struct {
	channelName string
	wsConn      *websocket.Conn
	//读取websocket的channel
	outChan  chan []byte
	mutex    sync.Mutex
	isClosed bool
}

type RegisterMap struct {
	sync.Map
}

type ConnSlice struct {
	mu    sync.Mutex
	conns []*Connection
}

type Message struct {
	ChannelName string `json:"channelName"`
	Content     string `json:"content"`
	MessageType int    `json:"messageType"`
}

// var registerMap sync.Map

var (
	upgrade = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var registerMap RegisterMap

var messageCh chan Message

func (conn *Connection) Close() {

	conn.wsConn.Close()
	registerMap.RemoveConn(conn)
	//只执行一次
	conn.mutex.Lock()
	if !conn.isClosed {
		conn.isClosed = true
	}
	// Om.Mmx.Unlock()
	conn.mutex.Unlock()
}

func (rm *RegisterMap) AddConn(conn *Connection) {
	conns, ok := rm.Load(conn.channelName)
	if !ok {
		connSlice := new(ConnSlice)
		connSlice.Add(conn)
		rm.Store(conn.channelName, connSlice)
	} else {
		connSlice := conns.(*ConnSlice)
		connSlice.Add(conn)
		rm.Store(conn.channelName, connSlice)
	}
}
func (rm *RegisterMap) RemoveConn(conn *Connection) {
	conns, ok := rm.Load(conn.channelName)
	if !ok {
		return
	} else {
		connSlice := conns.(*ConnSlice)
		connSlice.Remove(conn)
		rm.Store(conn.channelName, connSlice)
	}
}

func (cs *ConnSlice) Add(conn *Connection) {
	cs.mu.Lock()
	cs.conns = append(cs.conns, conn)
	cs.mu.Unlock()
}

func (cs *ConnSlice) Remove(conn *Connection) {
	cs.mu.Lock()
	key := -1
	for k, v := range cs.conns {
		if v == conn {
			key = k
			break
		}
	}
	if key >= 0 {
		cs.conns = append(cs.conns[:key], cs.conns[key:]...)
	}
	cs.mu.Unlock()
}

func ChannelRegister(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelName := params["channelName"]
	fmt.Print(channelName)
	wsConn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("fatal error: ", err)
		return
	}
	conn := Connection{channelName: channelName, wsConn: wsConn}
	registerMap.AddConn(&conn)
	go writeLoop(&conn)
}

func MessagePublish(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	var message Message
	err = json.Unmarshal([]byte(body), &message)
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}
	messageCh <- message
	io.WriteString(w, "accept")
}

func MessageReceive() {
	var message Message
	for {
		message = <-messageCh
		connSlice, ok := registerMap.Load(message.ChannelName)
		if ok {
			conns := connSlice.(*ConnSlice).conns
			for _, v := range conns {
				// if err := v.wsConn.WriteMessage(websocket.TextMessage, []byte(message.Content)); err != nil {
				// 	log.Fatalln("fatal error: ", err)
				// }
				v.outChan <- []byte(message.Content)
			}
		}
	}
}

func writeLoop(conn *Connection) {
	var (
		err error
	)
	for {
		data, ok := <-conn.outChan
		if ok {
			fmt.Println("11")
			if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
				goto ERR
			}
		}
	}
ERR:
	conn.Close()
}

func main() {
	messageCh = make(chan Message)
	go MessageReceive()
	router := mux.NewRouter()
	router.HandleFunc("/mq/channelRegister/{channelName}/", ChannelRegister)
	router.HandleFunc("/mq/messagePublish/", MessagePublish)
	// router.HandleFunc("/mq/publish/", timingSendMessage)
	err := http.ListenAndServe("0.0.0.0:9630", router)
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}
	fmt.Fprintf(os.Stdout, "%s", "start connection")

}
