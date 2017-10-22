package lastseen

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"idoink"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const LastSeenCmd = "seen"

// these are bucketed by "lastseen" bucket name
// then the kvp is nick -> last seen info
type ls struct {
	When time.Time
	Chan string
	Msg  string
}

var db *bolt.DB
var opened bool

func init() {
	var err error
	db, err = bolt.Open("lastseen.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err == nil {
		opened = true
	}
}

func UpdateLastSeen(e *idoink.E) (bool, error) {
	if !opened {
		return false, nil
	}

	// HACK:
	msg := e.Cmd + " " + strings.Join(e.Rest, "")

	s := &ls{
		When: time.Now(),
		Chan: e.To,
		Msg:  msg,
	}

	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s)
	if err != nil {
		return false, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("ls"))
		if err != nil {
			return err
		}
		return b.Put([]byte(e.From), buf.Bytes())
	})

	return false, err
}

func QueryLastSeen(e *idoink.E) (bool, error) {
	if len(e.Rest) < 1 {
		return false, nil
	}

	n := e.Rest[0]
	// lastSeenMx.Lock()
	// ls, ok := lastSeenState[n]
	// lastSeenMx.Unlock()

	vb := []byte{}
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("ls"))
		vb = b.Get([]byte(n))

		return nil
	})

	if err != nil {
		return false, err
	}

	if vb == nil {
		// not yet seen
		e.I.Message(e.To, fmt.Sprintf("%s: lastseen - no info for %s", e.From, n))
		return false, nil
	}

	br := bytes.NewReader(vb)
	enc := gob.NewDecoder(br)
	s := &ls{}
	err = enc.Decode(s)
	if err != nil {
		return false, err
	}

	e.I.Message(e.To,
		fmt.Sprintf("%s: lastseen - %s last seen on %s in %s saying %s",
			e.From, n, s.When, s.Chan, s.Msg))

	return false, nil
}
