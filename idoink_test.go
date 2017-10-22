package idoink_test

import (
	"idoink"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var nick = "testnick"
var server = "us.quakenet.org:6667"
var chans = "#idoinkbottest"

func TestNew(t *testing.T) {

	i := idoink.New(nick, server, chans)
	assert.NotNil(t, i)
	assert.Equal(t, nick, i.Nick())
	assert.Equal(t, []string{chans}, i.Chans())
}

func TestAddRemoveHandler(t *testing.T) {
	i := idoink.New(nick, server, chans)
	assert.NotNil(t, i)

	hid, err := i.AddHandler("foo", func(e *idoink.E) (bool, error) {
		return false, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, hid)

	err = i.RemoveHandler(hid)
	assert.NoError(t, err)
}

func TestStartStop(t *testing.T) {
	i := idoink.New(nick, server, chans)
	assert.NotNil(t, i)

	err := i.Start()
	assert.NoError(t, err)

	// give it a few seconds then slam it off
	time.Sleep(time.Second * 2)

	err = i.Stop()
	// dont assert on this error, sometimes we get errors from runnign test too fast
	//assert.NoError(t, err)

	// this sucks, even if we dont care about ret, if errors come in the test fails
	// if they come in after some how..
	// 2017/10/21 18:01:34 read tcp 192.168.1.5:13912->170.178.184.36:6667: use of closed network connection
	// 2017/10/21 18:01:34 read tcp 192.168.1.5:13912->170.178.184.36:6667: use of closed network connection
}
