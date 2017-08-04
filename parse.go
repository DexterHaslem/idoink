package main

import (
	"idoink/irc"
	"log"
	"strings"
)

func parseMsg(msg string, i *irc.IRC) {
	log.Printf("parseMsg: %s\n", msg)

	chunks := strings.Split(msg, " ")
	switch chunks[0] {
	case "PING":
		// chunk 1 contains the sep :12312412341
		challenge := chunks[1] //strings.Replace(chunks[1], ":", "", -1)
		//i.Pong(challenge)
		i.Pong(challenge[1:])
		break
	case "NOTICE":
		if chunks[1] == "AUTH" {
			//i.Register()
			if chunks[len(chunks)-1] == "response" {
				i.Register()
			}
		}
		break
	}
}
