# WPool

Connection pooling for backend server in Go using UDP for broadcasting data and TCP for relaying commands from the client

## Usage

#### type WPool

``` go
type WPool struct {
	dataAddr    string
	commandAddr string
	commandIn   chan *wjson.CommPacketJson
	dataOut     chan *wjson.CommPacketJson
	numConns    int
	connections []net.Conn
}
```

#### func CreateWPool

``` go
func CreateWPool(dataAddr string, commandAddr string) *WPool
```

Creates a new WPool that broadcasts and recieves commands from the specified addresses

#### func Serve

``` go
func (pool *WPool) Serve()
```

Function that starts the connection pool for TCP and UDP, calls ServeDataChannel and ServeCommandChannel

#### func ServeDataChannel

``` go 
func (pool *WPool) ServeDataChannel()
```

Broadcasts data to the WPool's port using UDP via channels

#### func ServeCommandChannel

``` go 
func (pool *WPool) ServeCommandChannel()
```

Recieves TCP connections for commands and adds them to the connection pool

#### func GetNextCommand

``` go
func (pool *WPool) GetNextCommand() *wjson.CommPacketJson 
```

Returns the next command in the commandIn channel as a JSON

#### func CommandHandler

``` go
func (pool *WPool) CommandHandler(conn net.Conn, connIndex int) 
```

Handles commands via TCP from the specified connection in the pool and passes them to the commandIn channel

#### func CloseConn

``` go
func (pool *WPool) CloseConn(conn net.Conn, index int) 
```
Closes the specified TCP connection in the pool at the specified index

#### func BroadcastPacket

``` go
func (pool *WPool) BroadcastPacket(packet *wjson.CommPacketJson)
```
Passes the given packet into the WPool's dataOut channel to be broadcasted

#### func SendPacketByteArray

``` go
func SendPacketByteArray(dataConn *net.UDPConn, data *wjson.CommPacketJson) 
```

Writes the next packet in the dataOut channel to the WPool's broadcast address 
