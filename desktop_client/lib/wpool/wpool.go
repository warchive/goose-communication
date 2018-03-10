package wpool

import (
	"net"
	// "github.com/mogball/wcomms/wjson"
)

// MaxConns stores the maximum number of clients allowed to connect
const MaxConns = 10

// Handler for recieving data with WPool
type Handler func(*WPool, net.Conn)

// WPool is a connection pool manager for TCP using net.Conn
// TODO change dataOut to wjson.CommPacketJSON
type WPool struct {
	port     string
	handler  Handler
	dataIn   chan []byte
	dataOut  chan []byte
	numConns int
}

// CreateWPool initializes and returns a WPool with a provided port
func CreateWPool(port string, handler Handler) *WPool {
	return &WPool{
		port:     port,
		handler:  handler,
		dataIn:   make(chan []byte),
		dataOut:  make(chan []byte),
		numConns: 0,
	}
}

// Serve starts the connection pool and adds / closes connections
// TODO implement connections being closed
func (pool *WPool) Serve() {
	connChannel := make(chan net.Conn)
	newConn, err := net.Listen("tcp", pool.port)
	if err != nil {
		panic(err)
	}

	// Goroutine for connection queue
	go func(connChannel chan net.Conn, output chan []byte) {
		var connections [MaxConns]net.Conn
		for {
			select {
			case conn := <-connChannel:
				connections[pool.numConns] = conn
				pool.numConns++
			case data := <-output:
				for i := 0; i < pool.numConns; i = i + 1 {
					connections[i].Write(data)
				}
			}
		}
	}(connChannel, pool.dataOut)

	// Add new connections
	for {
		if pool.numConns >= MaxConns {
			continue
		}
		conn, err := newConn.Accept()
		if err != nil {
			panic(err)
		}
		go pool.handler(pool, conn)
		connChannel <- conn
	}
}

// DataHandler reads data from the connected clients
func DataHandler(pool *WPool, conn net.Conn) {
	recvChannel := make(chan []byte)
	errChannel := make(chan error)

	// Goroutine for receiving from client
	go func(recvChannel chan []byte, errChannel chan error) {
		for {
			data := make([]byte, 1024)
			_, err := conn.Read(data)
			if err != nil {
				errChannel <- err
				return
			}
			recvChannel <- data
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
