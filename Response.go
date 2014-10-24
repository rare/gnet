package gnet

import (
	"github.com/rare/gnet/gnproto"
)

type Response struct {
	client		*Client
	header		*gnproto.Header
}

func NewResponse(client *Client, header *gnproto.Header) *Response {
	return &Response {
		client: client,
		header: header,
	}
}

func (this *Response) Flush() {
}
