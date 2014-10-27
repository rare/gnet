package gnet

import (
	"errors"
	"github.com/rare/gnet/gnproto"
	"github.com/rare/gnet/gnutil"
)

type Request struct {
	client	*Client
	header	*gnproto.Header
}

func NewRequest(client *Client, header *gnproto.Header) *Request {
	return &Request {
		client: client,
		header: header,
	}
}

func (this *Request) Cmd() uint16 {
	return this.header.Cmd
}

func (this *Request) Client() *Client {
	return this.client
}

func (this *Request) Body() ([]byte, error) {
	if this.header.Len > 0 {
		bodybuf := make([]byte, this.header.Len)
		if len(bodybuf) != gnutil.ReadFull(this.client.conn, bodybuf) {
			return bodybuf, errors.New("read request body error")
		}
		return bodybuf, nil
	}

	return nil, nil
}

func (this *Request) BodyLen() uint32 {
	return this.header.Len
}
