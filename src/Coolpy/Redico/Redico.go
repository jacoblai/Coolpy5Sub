// There are no dependencies on system binaries, and every server you start
// will be empty.
//
// Start a server with `s, err := miniredis.Run()`.
// Stop it with `defer s.Close()`.
//
// Point your Redis client to `s.Addr()` or `s.Host(), s.Port()`.
//
// Set keys directly via s.Set(...) and similar commands, or use a Redis client.
//
// For direct use you can select a Redis database with either `s.Select(12);
// s.Get("foo")` or `s.DB(12).Get("foo")`.
//
package Redico

import (
	"sync"
	"net"
	"github.com/bsm/redeo"
	"fmt"
)

type Redico struct {
	sync.Mutex
	srv        *redeo.Server
	listenAddr string
	password   string
	closed     chan struct{}
	listen     net.Listener
	dbs        map[int]*RedicoDB
	selectedDB int // DB id used in the direct Get(), Set() &c.
}

type txCmd func(*redeo.Responder, *connCtx)

// connCtx has all state for a single connection.
type connCtx struct {
	selectedDB       int            // selected DB
	authenticated    bool           // auth enabled and a valid AUTH seen
	transaction      []txCmd        // transaction callbacks. Or nil.
	dirtyTransaction bool           // any error during QUEUEing.
	watch            map[int]uint // WATCHed keys.
}

func NewRedico() *Redico {
	return &Redico{
		closed: make(chan struct{}),
		dbs:    map[int]*RedicoDB{},
	}
}

func Run() (*Redico, error) {
	m := NewRedico()
	return m, m.Start()
}

// Start starts a server. It listens on a random port on localhost. See also
// Addr().
func (m *Redico) Start() error {
	m.Lock()
	defer m.Unlock()

	l, err := listen("127.0.0.1:0")
	if err != nil {
		return err
	}
	m.listen = l
	m.listenAddr = l.Addr().String()
	m.srv = redeo.NewServer(&redeo.Config{Addr: m.listenAddr})

	commandsConnection(m, m.srv)
	commandsGeneric(m, m.srv)
	commandsString(m, m.srv)

	go func() {
		m.srv.Serve(m.listen)
		m.closed <- struct{}{}
	}()
	return nil
}

// Restart restarts a Close()d server on the same port. Values will be
// preserved.
func (m *Redico) Restart() error {
	m.Lock()
	defer m.Unlock()

	l, err := listen(m.listenAddr)
	if err != nil {
		return err
	}
	m.listen = l

	go func() {
		m.srv.Serve(m.listen)
		m.closed <- struct{}{}
	}()

	return nil
}

func listen(addr string) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		if l, err = net.Listen("tcp6", addr); err != nil {
			return nil, fmt.Errorf("failed to listen on a port: %v", err)
		}
	}
	return l, nil
}

// Close shuts down a Miniredis.
func (m *Redico) Close() {
	m.Lock()
	defer m.Unlock()
	if m.listen == nil {
		return
	}
	if m.listen.Close() != nil {
		return
	}
	m.srv.Close()
	<-m.closed
	m.listen = nil
}

// RequireAuth makes every connection need to AUTH first. Disable again by
// setting an empty string.
func (m *Redico) RequireAuth(pw string) {
	m.Lock()
	defer m.Unlock()
	m.password = pw
}

// DB returns a DB by ID.
func (m *Redico) DB(i int) *RedicoDB {
	m.Lock()
	defer m.Unlock()
	return m.db(i)
}

// get DB. No locks!
func (m *Redico) db(i int) *RedicoDB {
	if db, ok := m.dbs[i]; ok {
		return db
	}
	db := newRedicoDB(i, &m.Mutex) // the DB has our lock.
	m.dbs[i] = &db
	return &db
}

// Addr returns '127.0.0.1:12345'. Can be given to a Dial(). See also Host()
// and Port(), which return the same things.
func (m *Redico) Addr() string {
	m.Lock()
	defer m.Unlock()
	return m.listenAddr
}

// Host returns the host part of Addr().
func (m *Redico) Host() string {
	m.Lock()
	defer m.Unlock()
	host, _, _ := net.SplitHostPort(m.listenAddr)
	return host
}

// Port returns the (random) port part of Addr().
func (m *Redico) Port() string {
	m.Lock()
	defer m.Unlock()
	_, port, _ := net.SplitHostPort(m.listenAddr)
	return port
}

// CommandCount returns the number of processed commands.
func (m *Redico) CommandCount() int {
	m.Lock()
	defer m.Unlock()
	return int(m.srv.Info().TotalCommands())
}

// CurrentConnectionCount returns the number of currently connected clients.
func (m *Redico) CurrentConnectionCount() int {
	m.Lock()
	defer m.Unlock()
	return m.srv.Info().ClientsLen()
}

// TotalConnectionCount returns the number of client connections since server start.
func (m *Redico) TotalConnectionCount() int {
	m.Lock()
	defer m.Unlock()
	return int(m.srv.Info().TotalConnections())
}

// handleAuth returns false if connection has no access. It sends the reply.
func (m *Redico) handleAuth(cl *redeo.Client, out *redeo.Responder) bool {
	m.Lock()
	defer m.Unlock()
	if m.password == "" {
		return true
	}
	if cl.Ctx == nil || !getCtx(cl).authenticated {
		out.WriteErrorString("NOAUTH Authentication required.")
		return false
	}
	return true
}

func getCtx(cl *redeo.Client) *connCtx {
	if cl.Ctx == nil {
		cl.Ctx = &connCtx{}
	}
	return cl.Ctx.(*connCtx)
}

func inTx(ctx *connCtx) bool {
	return ctx.transaction != nil
}

func addTxCmd(ctx *connCtx, cb txCmd) {
	ctx.transaction = append(ctx.transaction, cb)
}

// setDirty can be called even when not in an tx. Is an no-op then.
func setDirty(cl *redeo.Client) {
	if cl.Ctx == nil {
		// No transaction. Not relevant.
		return
	}
	getCtx(cl).dirtyTransaction = true
}

func setAuthenticated(cl *redeo.Client) {
	getCtx(cl).authenticated = true
}