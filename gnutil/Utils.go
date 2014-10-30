package gnutil

import (
	"io"
	"net"
	"os"
)

func ReadFull(conn *net.TCPConn, buf []byte) int {
	n, err := io.ReadFull(conn, buf)
	if err != nil {
		if e, ok := err.(*net.OpError); ok && e.Timeout() {
			//TODO
			//trace
			return n
		}
		//TODO
		//trace
		return -1
	}
	return n
}

func DirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func FileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil &&  fi.Mode().IsRegular()
}
