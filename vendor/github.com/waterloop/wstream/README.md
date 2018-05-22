# wstream
```import "github.com/shunr/wstream"```


## Usage

#### type Conn

```go
type Conn interface {
	Close()
	Streams() map[string]Stream
}
```


#### func  AcceptConn

```go
func AcceptConn(session *quic.Session, channels []string) Conn
```
AcceptConn accepts a muliplexed connection from QUIC session and returns Conn
interface

#### func  OpenConn

```go
func OpenConn(session *quic.Session, channels []string) Conn
```
OpenConn opens a new muliplexed connection from QUIC session and returns Conn
interface

#### type OrderedStream

```go
type OrderedStream struct {
	ID int
}
```


#### func (*OrderedStream) Close

```go
func (wstream *OrderedStream) Close()
```
Close closes the OrderedStream

#### func (*OrderedStream) ReadCommPacketSync

```go
func (wstream *OrderedStream) ReadCommPacketSync() (*wjson.CommPacketJson, error)
```
ReadCommPacketSync returns the next CommPacketJson in the OrderedStream,
blocking until completion

#### func (*OrderedStream) WriteCommPacketSync

```go
func (wstream *OrderedStream) WriteCommPacketSync(packet *wjson.CommPacketJson) error
```
WriteCommPacketSync takes a pointer to a CommPacketJson and writes it to the
OrderedStream, blocking until completion

#### type Stream

```go
type Stream interface {
	Close()
	ReadCommPacketSync() (*wjson.CommPacketJson, error)
	WriteCommPacketSync(packet *wjson.CommPacketJson) error
}
```


#### func  OpenStream

```go
func OpenStream(stream *quic.Stream) Stream
```
OpenStream creates a new OrderedStream from an existing QUIC stream

#### type WConn

```go
type WConn struct {
	ID int
}
```


#### func (*WConn) Close

```go
func (wconn *WConn) Close()
```
Close closes the OrderedStream

#### func (*WConn) Streams

```go
func (wconn *WConn) Streams() map[string]Stream
```
Streams returns all streams making up the connection
