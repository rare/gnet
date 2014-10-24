package gnet

import (
	"errors"
	"net"
	"time"
	"github.com/rare/gnet/gnproto"
)

type Server struct {
	exit		chan bool				
	ln			*net.TCPListener		
	cn			uint32							//number of connected clients
	hm			map[uint16]HandlerFuncType		//request handler map
}

func (this *Server) handleConnection(conn *net.TCPConn) {
	//logger.Printf("accept connection from (%s) (%p)", conn.RemoteAddr(), conn)
	
	if this.cn >= Conf.MaxClients {
		//logger.Printf("too many clients")
		//TODO
		conn.Close()
		return
	}

	this.cn++
	defer func(){
		conn.Close()
		this.cn--
	}()

	cli := NewClient()
	err := cli.Init(conn, this)
	if err != nil {
		//logger.Printf("init client error")
		//TODO
		return
	}

	cli.Process()
}

func NewServer() *Server {
	return &Server{
		exit:		make(chan bool),
		ln:			nil,	
		cn:			0,
		hm:			make(map[uint16]HandlerFuncType),
	}
}

func (this *Server) Init(conf *Config) error {
	if conf == nil {
		return errors.New("nil parameter")
	}

	Conf = *conf

	tcpAddr, err := net.ResolveTCPAddr("tcp4", Conf.Addr)
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		//logger.Fatalf("failed to listen, (%v)", err)
		//TODO
		return err
	}
	this.ln = ln

	return nil
}

func (this *Server) HandleFunc(cmd uint16, handler HandlerFuncType) error {
	if cmd == gnproto.CMD_HEART_BEAT || handler == nil {
		return errors.New("cmd is 0(heartbeat cmd) or handler is nil")
	}

	this.hm[cmd] = handler
	return nil
}

func (this *Server) Dispatch(client *Client, header *gnproto.Header, body []byte) error {
	if handler, ok := this.hm[header.Cmd]; ok {
		return handler(client, header, body)
	}
	return errors.New("command handler not found")
}

func (this *Server) Run() {
	defer func() {
		this.ln.Close()
		this.ln = nil
		this.cn = 0
	}()

	for {
		select {
		case <-this.exit:
			//logger.Printf("ask me to quit")
			//TODO
			return
		default:
		}

		this.ln.SetDeadline(time.Now().Add(2 * time.Second))
		conn, err := this.ln.AcceptTCP()
		if err != nil {
			if e, ok := err.(*net.OpError); ok && e.Timeout() {
				continue
			}
			//logger.Printf("accept failed: %v\n", err)
			//TODO
			continue
		}

		go this.handleConnection(conn)
	}
}

func (this *Server) Stop() {
	close(this.exit)
}
