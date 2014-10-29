package gnet

import (
	"io"
	"github.com/rare/gnet/gnproto"
)

type Response struct {
	client		*Client
	header		*gnproto.Header
	body		io.Reader
	closeflag	bool
}

func NewResponse(client *Client, header *gnproto.Header) *Response {
	return &Response {
		client: client,
		header: header,
		body: nil,
		closeflag: false,
	}
}

func (this *Response) Client() *Client {
	return this.client
}

func (this *Response) SetBodyLen(l uint32) {
	this.header.Len = l
}

func (this *Response) SetBody(rd io.Reader) {
	this.body = rd
}

func (this *Response) Flush() {
	if this.body != nil {
		b, _ := this.header.Serialize()
		this.client.Write(b)
		this.client.ReadFrom(this.body)
	}

	if this.closeflag {
		this.client.Close()
	}
}

func (this *Response) SetCloseAfterSending() {
	this.closeflag = true
}
