package idoink

import (
	"errors"
	"idoink/irc"
	"log"
	"strings"
	"sync"
)

// is the main idoink type, it contains all settings for the irc bot
type idoink struct {
	irc           *irc.IRC
	hc            int
	handlers      map[int]*hm
	nick          string
	server        string
	chansList     string
	parsedChans   []string
	stopRequested bool
	m             *sync.Mutex
}

// I is an instance of the idoink bot, it has functions to control it
type I interface {

	// These are wrappers above irc which
	// will be exposed to the handlers

	Nick() string
	SetNick(string) error
	Chans() []string
	JoinChan(string) error
	PartChan(string) error
	Message(to string, msg string) error
	Raw(string) error
	// server is fixed tho

	Start() error
	Stop() error
	AddHandler(string, H) (int, error)
	RemoveHandler(int) error
}

// New creates a new IDoink bot.
// chans is a comma separated list of channels to join
func New(nick, server, chans string) I {
	parsedChans := []string{}
	if chans != "" {
		parsedChans = strings.Split(chans, ",")
	}

	return &idoink{
		m:           &sync.Mutex{},
		handlers:    map[int]*hm{},
		nick:        nick,
		server:      server,
		chansList:   chans,
		parsedChans: parsedChans,
	}
}

func (i *idoink) Nick() string {
	return i.nick
}

func (i *idoink) SetNick(nn string) error {
	if err := i.irc.SetNick(nn); err != nil {
		return err
	}
	i.nick = nn
	return nil
}

type errFunc func() error

func chkConnected(i *idoink, f errFunc) error {
	if i.irc == nil {
		return errors.New("not connected")
	}
	return f()
}

func (i *idoink) Chans() []string {
	return i.parsedChans
}

func (i *idoink) JoinChan(c string) error {
	return chkConnected(i, func() error {
		return i.irc.Join(c)
	})
}

func (i *idoink) PartChan(c string) error {
	return chkConnected(i, func() error {
		return i.irc.Part(c)
	})
}

func (i *idoink) Message(to, msg string) error {
	return chkConnected(i, func() error {
		return i.irc.PrivMsg(to, msg)
	})
}

func (i *idoink) Raw(cmd string) error {
	return chkConnected(i, func() error {
		return i.irc.Cmd(cmd)
	})
}

// Start will start the irc bot on a new goroutine, it will not block caller
func (i *idoink) Start() error {
	newIrc, err := irc.New(i.nick, i.server, i.parsedChans)
	if err != nil {
		//log.Fatal(err)
		return err
	}

	i.irc = newIrc

	go func() {
		mc := make(chan string, 5)

		i.irc.Start(func(m string) {
			mc <- m
		}, func(e error) {
			log.Fatal(e)
			close(mc)
		})

		for m := range mc {
			if i.stopRequested {
				break
			}
			i.parseMsg(m)
		}
	}()
	return nil
}

func (i *idoink) Stop() error {
	i.stopRequested = true
	// just slam it
	return i.irc.Close()
}

// AddHandler will add a new hook for every privmsg line received
// it will only fire for a given cmd if it is not an empty string
func (i *idoink) AddHandler(cmd string, h H) (int, error) {
	i.m.Lock()
	defer i.m.Unlock()

	// note: cmd duplicates are allowed currently.
	// add a filter here if they should be unique
	i.hc++
	id := i.hc
	i.handlers[id] = &hm{
		id:  id,
		h:   h,
		cmd: cmd,
	}

	return id, nil
}

func (i *idoink) RemoveHandler(id int) error {
	i.m.Lock()
	defer i.m.Unlock()

	_, ok := i.handlers[id]
	if !ok {
		return errors.New("not found")
	}

	delete(i.handlers, id)
	return nil
}
