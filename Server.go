package gnet

import (
	"container/list"
	"errors"
	"net"
	"time"
	"github.com/rare/gnet/gnproto"
	"github.com/rare/gnet/gnfilter"
	"github.com/rare/gnet/gnutil"
	log "github.com/cihub/seelog"
)

type HandlerFuncType func(req *Request, resp *Response) error

type Server struct {
	exit		chan bool
	ln			*net.TCPListener
	hm			map[uint16]HandlerFuncType		//request handler map
	filters		*list.List						//filters 
	storage		*gnutil.Storage
}

func (this *Server) addSysFilters() {
	max_conn_filter := gnfilter.NewMaxConnFilter(Conf.MaxClients)
	this.FilterFunc(*max_conn_filter)
	bwl_filter := gnfilter.NewBlackWhiteListFilter()
	err := bwl_filter.Init(Conf.BlackListFile, Conf.WhiteListFile)
	if err != nil {
		log.Warnf("init black white list filter error: (%v)", err)
		return
	}
	this.FilterFunc(*bwl_filter)

	//More sys filters
	//TODO
}

func (this *Server) doFilters(evt gnfilter.EventType, obj interface{}) gnfilter.FilterResult {
	for f := this.filters.Front(); f != nil; f = f.Next() {
		filter, _ := f.Value.(gnfilter.Filter)

		if filter.CareEvent(evt) {
			fr := filter.DoFilter(gnfilter.EVT_CONN_ACCEPTED, obj)

			if fr == gnfilter.FR_ABORT {
				return gnfilter.FR_ABORT
			}

			if fr == gnfilter.FR_END {
				break
			}
		}
	}

	return gnfilter.FR_OK
}

func (this *Server) handleConnection(conn *net.TCPConn) {
	log.Tracef("accept connection from (%s) (%p)", conn.RemoteAddr(), conn)

	fr := this.doFilters(gnfilter.EVT_CONN_ACCEPTED, conn)
	if fr != gnfilter.FR_OK {
		log.Warnf("conn(%s) aborted by filter", conn.RemoteAddr())
		conn.Close()
		return
	}

	defer func(){
		this.doFilters(gnfilter.EVT_CONN_CLOSING, conn)
		conn.Close()
		log.Tracef("conn(%s) closed", conn.RemoteAddr())
	}()

	cli := NewClient()
	err := cli.Init(conn, this)
	if err != nil {
		log.Warnf("conn(%s) init client error: (%v)", conn.RemoteAddr(), err)
		return
	}

	cli.Process()
}

func NewServer() *Server {
	return &Server{
		exit:		make(chan bool),
		ln:			nil,
		hm:			make(map[uint16]HandlerFuncType),
		filters:	list.New(),
		storage:	gnutil.NewStorage(),
	}
}

func (this *Server) Init(conf *Config) error {
	if conf == nil {
		log.Critical("init server error, nil conf")
		return errors.New("nil parameter")
	}

	Conf = *conf

	tcpAddr, err := net.ResolveTCPAddr("tcp4", Conf.Addr)
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Critical("failed to listen on %s, (%v)", Conf.Addr, err)
		return err
	}
	this.ln = ln

	this.addSysFilters()

	return nil
}

func (this *Server) FilterFunc(filter gnfilter.Filter) {
	this.filters.PushBack(filter)
}

func (this *Server) HandleFunc(cmd uint16, handler HandlerFuncType) error {
	if cmd == gnproto.CMD_HEART_BEAT || handler == nil {
		log.Warnf("set handler error, bad cmd(%d) or handler is nil", cmd)
		return errors.New("cmd is 0(heartbeat cmd) or handler is nil")
	}

	this.hm[cmd] = handler
	return nil
}

func (this *Server) Dispatch(req *Request, resp *Response) error {
	if handler, ok := this.hm[req.Cmd()]; ok {
		return handler(req, resp)
	}
	return errors.New("command handler not found")
}

func (this *Server) Run() {
	defer func() {
		log.Tracef("server(%v) exiting", this.ln.Addr())

		this.ln.Close()
		this.ln = nil
	}()

	for {
		select {
		case <-this.exit:
			log.Tracef("someone ask me to quit")
			return
		default:
		}

		this.ln.SetDeadline(time.Now().Add(2 * time.Second))
		conn, err := this.ln.AcceptTCP()
		if err != nil {
			if e, ok := err.(*net.OpError); ok && e.Timeout() {
				log.Tracef("server(%v) accept timeout", this.ln.Addr())
				continue
			}
			log.Warnf("server(%v) accept error: (%v)", this.ln.Addr(), err)
			continue
		}

		go this.handleConnection(conn)
	}
}

func (this *Server) Stop() {
	close(this.exit)
}

func (this *Server) Storage() *gnutil.Storage {
	return this.storage
}
