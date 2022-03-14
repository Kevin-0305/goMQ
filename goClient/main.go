package main

import (
	"fmt"
	"goClient/goClient"
	"sync"
)

func onMessage(msg string) {
	fmt.Println("func:", msg)
}

func main() {
	wsc := goClient.NewWsClient("127.0.0.1", "9630", "123", 10)
	wsc.Start(onMessage)
	var w1 sync.WaitGroup
	w1.Add(1)
	w1.Wait()
}
