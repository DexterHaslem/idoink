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
	Start() error
	Stop() error
	AddHandler(string, H) (int, error)
	RemoveHandler(int) error
}

// New creates a new IDoink bot.
// chans is a comma separated list of channels to join
func New(nick, server, chans string) I {
	return &idoink{
		nick:      nick,
		server:    server,
		chansList: chans,
	}
}

// Start will start the irc bot on a new goroutine, it will not block caller
func (i *idoink) Start() error {
	i.parsedChans = []string{}
	if i.chansList != "" {
		i.parsedChans = strings.Split(i.chansList, ",")
	}

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
	return nil
}

// AddHandler will add a new hook for every privmsg line received
// it will only fire for a given prefix if it is not an empty string
func (i *idoink) AddHandler(prefix string, h H) (int, error) {
	i.m.Lock()
	defer i.m.Unlock()

	// note: prefix duplicates are allowed currently.
	// add a filter here if they should be unique
	i.hc++

	i.handlers[i.hc] = &hm{
		id:     i.hc,
		h:      h,
		prefix: prefix,
	}

	return 0, nil
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
