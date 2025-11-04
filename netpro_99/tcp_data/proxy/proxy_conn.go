package proxy

import (
	"io"
	"net"
)

func proxyConn(src, dst string) error {
	connSrc, err := net.Dial("tcp", src)
	if err != nil {
		return err
	}

	defer connSrc.Close()

	connDst, err := net.Dial("tcp", dst)
	if err != nil {
		return err
	}
	defer connDst.Close()

	go func() {
		_, _ = io.Copy(connSrc, connDst)
	}()

	_, err = io.Copy(connDst, connSrc)

	return err
}
