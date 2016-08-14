package main

import (
	G "github.com/gilmour-libs/gilmour-e-go"
	"gopkg.in/gilmour-libs/gilmour-e-go.v4/backends"
	"log"
	"time"
)

func echoEngine(host_port string) *G.Gilmour {
	redis := backends.MakeRedis(host_port, "")
	engine := G.Get(redis)
	return engine
}

func ExecuteRequest(request *G.RequestComposer, count *int) {
	req_msg := G.NewMessage()
	resp, err := request.Execute(req_msg)

	if resp == nil {
		log.Println("nil response due to ", err)
	}

	msg := resp.Next()

	var data string
	msg.GetData(&data)
	log.Println(data)
	atomic.AddInt64(counter, 1)
}

func main() {
	e1 := echoEngine("127.0.0.1:7000")
	e1.Start()

	request1 := e1.NewRequest("test.handler.one")

	count := 0

	for i := 0; i < 50; i += 1 {
		go func(count *int64, request *G.RequestComposer) {
			for {
				ExecuteRequest(request, count)
			}
		}(&count, request1)
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("Total requests processed in 10 seconds = ", count)
			count = 0
		}
	}
}
