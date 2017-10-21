package idoink

import (
	"idoink/irc"
)

// note, pulling handler args to seperate struct makes it easier to change

// E is an irc event (currently only private messages) passed to handlers.
// it has a reference to the IRC client so it can send message responses, etc.
type E struct {
	From string
	To   string
	Rest []string
	IRC  *irc.IRC
}

// H is the main hook for handler packages to do something interesting
// with a chat line.
// It returns true if it stops any further processing, and any error that occured.
type H func(e *E) (bool, error)

// internal bookkeeping for registered prefix
type hm struct {
	id     int
	prefix string
	h      H
}
