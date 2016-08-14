package main

import (
	"fmt"
	G "github.com/gilmour-libs/gilmour-e-go"
	"gopkg.in/gilmour-libs/gilmour-e-go.v4/backends"
	"log"
	"sync"
)

func echoEngine(host_port string) *G.Gilmour {
	redis := backends.MakeRedis("127.0.0.1:6379", "")
	engine := G.Get(redis)
	return engine
}

//Fetch a remote file from the URL received in Request.
func fetchReply(g *G.Gilmour) func(req *G.Request, resp *G.Message) {
	return func(req *G.Request, resp *G.Message) {
		var num int
		if err := req.Data(&num); err != nil {
			panic(err)
		}

		//Send back the contents.
		log.Println("Responding to ", num)
		resp.SetData(fmt.Sprint("Response of %d", num))
	}
}

//Bind all service endpoints to their topics.
func bindListener(g *G.Gilmour) {
	g.ReplyTo("test.handler.one", fetchReply(g), nil)
}

func main() {
	engine := echoEngine("127.0.0.1:7000")
	bindListener(engine)
	engine.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
