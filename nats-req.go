// Copyright 2012-2016 Apcera Inc. All rights reserved.
// +build ignore

package main

import (
	// "flag"
	"fmt"
	"github.com/nats-io/nats"
	"log"
	"sync/atomic"
	"time"
)

func main() {
	urls := "nats://127.0.0.1:4222, nats://127.0.0.1:5222, nats://127.0.0.1:6222"
	nc, err := nats.Connect(urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}
	defer nc.Close()

	var counter int64

	for i := 0; i < 1000; i += 1 {
		go func(i int, sub string, payload_arg string, nc *nats.Conn, counter *int64) {
			j := 1
			for {
				raw_payload := fmt.Sprintf("%s_%d_%d", payload_arg, i, j)
				subj, payload := sub, []byte(raw_payload)

				msg, err := nc.Request(subj, []byte(payload), 1000*time.Millisecond)
				if err != nil {
					log.Fatalf("Error in Request: %v\n", err)
				} else if string(msg.Data) == raw_payload+"_bar" {
					atomic.AddInt64(counter, 1)
				}
				j += 1
			}
		}(i, "testsubject", "payload", nc, &counter)
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("Total requests processed = ", counter)
			counter = 0
		}
	}
}
