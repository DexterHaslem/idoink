package main

import (
	"flag"
	"idoink/irc"
	"log"
	"strings"
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

	parsedChans := []string{}
	if *chans != "" {
		parsedChans = strings.Split(*chans, ",")
	}

	i, err := irc.New(*nick, *server, parsedChans)
	if err != nil {
		log.Fatal(err)
	}

	mc := make(chan string, 5)

	i.Start(func(m string) {
		mc <- m
	}, func(e error) {
		log.Fatal(e)
		close(mc)
	})

	for m := range mc {
		parseMsg(m, i)
	}
}
