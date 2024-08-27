package websocket

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const (
	TextMessage = websocket.TextMessage

	BinaryMessage = websocket.BinaryMessage
)

type Conn struct {
	srv  *Server
	conn *websocket.Conn
}

func newConn(srv *Server, c *websocket.Conn) *Conn {
	return &Conn{
		srv:  srv,
		conn: c,
	}
}

func (c *Conn) Read(p []byte) (int, error) {
	_, r, err := c.conn.NextReader()
	if err != nil {
		return 0, err
	}
	return r.Read(p)
}

func (c *Conn) Write(p []byte) (n int, err error) {
	err = c.conn.WriteMessage(c.srv.Option.MessageType, p)
	n = len(p)
	return
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.NetConn().SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
