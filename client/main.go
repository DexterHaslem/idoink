package main

import (
	"flag"
	"idoink"
	"idoink/handlers/admin"
	"idoink/handlers/darksky"
	"idoink/handlers/ddg"
	"idoink/handlers/lastfm"
	"idoink/handlers/lastseen"
	"log"
)

func main() {
	nick := flag.String("nick", "", "IRC nickname to use")
	server := flag.String("server", "", "IRC server:port")
	chans := flag.String("chans", "", "IRC channels to join comma separated")

	flag.Parse()

	if *nick == "" || *server == "" {
		flag.Usage()
		return
	}

	i := idoink.New(*nick, *server, *chans)

	// just register all handlers for now. realistically we could register one simple one
	// that allows registering rest via admin cmds

	i.AddHandler(admin.AdminCommand, admin.Admin)
	i.AddHandler(ddg.DDGCmd, ddg.DDG)
	i.AddHandler(lastfm.LastfmCmd, lastfm.LastFM)
	i.AddHandler(lastseen.LastSeenCmd, lastseen.LastSeen)
	i.AddHandler(darksky.Cmd, darksky.DarkSky)

	if err := i.Start(); err != nil {
		log.Fatal(err)
	}
}
