## idoink

_idoink_ is a go library for writing IRC bots.


### usage

Using idoink requires just a few steps:

1. create a new instance of idoink.I using New()
1. add message handlers
1. Start!


Here is an example using an existing handler:


```go
import "github.com/DexterHaslem/idoink"

// ....

// note the server has :port
i := idoink.New("mynick", "irc.freenode.net:6667", "flatearth")

// this is a provided handler in idoink/handlers/admin
i.AddHandler(admin.AdminCommand, admin.Admin)

// start will block caller on success. if you want async, run in a goroutine
if err := i.Start(); err != nil {
	log.Fatal(err)
}
	// profit!
```

> NOTE: there is also a client using all provided handlers in `client/main.go`



### Creating handlers

to create a handler, you simply need to make a function
with the following signature:

```go
func(e *idoink.E) (bool, error)
```

e is the event, containing everything from the received message:

```go
// E is an irc event (currently only private messages) passed to handlers.
// it has a reference to the bot client so it can send message responses, etc.
type E struct {
	From string
	To   string
	Rest []string
	I    I
}

```

for the return values of the function _`(bool, error)`_,
return true if you want your function to STOP further commands
from processing. This will be unlikely for most commands.
Return any errors that occur during processing if
possible and the error will be reported.

To interact with IRC you can use the provided `I` instance:

```go
type I interface {

	// These are wrappers above irc which
	// will be exposed to the handlers

	Nick() string
	SetNick(string) error
	Chans() []string
	JoinChan(string) error
	PartChan(string) error
	Message(to string, msg string) error

	// raw irc line, left as is, but \r\n is added to end
	Raw(string) error

	// ...
}

```

**Command prefix**

When registering a handler, you can give it a prefix or empty line.

For example if you want the bot to handle the command "foo"
for your handler, you would do something like this

```go
i.AddHandler("foo", myHandlerFunc)
```

If the prefix is empty, your handler will always execute for
every private message received in a channel. This is where
returning true in your handler would be important if you wanted no further handlers to process the message.

**NOTE:** Handlers are evaluated in order they are registered!
