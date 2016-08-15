package main

import (
	G "../gilmour-e-go"
	"../gilmour-e-go/backends"
	"log"
	"time"
	"sync/atomic"
	"fmt"
)

func echoEngine() *G.Gilmour {
	redis := backends.MakeRedisSentinel("mymaster", "", []string{":16380", ":16381", ":16382"})
	engine := G.Get(redis)
	return engine
}

func ExecuteRequest(request *G.RequestComposer, msg string, counter *int64) {
	req_msg := G.NewMessage().SetData(msg)
	resp, err := request.Execute(req_msg)

	if resp == nil || err != nil{
		log.Println("Error Occured: ", err)
		return
	}

	atomic.AddInt64(counter, 1)
}

func main() {
	e := echoEngine()
	e.Start()

	request := e.NewRequest("test.handler.one")

	var counter int64

	for i := 0; i < 1; i += 1 {
		go func(counter *int64, request *G.RequestComposer, i int) {
			j := 0
			for {
				payload := fmt.Sprintf("foo_%d_%d", i, j)

				ExecuteRequest(request, payload, counter)
				j += 1
			}
		}(&counter, request, i)
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			log.Println("Total requests processed in 10 seconds = ", counter)
			counter = 0
		}
	}
}
