package admin

import (
	"context"
	"fmt"
	"github.com/elnoro/foxylock/m/v2/db"
	"github.com/tidwall/redcon"
	"log"
	"strings"
	"sync"
)

type redisProtoServer struct {
	db   db.HostDb
	addr string
	auth *auth
}

func NewRedisLikeServer(db db.HostDb, addr string, password string) DbServer {
	auth := &auth{authList: make(map[redcon.Conn]bool), password: password}

	return &redisProtoServer{
		db:   db,
		addr: addr,
		auth: auth,
	}
}

type auth struct {
	lock     sync.RWMutex
	password string
	authList map[redcon.Conn]bool
}

func (a *auth) isAuthenticated(c redcon.Conn) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.authList[c]
}

func (a *auth) authenticate(c redcon.Conn, userPass string) bool {
	// TODO use constant time compare
	if userPass == a.password {
		a.lock.Lock()
		a.authList[c] = true
		a.lock.Unlock()

		return true
	}

	return false
}

func (a *auth) delete(c redcon.Conn) {
	a.lock.Lock()
	defer a.lock.Unlock()
	delete(a.authList, c)
}

func (r *redisProtoServer) Run(_ context.Context) error {
	// TODO use context
	return redcon.ListenAndServe(r.addr, func(conn redcon.Conn, cmd redcon.Command) {
		switch strings.ToLower(string(cmd.Args[0])) {
		default:
			conn.WriteError("ERR server supports only a subset of redis commands: ADD, PING, GET, DEL")
		case "ping":
			conn.WriteString("PONG")
		case "quit":
			conn.WriteString("OK")
			err := conn.Close()
			if err != nil {
				log.Print(err)
			}
		case "auth":
			if len(cmd.Args) != 2 {
				conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
				return
			}

			pass := string(cmd.Args[1])
			result := r.auth.authenticate(conn, pass)
			if result {
				conn.WriteString("OK")
			} else {
				conn.WriteString("ERR wrong password")
			}
		case "add":
			if !r.auth.isAuthenticated(conn) {
				conn.WriteError("ERR please enter password")
				return
			}
			if len(cmd.Args) != 2 {
				conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
				return
			}

			host := string(cmd.Args[1])
			err := r.db.Save(host)
			if err != nil {
				conn.WriteString(fmt.Sprintf("ERR %v", err))
			} else {
				conn.WriteString("OK")
			}

		case "get":
			if !r.auth.isAuthenticated(conn) {
				conn.WriteError("ERR please enter password")
				return
			}
			if len(cmd.Args) != 2 {
				conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
				return
			}
			host := string(cmd.Args[1])
			ok, err := r.db.Get(host)
			if err != nil {
				conn.WriteString(fmt.Sprintf("ERR %v", err))
			} else if ok {
				conn.WriteString("Host is blocked")
			} else {
				conn.WriteString("Host is not blocked")
			}
		case "del":
			if !r.auth.isAuthenticated(conn) {
				conn.WriteError("ERR please enter password")
				return
			}
			if len(cmd.Args) != 2 {
				conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
				return
			}
			host := string(cmd.Args[1])
			err := r.db.Delete(host)
			if err != nil {
				conn.WriteString(fmt.Sprintf("ERR %v", err))
			} else {
				conn.WriteString("OK")
			}
		}
	},
		func(conn redcon.Conn) bool {
			log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			r.auth.delete(conn)
			log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
}
