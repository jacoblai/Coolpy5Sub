package Redico

import (
	"github.com/bsm/redeo"
	"strings"
	"strconv"
)

func commandsString(m *Redico, srv *redeo.Server) {
	srv.HandleFunc("GET", m.cmdGet)
	srv.HandleFunc("SET", m.cmdSet)
	srv.HandleFunc("INCR", m.cmdIncr)
}

// SET
func (m *Redico) cmdSet(out *redeo.Responder, r *redeo.Request) error {
	if len(r.Args) < 2 {
		setDirty(r.Client())
		return r.WrongNumberOfArgs()
	}
	if !m.handleAuth(r.Client(), out) {
		return nil
	}

	var (
		nx     = false // set iff not exists
		xx     = false // set iff exists
	)

	key := r.Args[0]
	value := r.Args[1]
	r.Args = r.Args[2:]
	for len(r.Args) > 0 {
		switch strings.ToUpper(r.Args[0]) {
		case "NX":
			nx = true
			r.Args = r.Args[1:]
			continue
		case "XX":
			xx = true
			r.Args = r.Args[1:]
			continue
		case "EX", "PX":
			if len(r.Args) < 2 {
				setDirty(r.Client())
				out.WriteErrorString(msgInvalidInt)
				return nil
			}
			var err error
			_, err = strconv.Atoi(r.Args[1])
			if err != nil {
				setDirty(r.Client())
				out.WriteErrorString(msgInvalidInt)
				return nil
			}
			r.Args = r.Args[2:]
			continue
		default:
			setDirty(r.Client())
			out.WriteErrorString(msgSyntaxError)
			return nil
		}
	}

	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)

		if nx {
			if db.exists(key) {
				out.WriteNil()
				return
			}
		}
		if xx {
			if !db.exists(key) {
				out.WriteNil()
				return
			}
		}

		db.del(key, true) // be sure to remove existing values of other type keys.
		// a vanilla SET clears the expire
		db.stringSet(key, value)
		//if expire != 0 {
		//	db.expire[key] = expire
		//}
		out.WriteOK()
	})
}

// GET
func (m *Redico) cmdGet(out *redeo.Responder, r *redeo.Request) error {
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
		if !db.exists(key) {
			out.WriteNil()
			return
		}
		out.WriteString(db.stringGet(key))
	})
}

// INCR
func (m *Redico) cmdIncr(out *redeo.Responder, r *redeo.Request) error {
	if len(r.Args) != 1 {
		setDirty(r.Client())
		return r.WrongNumberOfArgs()
	}
	if !m.handleAuth(r.Client(), out) {
		return nil
	}

	return withTx(m, out, r, func(out *redeo.Responder, ctx *connCtx) {
		db := m.db(ctx.selectedDB)

		key := r.Args[0]
		if !db.exists(key) {
			out.WriteErrorString(msgWrongType)
			return
		}
		v, err := db.stringIncr(key, +1)
		if err != nil {
			out.WriteErrorString(err.Error())
			return
		}
		// Don't touch TTL
		out.WriteInt(v)
	})
}
