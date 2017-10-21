package idoink

// H is the main hook for handler packages to do something interesting
// with a chat line.
// It returns true if it stops any further processing, and any error that occured.
type H func(from, to string, rest ...string) (bool, error)

// internal bookkeeping for registered prefix
type hm struct {
	id     int
	prefix string
	h      H
}
