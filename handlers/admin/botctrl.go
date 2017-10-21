package main

const bamCmd = "!"

func bam(from, to string, rest ...string) {
	if len(rest) < 1 {
		return
	}

	// TODO: point of the bang! (renamed bam here)
	// is to run admin'ish level commands that may
	// adversely affect the bot
	// so it needs an AUTH SYSTEM
	subcmd := rest[0]
	switch subcmd {
	case "j":
		if len(rest) >= 2 {
			c := rest[1]
			i.Join(c)
		}
		break
	case "p":
		if len(rest) >= 2 {
			c := rest[1]
			i.Part(c)
		}
		break
	case "d":
		// debug dump
		break
	case "q":
		i.Close()
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
}
