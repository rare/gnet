package gnet

import (
	gnproto "./gnproto"
)

type HandlerFuncType func(client *Client, header *gnproto.Header, body []byte) error

type Header interface {
	HandleFunc(cmd uint16, handler HandlerFuncType) error
	Dispatch(client *Client, header *gnproto.Header, body []byte) error
}


