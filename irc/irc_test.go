package irc_test

import (
	"idoink/irc"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	irc, err := irc.New("idoink"+strconv.Itoa(time.Now().Minute()), "irc.quakenet.org:6667")
	assert.NoError(t, err)
	assert.NotNil(t, irc)
	//assert.NoError(t, irc.Close())
}
