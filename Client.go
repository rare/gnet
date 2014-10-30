package gnet

import (
	"errors"
	"net"
	"io"
	"time"
	"github.com/rare/gnet/gnproto"
	"github.com/rare/gnet/gnutil"
	log "github.com/cihub/seelog"
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
	log.Tracef("conn(%s)'s input routine start", this.conn.RemoteAddr())

	for {
		var (
			now = time.Now()
			headbuf = make([]byte, gnproto.HEADER_SIZE)
			header		gnproto.Header
		)

		if now.After(this.hbtime.Add(time.Duration(Conf.HBTimeout) * time.Second)) {
			log.Warnf("conn(%s) heartbeat timeout", this.conn.RemoteAddr())
			break
		}

		this.conn.SetReadDeadline(now.Add(time.Duration(Conf.ReadTimeout) * time.Second))
		if len(headbuf) != gnutil.ReadFull(this.conn, headbuf) {
			log.Warnf("conn(%s) read head timeout", this.conn.RemoteAddr())
			break
		}

		if err := header.Deserialize(headbuf); err != nil {
			log.Warnf("conn(%s) parse head error", this.conn.RemoteAddr())
			break
		}

		if Conf.MaxBodyLen != 0 && header.Len > Conf.MaxBodyLen {
			log.Warnf("conn(%s) body too large, head len: %d, max len: %d\n", this.conn.RemoteAddr(), header.Len, Conf.MaxBodyLen)
			break
		}

		if header.Cmd == gnproto.CMD_HEART_BEAT {
			log.Tracef("conn(%s) received heartbeat", this.conn.RemoteAddr())
			this.hbtime = time.Now()
		} else {
			req := NewRequest(this, &header)
			resp := NewResponse(this, &header)
			if err := this.server.Dispatch(req, resp); err != nil {
				log.Warnf("conn(%s) dispatch command(%d) error", this.conn.RemoteAddr(), header.Cmd)
				break
			}
			this.outchan <- resp
		}
	}

	log.Tracef("conn(%s)'s input routine end", this.conn.RemoteAddr())
	this.exit <- true
}

func (this *Client) handleOutput() {
	log.Tracef("conn(%s)'s output routine start", this.conn.RemoteAddr())

	for {
		select {
			case resp, _ := <-this.outchan:
				resp.Flush()

			case <-this.exit:
				log.Tracef("conn(%s)'s output routine end", this.conn.RemoteAddr())
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
		log.Warnf("Client.Init error, nil parameter")
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
		log.Warnf("conn(%s) write error: (%v)", this.conn.RemoteAddr(), err)
		return -1
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
			if n != len(buf) || err != nil {
				log.Warnf("conn(%s) write error(%v) in ReadFrom", this.conn.RemoteAddr(), err)
				break
			}
			nn += n
		}
		if err == io.EOF {
			log.Trace("conn(%s) ReadFrom rd(%p) end", this.conn.RemoteAddr(), rd)
			break
		}
		if err != nil {
			log.Warnf("conn(%s) ReadFrom rd(%p) error: (%v)", this.conn.RemoteAddr(), rd, err)
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
