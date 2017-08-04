package irc

import (
	"bufio"
	"net"
	"net/textproto"
)

type conn struct {
	tcp net.Conn
	r   *textproto.Reader
	w   *textproto.Writer
}

func connect(host string) (*conn, error) {
	nconn, err := net.Dial("tcp", host) //+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	rb := bufio.NewReader(nconn)
	wb := bufio.NewWriter(nconn)
	return &conn{
		tcp: nconn,
		r:   textproto.NewReader(rb),
		w:   textproto.NewWriter(wb),
	}, nil
}

func (c *conn) write(fs string, args ...interface{}) error {
	return c.w.PrintfLine(fs, args...)
}

func (c *conn) read() (string, error) {
	return c.r.ReadLine()
}
