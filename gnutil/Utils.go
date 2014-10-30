package gnutil

import (
	"io"
	"net"
	"os"
	log "github.com/cihub/seelog"
)

func ReadFull(conn *net.TCPConn, buf []byte) int {
	n, err := io.ReadFull(conn, buf)
	if err != nil {
		if e, ok := err.(*net.OpError); ok && e.Timeout() {
			log.Warnf("conn(%s) read timeout", conn.RemoteAddr())
			return n
		}
		log.Warnf("conn(%s) read error: (%v)", conn.RemoteAddr(), err)
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
