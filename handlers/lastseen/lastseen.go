package lastseen

import (
	"fmt"
	"idoink"
	"sync"
	"time"
)

const LastSeenCmd = "seen"

// TODO: consider if this is still needed
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

// TODO boltdb persistence

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

func LastSeen(e *idoink.E) (bool, error) {
	if len(e.Rest) < 1 {
		return false, nil
	}

	n := e.Rest[0]
	lastSeenMx.Lock()
	ls, ok := lastSeenState[n]
	lastSeenMx.Unlock()

	if !ok {
		e.I.Message(e.To, fmt.Sprintf("%s: lastseen - no info for %s", e.From, n))
	} else {
		date := ls.when // todo : nice format
		c := ls.channel
		m := ls.msg
		e.I.Message(e.To, fmt.Sprintf("%s: lastseen - %s last seen on %s in %s saying %s",
			e.From, n, date, c, m))
	}

	return false, nil
}
