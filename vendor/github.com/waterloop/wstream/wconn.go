package wstream

import (
	"fmt"

	quic "github.com/lucas-clemente/quic-go"
)

type Conn interface {
	Close()
	Streams() map[string]Stream
}

type WConn struct {
	ID      int
	streams map[string]Stream
	session *quic.Session
}

// OpenConn opens a new muliplexed connection from QUIC session and returns Conn interface
func OpenConn(session *quic.Session, channels []string) Conn {
	wconn := new(WConn)
	wconn.streams = make(map[string]Stream)
	for _, id := range channels {
		stream, err := (*session).OpenStreamSync()
		if err != nil {
			panic(err)
		}
		wstream := OpenStream(&stream)
		wconn.streams[id] = wstream
	}
	return wconn
}

// AcceptConn accepts a muliplexed connection from QUIC session and returns Conn interface
func AcceptConn(session *quic.Session, channels []string) Conn {
	wconn := new(WConn)
	wconn.streams = make(map[string]Stream)
	fmt.Println(channels)
	for _, id := range channels {
		stream, err := (*session).AcceptStream()
		if err != nil {
			panic(err)
		}
		wstream := OpenStream(&stream)
		wconn.streams[id] = wstream
	}
	return wconn
}

// Close closes the OrderedStream
func (wconn *WConn) Close() {
	if wconn.session == nil {
		return
	}
	(*wconn.session).Close(nil)
}

// Streams returns all streams making up the connection
func (wconn *WConn) Streams() map[string]Stream {
	return wconn.streams
}
