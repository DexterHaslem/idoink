package admin

import (
	"fmt"
	"idoink"
	"runtime"
)

const AdminCommand = "admin"

func Admin(e *idoink.E) (bool, error) {
	// no arguments, we dont do anything
	if len(e.Rest) < 1 {
		e.I.Message(e.To, "Admin requires a command")
		return false, nil
	}

	rest := e.Rest
	// run admin'ish level commands that may
	// adversely affect the bot
	// so it needs an AUTH SYSTEM

	// TODO: real auth
	if e.From != "dexter" {
		// netsplit steal my nick and kill my bot :((
		e.I.Message(e.To, "yer unauthed")
		return false, nil
	}
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
		e.I.Message(e.To, fmt.Sprintf("%s: version %s", e.From, "TODO"))
		break
	case "i":
		// geninfo maybe merge into version
		break
	case "s":
		sysInfo(e)
	}

	return false, nil
}

func sysInfo(e *idoink.E) {
	// system info, current gc and heap info
	s := runtime.MemStats{}
	runtime.ReadMemStats(&s)

	// raw bytes -> nice string
	np := func(v uint64) string {
		// convert to kB/mB and add type
		suffix := "KB"
		converted := float64(v) / 1024.0
		if converted > 1024 {
			converted /= 1024
			suffix = "MB"
		}
		return fmt.Sprintf("%.1f%s", converted, suffix)
	}

	msg := fmt.Sprintf("%s: meminfo (HEAP): alloc=%s idle=%s inuse=%s objects=%s released=%s sys=%s",
		e.From, np(s.HeapAlloc), np(s.HeapIdle), np(s.HeapInuse), np(s.HeapObjects),
		np(s.HeapReleased), np(s.HeapSys))

	e.I.Message(e.To, msg)
}
