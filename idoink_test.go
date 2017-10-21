package idoink_test

import (
	"idoink"
	"testing"

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
