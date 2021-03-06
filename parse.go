package idoink

import (
	"log"
	"strings"
	"time"
)

const botMagic = "^bot"

// parseMsg is root level handler, it handles
// protocol level messages as well as ferrying the
// higher level messages off to subparsers
func (i *idoink) parseMsg(msg string) {
	//log.Printf("parseMsg: %s\n", msg)
	chunks := strings.Split(msg, " ")
	switch chunks[0] {
	case "PING":
		// chunk 1 contains the sep :12312412341
		challenge := chunks[1] //strings.Replace(chunks[1], ":", "", -1)
		//i.Pong(challenge)
		i.irc.Pong(challenge[1:])

		// we get this after register, join our chans
		for _, c := range i.irc.Chans {
			i.irc.Join(c)
		}
		break
	case "NOTICE":
		switch chunks[1] {
		case "AUTH":
			if chunks[len(chunks)-1] == "response" {
				i.irc.Register()
			}
			break
		}
	default:
		if chunks[0][0] == byte(':') {
			i.tryServerMessage(chunks)
		}
		break
	}
}

func (i *idoink) tryServerMessage(chunks []string) {
	// qnet in particular sends these server messages
	// with codes, eg
	//:cymru.us.quakenet.org 001
	if len(chunks) < 2 {
		return
	}

	// TODO: abstract this somewhere too
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
		for _, c := range i.irc.Chans {
			i.irc.Join(c)
		}
	case "PRIVMSG":
		i.onPrivMsg(chunks[0], chunks[2], chunks[3:]...)
		break
	}
}

func (i *idoink) onPrivMsg(from, to string, rest ...string) {
	// start commands at beginning of msg
	cmd := rest[0][1:] // skip the message part first character because its ':'

	// from starts with ':' as well, trim it, and is full netmask

	fullnm := from[1:]
	fromchunks := strings.Split(fullnm, "!")
	nick := fromchunks[0]
	netmask := fromchunks[1]

	// construct one new event for all handlers.
	// if 3rd party handler modify it, ¯\_(ツ)_/¯
	e := &E{
		Cmd:      cmd,
		From:     nick,
		FullFrom: fullnm,
		Netmask:  netmask,
		To:       to,
		Rest:     rest[1:], // skip cmd
		I:        i,
		Time:     time.Now(),
	}

	i.handleBotMsg(e) //fn, to, cmd, params...)
}

func (i *idoink) handleBotMsg(e *E) { //from, to, cmd string, rest ...string) {
	// nothing to run
	if i.hc < 1 {
		return
	}

	run := []*hm{}
	// filter through handlers to see what we need to run
	// beware, iterating map is nondeterministic in go, so loop by id manually
	// this will preserve order they registered in in case a handler requests
	// to stop processing
	for idx := 0; idx < len(i.handlers); idx++ {
		hi, ok := i.handlers[idx+1]
		if !ok {
			break
		}
		if hi.cmd == "" || hi.cmd == e.Cmd {
			run = append(run, hi)
		}
	}

	if len(run) < 1 {
		i.irc.PrivMsg(e.To, "unknown command")
		return
	}

	for _, hi := range run {

		stop, err := hi.h(e)

		if err != nil {
			log.Printf("error on %s handler(#%d): %s\n", hi.cmd, hi.id, err)
		}

		if stop {
			break
		}
	}
}
