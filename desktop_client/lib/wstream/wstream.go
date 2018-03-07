package wstream

import (
	"encoding/json"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/mogball/wcomms/wjson"
)

// Stream provides read and write capabiltiies with CommPackets
type Stream interface {
	Close()
	ReadCommPacketSync() (*wjson.CommPacketJson, error)
	WriteCommPacketSync(packet *wjson.CommPacketJson) error
}

// OrderedStream is used for order-sensitive data
type OrderedStream struct {
	ID          int
	stream      *quic.Stream
	dataChannel chan []byte
}

// OpenStream creates a new OrderedStream from an existing QUIC stream
func OpenStream(stream *quic.Stream) Stream {
	wstream := new(OrderedStream)
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
	return wstream
}

// Close closes the OrderedStream
func (wstream *OrderedStream) Close() {
	if wstream.stream == nil {
		return
	}
	(*wstream.stream).Close()
	close(wstream.dataChannel)
}

// ReadCommPacketSync returns the next CommPacketJson in the OrderedStream, blocking until completion
func (wstream *OrderedStream) ReadCommPacketSync() (*wjson.CommPacketJson, error) {
	encoded := wstream.ReadBytes()
	packet := &(wjson.CommPacketJson{})
	err := json.Unmarshal(encoded, packet)
	return packet, err
}

// WriteCommPacketSync takes a pointer to a CommPacketJson and writes it to the OrderedStream, blocking until completion
func (wstream *OrderedStream) WriteCommPacketSync(packet *wjson.CommPacketJson) error {
	bytes, err := json.Marshal(packet)
	wstream.WriteBytes(bytes)
	return err
}

// ReadBytes reads bytes from a stream
func (wstream *OrderedStream) ReadBytes() []byte {
	bytes := <-wstream.dataChannel
	return bytes
}

// WriteBytes writes bytes to a stream
func (wstream *OrderedStream) WriteBytes(bytes []byte) {
	(*wstream.stream).Write(bytes)
}
