package admin

import "idoink"

const AdminCommand = "admin"

func Admin(e *idoink.E) (bool, error) {
	// no arguments, we dont do anything
	if len(e.Rest) < 1 {
		e.I.Message(e.From, "Admin requires a command")
		return false, nil
	}

	rest := e.Rest
	// run admin'ish level commands that may
	// adversely affect the bot
	// so it needs an AUTH SYSTEM
	subcmd := rest[0]

	switch subcmd {
	case "j":
		if len(rest) >= 2 {
			c := rest[1]
			e.I.JoinChan(c)
		}
		break
	case "p":
		if len(rest) >= 2 {
			c := rest[1]
			e.I.PartChan(c)
		}
		break
	case "d":
		// debug dump
		break
	case "q":
		e.I.Stop()
		break
	case "v":
		// version info
		break
	case "i":
		// geninfo maybe merge into version
		break
	case "s":
		// system info, current gc and heap info
		break
	}

	return false, nil
}
