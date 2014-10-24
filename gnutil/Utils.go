package gnutil

import (
	"io"
	"net"
)

func ReadFull(conn *net.TCPConn, buf []byte) int {
	n, err := io.ReadFull(conn, buf)
	if err != nil {
		if e, ok := err.(*net.OpError); ok && e.Timeout() {
			return n
		}
		return -1
	}
	return n
}

