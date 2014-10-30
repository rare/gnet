package gnet

import (
	"fmt"		//debug

	"errors"
	"net"
	"io"
	"time"
	"github.com/rare/gnet/gnproto"
	"github.com/rare/gnet/gnutil"
)

type Client struct {
	exit		chan bool
	conn		*net.TCPConn
	server		*Server
	hbtime		time.Time
	outchan		chan *Response
	storage		*gnutil.Storage
}

func (this *Client) handleInput() {
	for {
		var (
			now = time.Now()
			headbuf = make([]byte, gnproto.HEADER_SIZE)
			header		gnproto.Header
		)

		if now.After(this.hbtime.Add(time.Duration(Conf.HBTimeout) * time.Second)) {
			//debug
			fmt.Println("conn heartbeat timeout")
			//TODO
			break
		}

		this.conn.SetReadDeadline(now.Add(time.Duration(Conf.ReadTimeout) * time.Second))
		if len(headbuf) != gnutil.ReadFull(this.conn, headbuf) {
			//debug
			fmt.Println("conn read head timeout")
			//TODO
			break
		}

		if err := header.Deserialize(headbuf); err != nil {
			//debug
			fmt.Println("conn parse head error")
			//TODO
			break
		}

		if Conf.MaxBodyLen != 0 && header.Len > Conf.MaxBodyLen {
			//debug
			fmt.Printf("head len: %d, max len: %d\n", header.Len, Conf.MaxBodyLen)
			fmt.Println("conn body len too large")
			//TODO
			break
		}

		if header.Cmd == gnproto.CMD_HEART_BEAT {
			this.hbtime = time.Now()
		} else {
			req := NewRequest(this, &header)
			resp := NewResponse(this, &header)
			if err := this.server.Dispatch(req, resp); err != nil {
				//debug
				fmt.Println("dispatch error")
				//TODO
				break
			}
			this.outchan <- resp
		}
	}

	this.exit <- true
}

func (this *Client) handleOutput() {
	for {
		select {
			case resp, ok := <-this.outchan:
				if !ok {
					//TODO
					return
				}
				resp.Flush()

			case <-this.exit:
				//TODO
				//trace
				return
		}
	}
}

func NewClient() *Client {
	return &Client{
		exit:		make(chan bool),
		conn:		nil,
		server:		nil,
		hbtime:		time.Now(),
		outchan:	make(chan *Response, Conf.OutChanBufSize),
		storage:	gnutil.NewStorage(),
	}
}

func (this *Client) Init(conn *net.TCPConn, server *Server) error {
	if conn == nil || server == nil {
		//TODO
		//trace
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

func (this *Client) Write(buf []byte) int {
	this.conn.SetWriteDeadline(time.Now().Add(time.Duration(Conf.WriteTimeout) + time.Second))
	n, err := this.conn.Write(buf)
	if n != len(buf) || err != nil {
		//TODO
	}
	return n
}

func (this *Client) ReadFrom(rd io.Reader) int {
	var nn = 0
	buf := make([]byte, 32*1024)
	for {
		n, err:= rd.Read(buf)
		if n > 0 {
			this.conn.SetWriteDeadline(time.Now().Add(time.Duration(Conf.WriteTimeout) + time.Second))
			n, err := this.conn.Write(buf)
			if n != len(buf) || err !=nil {
				//TODO
				break
			}
			nn += n
		}
		if err == io.EOF {
			//TODO
			//trace
			break
		}
		if err != nil {
			//TODO
			break
		}
	}

	return nn
}

func (this *Client) Close() {
	this.conn.Close()
}

func (this *Client) Storage() *gnutil.Storage {
	return this.storage
}
