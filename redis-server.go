package main

import (
	"fmt"
	G "gopkg.in/gilmour-libs/gilmour-e-go.v4"
	"gopkg.in/gilmour-libs/gilmour-e-go.v4/backends"
	"log"
	"sync"
)

func echoEngine() *G.Gilmour {
	// redis := backends.MakeRedisSentinel("mymaster", "", []string{":16380", ":16381", ":16382"})
	redis := backends.MakeRedis("127.0.0.1:6379", "")
	engine := G.Get(redis)
	return engine
}

//Fetch a remote file from the URL received in Request.
func fetchReply(g *G.Gilmour) func(req *G.Request, resp *G.Message) {
	return func(req *G.Request, resp *G.Message) {
		var data string
		if err := req.Data(&data); err != nil {
			panic(err)
		}

		//Send back the contents.
		log.Println("Responding to ", data)
		resp.SetData(fmt.Sprint("%s_bar", data))
	}
}

//Bind all service endpoints to their topics.
func bindListener(g *G.Gilmour) {
	g.ReplyTo("test.handler.one", fetchReply(g), nil)
}

func main() {
	engine := echoEngine()
	bindListener(engine)
	engine.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
