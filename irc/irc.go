package irc

import "github.com/labstack/gommon/log"

type IRC struct {
	c     *conn
	Nick  string
	Host  string
	Chans []string
}

func (i *IRC) cmd(f string, args ...interface{}) {
	//log.Printf("sending -> %s %+v", f, args...)
	i.c.write(f, args...)
}

func New(nick, host string, chans []string) (*IRC, error) {
	r := &IRC{
		Nick:  nick,
		Host:  host,
		Chans: chans,
	}
	return r, nil
}

func (i *IRC) Cmd(c string) {
	i.cmd(c)
}

func (i *IRC) Pong(r string) {
	i.cmd("PONG %s", r)
}

func (i *IRC) Register() {
	i.User(i.Nick, "8", "*", i.Nick)
	i.SetNick(i.Nick)
}

type DoneCallback func(error)

type MessageCallback func(string)

func (i *IRC) Start(m MessageCallback, d DoneCallback) error {
	conn, err := connect(i.Host)
	if err != nil {
		d(err)
		return err
	}
	i.c = conn

	go i.readLoop(m, d)

	return nil
}

func (i *IRC) Close() error {
	return i.c.disconnect()
}

func (i *IRC) User(user, host, server, realname string) {
	i.cmd("USER %s %s %s :%s", user, host, server, realname)
}

func (i *IRC) SetNick(nn string) {
	i.cmd("NICK %s", nn)
}

func (i *IRC) parse(rawstr string) {
}

func (i *IRC) readLoop(m MessageCallback, done DoneCallback) {
	for {
		r, err := i.c.read()
		if err != nil {
			log.Error(err)
			if done != nil {
				done(err)
			}
			return
		}
		m(r)
	}
}
