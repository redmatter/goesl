// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package main

import (
	"fmt"
	"github.com/0x19/goesl"
	"github.com/0x19/goesl/log"
	"os"
	"runtime"
	"strings"
)

var welcomeFile = "%s/media/welcome.wav"

func main() {
	goesl.SetLogger(
		log.NewLogger())

	defer func() {
		if r := recover(); r != nil {
			goesl.Error("Recovered in f", r)
		}
	}()

	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	wd, err := os.Getwd()

	if err != nil {
		goesl.Error("Error while attempt to get WD: %s", wd)
		os.Exit(1)
	}

	welcomeFile = fmt.Sprintf(welcomeFile, wd)

	if s, err := goesl.NewOutboundServer(":8084"); err != nil {
		goesl.Error("Got error while starting Freeswitch outbound server: %s", err)
	} else {
		go handlePlayback(s)
		_ = s.Start()
	}

}

// handlePlayback - Running under goroutine here to explain how to handle playback ( play to the caller )
func handlePlayback(s *goesl.OutboundServer) {

	for {

		select {

		case conn := <-s.Conns:
			goesl.Notice("New incoming connection: %v", conn)

			if err := conn.Connect(); err != nil {
				goesl.Error("Got error while accepting connection: %s", err)
				break
			}

			answer, err := conn.ExecuteAnswer("", false)

			if err != nil {
				goesl.Error("Got error while executing answer: %s", err)
				break
			}

			goesl.Debug("Answer Message: %s", answer)
			goesl.Debug("Caller UUID: %s", answer.GetHeader("Caller-Unique-Id"))

			cUUID := answer.GetCallUUID()

			if sm, err := conn.Execute("playback", welcomeFile, true); err != nil {
				goesl.Error("Got error while executing playback: %s", err)
				break
			} else {
				goesl.Debug("Playback Message: %s", sm)
			}

			if hm, err := conn.ExecuteHangup(cUUID, "", false); err != nil {
				goesl.Error("Got error while executing hangup: %s", err)
				break
			} else {
				goesl.Debug("Hangup Message: %s", hm)
			}

			go func() {
				for {
					msg, err := conn.ReadMessage()

					if err != nil {

						// If it contains EOF, we really dont care...
						if !strings.Contains(err.Error(), "EOF") {
							goesl.Error("Error while reading Freeswitch message: %s", err)
						}
						break
					}

					goesl.Debug("%s", msg)
				}
			}()

		default:
			// YabbaDabbaDooooo!
			//Flintstones. Meet the Flintstones. They're the modern stone age family. From the town of Bedrock,
			// They're a page right out of history. La la,lalalalala la :D
		}
	}

}
