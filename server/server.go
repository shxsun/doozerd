package server

import (
	"github.com/shxsun/doozerd/consensus"
	"github.com/shxsun/doozerd/store"
	"log"
	"net"
	"syscall"
)

// ListenAndServe listens on l, accepts network connections, and
// handles requests according to the doozer protocol.
func ListenAndServe(l net.Listener, canWrite chan bool, st *store.Store, p consensus.Proposer, rwsk, rosk string, self string) {
	var w bool
	for {
		c, err := l.Accept()
		if err != nil {
			if err == syscall.EINVAL {
				break
			}
			if e, ok := err.(*net.OpError); ok && !e.Temporary() {
				break
			}
			log.Println(err)
			continue
		}

		// has this server become writable?
		select {
		case w = <-canWrite:
			canWrite = nil
		default:
		}

		go serve(c, st, p, w, rwsk, rosk, self)
	}
}

func serve(nc net.Conn, st *store.Store, p consensus.Proposer, w bool, rwsk, rosk string, self string) {
	c := &conn{
		c:        nc,
		addr:     nc.RemoteAddr().String(),
		st:       st,
		p:        p,
		canWrite: w,
		rwsk:     rwsk,
		rosk:     rosk,
		self:     self,
	}

	c.grant("") // start as if the client supplied a blank password
	c.serve()
	nc.Close()
}
