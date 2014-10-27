package gnet

import (
	"io"
	"github.com/rare/gnet/gnproto"
)

type Response struct {
	client		*Client
	header		*gnproto.Header
	body		io.Reader
}

func NewResponse(client *Client, header *gnproto.Header) *Response {
	return &Response {
		client: client,
		header: header,
		body: nil,
	}
}

func (this *Response) Client() *Client {
	return this.client
}

func (this *Response) SetBody(rd io.Reader) {
	this.body = rd
}

func (this *Response) Flush() {
	b, _ := this.header.Serialize()
	this.client.Write(b)
	this.client.ReadFrom(this.body)
}
