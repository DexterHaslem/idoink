package irc

type IRC struct {
	c    *conn
	nick string
	host string
}

func New(nick, host string, port int) (*IRC, error) {
	conn, err := connect(host)
	if err != nil {
		return nil, err
	}
	r := &IRC{
		c: conn,
	}
	return r, nil
}

func (i *IRC) User(user, host, server, realname string) {
	i.c.write("USER %s %s %s :%s", user, host, server, realname)
}
