package websocket

import (
	"bytes"
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
	buff *bytes.Buffer
}

func newConn(srv *Server, c *websocket.Conn) *Conn {
	return &Conn{
		srv:  srv,
		conn: c,
		buff: &bytes.Buffer{},
	}
}

func (c *Conn) Read(p []byte) (int, error) {
	_, b, err := c.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	if _, err = c.buff.Write(b); err != nil {
		return 0, err
	}
	return c.buff.Read(p)
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
