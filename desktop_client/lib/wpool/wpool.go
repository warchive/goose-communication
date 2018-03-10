package wpool

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/mogball/wcomms/wjson"
)

// MaxConns stores the maximum number of clients allowed to connect
const MaxConns = 10

// BroadcastAddr is the address where data is broadcasted to
const BroadcastAddr = "255.255.255.255:12345"

// Handler for recieving data with WPool
type Handler func(*WPool, net.Conn)

// WPool is a connection pool manager for UDP using net.Conn
type WPool struct {
	dataAddr    string
	commandAddr string
	handler     Handler
	dataIn      chan wjson.CommPacketJson
	dataOut     chan wjson.CommPacketJson
	numConns    int
}

// CreateWPool initializes and returns a WPool with a provided port
func CreateWPool(dataAddr string, commandAddr string, handler Handler) *WPool {
	return &WPool{
		dataAddr:    dataAddr,
		commandAddr: commandAddr,
		handler:     handler,
		dataIn:      make(chan wjson.CommPacketJson),
		dataOut:     make(chan wjson.CommPacketJson),
		numConns:    0,
	}
}

// Serve starts the connection pool and adds / closes connections
func (pool *WPool) Serve() {
	DataAddr, err := net.ResolveUDPAddr("udp", BroadcastAddr)
	dataConn, err := net.DialUDP("udp", nil, DataAddr)
	fmt.Println(dataConn)
	commandConn, err := net.Listen("tcp", pool.commandAddr)

	connChannel := make(chan net.Conn)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			data := <-pool.dataOut
			fmt.Println(data)
			SendPacketByteArray(dataConn, data)
		}
	}()

	// Goroutine for connection queue
	go func(connChannel chan net.Conn, output chan wjson.CommPacketJson) {
		var connections [MaxConns]net.Conn
		for {
			select {
			case conn := <-connChannel:
				connections[pool.numConns] = conn
				pool.numConns++
			case data := <-output:
				for i := 0; i < pool.numConns; i = i + 1 {
					packet, err := json.Marshal(data)
					if err == nil {
						connections[i].Write(packet)
					}
				}
			}
		}
	}(connChannel, pool.dataOut)

	// Add new connections
	for {
		if pool.numConns >= MaxConns {
			continue
		}
		conn, err := commandConn.Accept()
		fmt.Println("Connected")
		if err != nil {
			panic(err)
		}
		go pool.handler(pool, conn)
		connChannel <- conn
	}
}

// CommandHandler reads data from the connected clients
func CommandHandler(pool *WPool, conn net.Conn) {
	recvChannel := make(chan wjson.CommPacketJson)
	errChannel := make(chan error)

	// Goroutine for receiving from client
	go func(recvChannel chan wjson.CommPacketJson, errChannel chan error) {
		for {
			data := make([]byte, 1024)
			_, err := conn.Read(data)
			if err != nil {
				errChannel <- err
				return
			}
			packet := wjson.CommPacketJson{}
			err = json.Unmarshal(data, &packet)
			if err == nil {
				recvChannel <- packet
			}
		}
	}(recvChannel, errChannel)

	// Listen to broadcast
	for {
		select {
		// Receive from client and write same data back to broadcast
		// TODO recieve dataOut from desktop_client.go
		case data := <-recvChannel:
			pool.dataIn <- data
		// Close connection on error
		case <-errChannel:
			pool.numConns--
			conn.Close()
			return
		}
	}
}

// BroadcastPacket sets the current dataOut to the provided packet
// TODO pass in data from desktop client
func (pool *WPool) BroadcastPacket() {
	packet := wjson.CommPacketJson{
		Time: 1323,
		Type: "State",
		Id:   122,
		Data: []float32{32.2323, 1222.22, 2323.11},
	}
	pool.dataOut <- packet
}

// SendPacketByteArray writes data to BroadcastAddr
func SendPacketByteArray(dataConn *net.UDPConn, data wjson.CommPacketJson) {
	packet, err := json.Marshal(data)
	_, err = dataConn.Write(packet)
	if err != nil {
		panic(err)
	}
	fmt.Println("Broadcasting!")
	fmt.Println(packet)
}
