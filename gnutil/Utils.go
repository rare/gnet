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
			return n
		}
		return -1
	}
	return n
}

func DirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.IsDir(), nil
	}
	return false, err
}

func FileExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.Mode().IsRegular(), nil
	}
	return false, err
}
