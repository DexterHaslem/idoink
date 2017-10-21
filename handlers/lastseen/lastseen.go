package main

import (
	"fmt"
	"sync"
	"time"
)

const lastSeenCmd = "lastseen"

var lastSeenMx *sync.Mutex

type lastSeenInfo struct {
	when    time.Time
	channel string
	msg     string
}

var lastSeenState map[string]*lastSeenInfo

func init() {
	lastSeenState = map[string]*lastSeenInfo{}
	lastSeenMx = &sync.Mutex{}
}

func updateLastSeen(nick, msg, channel string) {
	s := &lastSeenInfo{
		when:    time.Now(),
		channel: channel,
		msg:     msg,
	}
	lastSeenMx.Lock()
	lastSeenState[nick] = s
	lastSeenMx.Unlock()
}

func lastSeen(from, to string, chunks ...string) {
	if len(chunks) < 1 {
		return
	}

	n := chunks[0]
	lastSeenMx.Lock()
	ls, ok := lastSeenState[n]
	lastSeenMx.Unlock()

	if !ok {
		i.PrivMsg(to, fmt.Sprintf("%s: lastseen - no info for %s", from, n))
	} else {
		date := ls.when // todo : nice format
		c := ls.channel
		m := ls.msg
		i.PrivMsg(to, fmt.Sprintf("%s: lastseen - %s last seen on %s in %s saying %s", from, n, date, c, m))
	}

}
