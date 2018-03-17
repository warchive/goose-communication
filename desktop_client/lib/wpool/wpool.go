package wpool

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/mogball/wcomms/wjson"
)

// MaxConns stores the maximum number of clients allowed to connect
const MaxConns = 3

// PacketBufferSize is the maximum number of packets to keep in a data channel
const PacketBufferSize = 100

// WPool is a connection pool manager for UDP using net.Conn
type WPool struct {
	dataAddr    string
	commandAddr string
	dataIn      chan *wjson.CommPacketJson
	dataOut     chan *wjson.CommPacketJson
	numConns    int
	connections []net.Conn
}

// CreateWPool initializes and returns a WPool with a provided port
func CreateWPool(dataAddr string, commandAddr string) *WPool {
	return &WPool{
		dataAddr:    dataAddr,
		commandAddr: commandAddr,
		dataIn:      make(chan *wjson.CommPacketJson, PacketBufferSize),
		dataOut:     make(chan *wjson.CommPacketJson, PacketBufferSize),
		numConns:    0,
		connections: make([]net.Conn, 0),
	}
}

// Serve starts the connection pool and adds / closes connections
func (pool *WPool) Serve() {
	fmt.Println("Serving data channel")
	pool.serveDataChannel()
	fmt.Println("Serving commands channel")
	pool.serveCommandChannel()
}

func (pool *WPool) serveDataChannel() {
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

func (pool *WPool) serveCommandChannel() {
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
		go pool.commandHandler(conn, pool.numConns-1)
	}
}

// CommandHandler reads data from the connected clients
func (pool *WPool) commandHandler(conn net.Conn, connIndex int) {
	recvChannel := make(chan *wjson.CommPacketJson)
	buf := make([]byte, 1024)
	// Goroutine for receiving from client
	go func() {

		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				pool.closeConn(conn, connIndex)
				return
			}
			packet := &wjson.CommPacketJson{}
			err = json.Unmarshal(buf[:n], packet)
			if err == nil {
				recvChannel <- packet
			} else {
				fmt.Println(err)
				pool.closeConn(conn, connIndex)
				return
			}
		}
	}()

	// Broadcast recvd command to all clients
	// TODO: Send to pod
	for {
		packet := <-recvChannel
		pool.dataIn <- packet
		pool.BroadcastPacket(packet)
	}
}

func (pool *WPool) closeConn(conn net.Conn, index int) {
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
