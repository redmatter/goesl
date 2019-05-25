// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package main

import (
	"flag"
	"fmt"
	"github.com/0x19/goesl"
	"github.com/0x19/goesl/log"
	"runtime"
	"strings"
)

var (
	fshost   = flag.String("fshost", "localhost", "FreeSWITCH hostname. Default: localhost")
	fsport   = flag.Uint("fsport", 8021, "FreeSWITCH port. Default: 8021")
	password = flag.String("pass", "ClueCon", "FreeSWITCH password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "FreeSWITCH connection timeout in seconds. Default: 10")
)

// Small client that will first make sure all events are returned as JSON and second, will originate
func main() {
	goesl.SetLogger(
		log.NewLogger())

	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := goesl.NewClient(*fshost, *fsport, *password, *timeout)

	if err != nil {
		goesl.Error("Error while creating new client: %s", err)
		return
	}

	goesl.Debug("Yuhu! New client: %q", client)

	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crucial part in handling incoming messages :)
	go client.Handle()

	_ = client.Send("events json ALL")

	_ = client.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))

	for {
		msg, err := client.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				goesl.Error("Error while reading Freeswitch message: %s", err)
			}
			break
		}

		goesl.Debug("%s", msg)
	}
}
