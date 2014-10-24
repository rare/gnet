package gnet

import (
	"errors"
	"io"
	"net"
	"time"
	gnproto "./gnproto"
)

type Client struct {
	exit		chan bool
	conn		*net.TCPConn
	server		*Server
	hbtime		time.Time
}

func (this *Client) readFull(buf []byte) int {
	n, err := io.ReadFull(this.conn, buf)
	if err != nil {
		if e, ok := err.(*net.OpError); ok && e.Timeout() {
			return n
		}
		return -1
	}
	return n
}

func (this *Client) handleInput() {
	for {
		var (
			now = time.Now()
			headbuf = make([]byte, gnproto.HEADER_SIZE)
			bodybuf		[]byte
			header		gnproto.Header
		)

		if now.After(this.hbtime.Add(time.Duration(Conf.HBTimeout) * time.Second)) {
			//logger.Printf("%p: heartbeat timeout", this.conn)
			//TODO
			break
		}

		this.conn.SetReadDeadline(now.Add(time.Duration(Conf.ReadTimeout) * time.Second))
		if len(headbuf) != this.readFull(headbuf) {
			//logger.Printf("%p: read header timeout", this.conn)
			//TODO
			break
		}

		if err := header.Deserialize(headbuf); err != nil {
			//logger.Printf("%p: deserialize header error", this.conn)
			//TODO
			break
		}

		if header.Len > Conf.MaxBodyLen {
			//logger.Printf("%p: header len too big", this.conn)
			//TODO
			break
		}

		if header.Len > 0 {
			bodybuf = make([]byte, header.Len)
			if len(bodybuf) != this.readFull(bodybuf) {
				//logger.Printf("%p: read body timeout", this.conn)
				//TODO
				break
			}
		}

		if header.Cmd == gnproto.CMD_HEART_BEAT {
			this.hbtime = time.Now()
		} else {
			if err := this.server.Dispatch(this, &header, bodybuf); err != nil {
				//logger.Printf("%p: dispatch command error", this.conn)
				//TODO
				break
			}
		}
	}

	this.exit <- true
}

func (this *Client) handleOutput() {
	for {
		select {
			/*
			case pack := <-client.outMsgs:
				time.Sleep(1 * time.Second)
			*/
			case <-this.exit:
				return
		}
	}
}

func NewClient() *Client {
	return &Client{
		exit:		make(chan bool),
		conn:		nil,
		server:	nil,
		hbtime:		time.Now(),
	}
}

func (this *Client) Init(conn *net.TCPConn, server *Server) error {
	if conn == nil || server == nil {
		return errors.New("nil parameter")
	}

	this.conn = conn
	this.server = server
	return nil
}

func (this *Client) Process() {
	go this.handleOutput()
	this.handleInput()
}

