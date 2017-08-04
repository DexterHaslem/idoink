package main

import (
	"log"
	"strings"
)

const botMagic = "^bot"

// parseMsg is root level handler, it handles
// protocol level messages as well as ferrying the
// higher level messages off to subparsers
func parseMsg(msg string) {
	//log.Printf("parseMsg: %s\n", msg)
	chunks := strings.Split(msg, " ")
	switch chunks[0] {
	case "PING":
		// chunk 1 contains the sep :12312412341
		challenge := chunks[1] //strings.Replace(chunks[1], ":", "", -1)
		//i.Pong(challenge)
		i.Pong(challenge[1:])

		// we get this after register, join our chans
		for _, c := range i.Chans {
			i.Join(c)
		}
		break
	case "NOTICE":
		switch chunks[1] {
		case "AUTH":
			if chunks[len(chunks)-1] == "response" {
				i.Register()
			}
			break
		}
	default:
		if chunks[0][0] == byte(':') {
			tryServerMessage(chunks)
		}
		break
	}
}

func tryServerMessage(chunks []string) {
	// qnet in particular sends these server messages
	// with codes, eg
	//:cymru.us.quakenet.org 001
	if len(chunks) < 2 {
		return
	}

	switch chunks[1] {
	case "005":
		log.Printf("Server CAPS: %v", chunks[3:])
		break
	case "422":
		// MOTD
		break
	case "NOTICE":
		// on qnet these are some of the last statuslines to come in
		// server will parrot notice to nick, send joins after htis
		for _, c := range i.Chans {
			i.Join(c)
		}
	case "PRIVMSG":
		privMsg(chunks[0], chunks[2], chunks[3:]...)
		break
	}
}

func privMsg(from, to string, rest ...string) {
	//:dmh!sid189360@id-189360.charlton.irccloud.com PRIVMSG #warsow.na :haa ownage
	fn := strings.Split(from, "!")[0][1:]
	msg := strings.Join(rest, " ")[1:]
	log.Printf("privmsg from %s to %s: %s\n", fn, to, msg)

	// so first thing we check for is bot magic
	trimmed := rest[0][1:]
	if trimmed != botMagic {
		return
	}

	// next message must be the bot specific command
	cmd := rest[1]
	params := []string{}
	if len(rest) >= 3 {
		params = rest[2:]
	}
	handleBotMsg(fn, to, cmd, params...)
}

func handleBotMsg(from, to, cmd string, chunks ...string) {
	// maybe consider to containing only chans, kind of redundant tho
	switch cmd {
	case ddgCmd:
		ddg(from, to, chunks...)
		break
	case lastfmCmd:
		lastfm(from, to, chunks...)
		break
	}
}
