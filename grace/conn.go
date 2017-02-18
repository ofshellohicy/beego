package grace

import (
	"errors"
	"net"
	"sync"
)

type graceConn struct {
	net.Conn
	server *Server
	m      sync.Mutex
	closed bool
}

func (c *graceConn) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	// diff from here
	c.m.Lock()
	if c.closed {
		c.m.Unlock()
		return
	}
	c.server.wg.Done()
	c.closed = true
	c.m.Unlock()
	return c.Conn.Close()

	// the following is modify by ofshellohicy
	// err = c.Conn.Close()
	// if err == nil {
	// 	c.server.wg.Done()
	// } else {
	// 	log.Panicln("graceConn error:", err)
	// }
	// return err
}
