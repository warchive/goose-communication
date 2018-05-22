package wpool

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/waterloop/wcomms/wjson"
)

// MaxConns stores the maximum number of clients allowed to connect
const MaxConns = 3

// PacketBufferSize is the maximum number of packets to keep in a data channel
const PacketBufferSize = 1024

// WPool is a connection pool manager for UDP using net.Conn
type WPool struct {
	dataAddr    string
	commandAddr string
	commandIn   chan *wjson.CommPacketJson
	dataOut     chan *wjson.CommPacketJson
	numConns    int
	connections []net.Conn
}

// CreateWPool initializes and returns a WPool with a provided port
func CreateWPool(dataAddr string, commandAddr string) *WPool {
	return &WPool{
		dataAddr:    dataAddr,
		commandAddr: commandAddr,
		commandIn:   make(chan *wjson.CommPacketJson, PacketBufferSize),
		dataOut:     make(chan *wjson.CommPacketJson, PacketBufferSize),
		numConns:    0,
		connections: make([]net.Conn, 0),
	}
}

// Serve starts the connection pool and adds / closes connections for TCP and UDP
func (pool *WPool) Serve() {
	fmt.Println("Serving data channel")
	pool.ServeDataChannel()
	fmt.Println("Serving commands channel")
	pool.ServeCommandChannel()
}

// ServeDataChannel broadcasts data to the specified port using UDP
func (pool *WPool) ServeDataChannel() {
	broadcast, err := net.ResolveUDPAddr("udp", pool.dataAddr)
	dataConn, err := net.DialUDP("udp", nil, broadcast)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			data := <-pool.dataOut
			SendPacketByteArray(dataConn, data)
		}
	}()
}

// ServeCommandChannel takes connections through tcp and adds them to the connection pool
func (pool *WPool) ServeCommandChannel() {
	commandConn, err := net.Listen("tcp", pool.commandAddr)
	if err != nil {
		panic(err)
	}

	for {
		if pool.numConns >= MaxConns {
			continue
		}
		conn, err := commandConn.Accept()
		fmt.Println("Connected")
		if err != nil {
			continue
		}
		pool.numConns++
		pool.connections = append(pool.connections, conn)
		go pool.CommandHandler(conn, pool.numConns-1)
	}
}

// CommandHandler reads data from the connected clients
func (pool *WPool) CommandHandler(conn net.Conn, connIndex int) {
	buf := make([]byte, 1024)
	// Goroutine for receiving from client
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			pool.CloseConn(conn, connIndex)
			return
		}
		packet := &wjson.CommPacketJson{}
		err = json.Unmarshal(buf[:n], packet)
		if err == nil {
			pool.commandIn <- packet
			pool.BroadcastPacket(packet)
		} else {
			fmt.Println(err)
			pool.CloseConn(conn, connIndex)
			return
		}
	}
}

// GetNextCommand gets the next command
func (pool *WPool) GetNextCommand() *wjson.CommPacketJson {
	command := <-pool.commandIn
	return command
}

// CloseConn closes a TCP connection in the pool
func (pool *WPool) CloseConn(conn net.Conn, index int) {
	pool.connections[index] = pool.connections[pool.numConns-1]
	pool.connections = pool.connections[:pool.numConns-1]
	pool.numConns--
	conn.Close()
	fmt.Println("Connection", index, "closed")
}

// BroadcastPacket sets the current dataOut to the provided packet
func (pool *WPool) BroadcastPacket(packet *wjson.CommPacketJson) {
	pool.dataOut <- packet
}

// SendPacketByteArray writes data to BroadcastAddr
func SendPacketByteArray(dataConn *net.UDPConn, data *wjson.CommPacketJson) {
	packet, err := json.Marshal(*data)
	_, err = dataConn.Write(packet)
	fmt.Println("Broadcasting ", string(packet))
	if err != nil {
		panic(err)
	}
}
