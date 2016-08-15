// Copyright 2012-2016 Apcera Inc. All rights reserved.
// +build ignore

package main

import (
	"fmt"
	"github.com/nats-io/nats"
	"log"
	"runtime"
)

func main() {
	urls := "nats://127.0.0.1:4222, nats://127.0.0.1:5222, nats://127.0.0.1:6222"

	nc, err := nats.Connect(urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	subj, reply, i := "testsubject", "bar", 0

	nc.QueueSubscribe(subj, "worker_group", func(msg *nats.Msg) {
		i++
		nc.Publish(msg.Reply, []byte(fmt.Sprintf("%s_%s", msg.Data, reply)))
	})

	log.Printf("Listening on [%s]\n", subj)

	runtime.Goexit()
}
