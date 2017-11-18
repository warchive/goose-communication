package wstream

import quic "github.com/lucas-clemente/quic-go"

type Stream interface {
	Open(stream *quic.Stream)
	ReadSync() []byte
	WriteSync(bytes []byte)
}
