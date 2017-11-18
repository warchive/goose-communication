package wstream

import quic "github.com/lucas-clemente/quic-go"

type OrderedStream struct {
	ID          int
	stream      *quic.Stream
	dataChannel chan []byte
}

// Open creates a new OrderedStream from an existing QUIC stream
func (wstream *OrderedStream) Open(stream *quic.Stream) {
	wstream.stream = stream
	buf := make([]byte, 1024)
	wstream.dataChannel = make(chan []byte)
	go func() {
		for {
			n, err := (*stream).Read(buf)
			if err == nil {
				bytes := buf[0:n]
				wstream.dataChannel <- bytes
			} else {
				break
			}
		}
	}()
}

// ReadSync reads the next value synchronously from the OrderedStream
func (wstream *OrderedStream) ReadSync() []byte {
	bytes := <-wstream.dataChannel
	return bytes
}

// WriteSync writes a byte array synchronously to the stream
func (wstream *OrderedStream) WriteSync(bytes []byte) {
	(*wstream.stream).Write(bytes)
}
