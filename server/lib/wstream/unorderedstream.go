package wstream

import quic "github.com/lucas-clemente/quic-go"

type UnorderedStream struct {
	ID          int
	streams     []*quic.Stream
	dataChannel chan []byte
}

// Open creates a new UnorderedStream by multiplexing existing QUIC streams
func (wstream *UnorderedStream) Open(streams []*quic.Stream) {
}

// ReadSync reads the next value synchronously from the UnorderedStream
func (wstream *UnorderedStream) ReadSync() []byte {
	bytes := <-wstream.dataChannel
	return bytes
}

// WriteSync writes a byte array synchronously to the stream
func (wstream *UnorderedStream) WriteSync(bytes []byte) {
}
