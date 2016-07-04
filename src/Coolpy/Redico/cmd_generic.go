package Redico

import (
	"github.com/bsm/redeo"
	"strings"
	"strconv"
)

func commandsGeneric(m *Redico, srv *redeo.Server) {
	srv.HandleFunc("DEL", m.cmdDel)
	// DUMP
	srv.HandleFunc("EXISTS", m.cmdExists)
	srv.HandleFunc("KEYS", m.cmdKeys)
	srv.HandleFunc("KEYSSTART", m.cmdKeysStart)
	srv.HandleFunc("KEYSRANGE", m.cmdKeyRange)
	// OBJECT
	srv.HandleFunc("SCAN", m.cmdScan)
}

// DEL
func (m *Redico) cmdDel(out *redeo.Responder, r *redeo.Request) error {
	if !m.handleAuth(r.Client(), out) {
		return nil
	}
	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)

		count := 0
		for _, key := range r.Args {
			if db.exists(key) {
				count++
			}
			db.del(key, true) // delete expire
		}
		out.WriteInt(count)
	})
}

// EXISTS
func (m *Redico) cmdExists(out *redeo.Responder, r *redeo.Request) error {
	if len(r.Args) < 1 {
		setDirty(r.Client())
		return r.WrongNumberOfArgs()
	}
	if !m.handleAuth(r.Client(), out) {
		return nil
	}

	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)

		found := 0
		for _, k := range r.Args {
			if db.exists(k) {
				found++
			}
		}
		out.WriteInt(found)
	})
}

// KEYS
func (m *Redico) cmdKeys(out *redeo.Responder, r *redeo.Request) error {
	if len(r.Args) != 1 {
		setDirty(r.Client())
		return r.WrongNumberOfArgs()
	}
	if !m.handleAuth(r.Client(), out) {
		return nil
	}

	key := r.Args[0]

	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)

		keys := matchKeys(db.allKeys(), key)
		out.WriteBulkLen(len(keys))
		for _, s := range keys {
			out.WriteString(s)
		}
	})
}

// KEYSSTART
func (m *Redico) cmdKeysStart(out *redeo.Responder, r *redeo.Request) error {
	if len(r.Args) != 1 {
		setDirty(r.Client())
		return r.WrongNumberOfArgs()
	}
	if !m.handleAuth(r.Client(), out) {
		return nil
	}

	key := r.Args[0]

	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)

		keys := db.keyStart(key)
		out.WriteBulkLen(len(keys))
		for _, s := range keys {
			out.WriteString(s)
		}
	})
}

// KEYSRANGE
func (m *Redico) cmdKeyRange(out *redeo.Responder, r *redeo.Request) error {
	if len(r.Args) < 2 {
		setDirty(r.Client())
		return r.WrongNumberOfArgs()
	}
	if !m.handleAuth(r.Client(), out) {
		return nil
	}

	min := r.Args[0]
	max := r.Args[1]

	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)

		keys := db.keyRange(min, max)
		out.WriteBulkLen(len(keys))
		for _, s := range keys {
			out.WriteString(s)
		}
	})
}

// SCAN
func (m *Redico) cmdScan(out *redeo.Responder, r *redeo.Request) error {
	if len(r.Args) < 1 {
		setDirty(r.Client())
		return r.WrongNumberOfArgs()
	}
	if !m.handleAuth(r.Client(), out) {
		return nil
	}

	cursor, err := strconv.Atoi(r.Args[0])
	if err != nil {
		setDirty(r.Client())
		out.WriteErrorString(msgInvalidCursor)
		return nil
	}
	// MATCH and COUNT options
	var withMatch bool
	var match string
	args := r.Args[1:]
	for len(args) > 0 {
		if strings.ToLower(args[0]) == "count" {
			if len(args) < 2 {
				setDirty(r.Client())
				out.WriteErrorString(msgSyntaxError)
				return nil
			}
			_, err := strconv.Atoi(args[1])
			if err != nil {
				setDirty(r.Client())
				out.WriteErrorString(msgInvalidInt)
				return nil
			}
			// We do nothing with count.
			args = args[2:]
			continue
		}
		if strings.ToLower(args[0]) == "match" {
			if len(args) < 2 {
				setDirty(r.Client())
				out.WriteErrorString(msgSyntaxError)
				return nil
			}
			withMatch = true
			match = args[1]
			args = args[2:]
			continue
		}
		setDirty(r.Client())
		out.WriteErrorString(msgSyntaxError)
		return nil
	}

	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)
		// We return _all_ (matched) keys every time.

		if cursor != 0 {
			// Invalid cursor.
			out.WriteBulkLen(2)
			out.WriteString("0") // no next cursor
			out.WriteBulkLen(0)  // no elements
			return
		}

		keys := db.allKeys()
		if withMatch {
			keys = matchKeys(keys, match)
		}

		out.WriteBulkLen(2)
		out.WriteString("0") // no next cursor
		out.WriteBulkLen(len(keys))
		for _, k := range keys {
			out.WriteString(k)
		}
	})
}