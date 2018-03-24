package main

import (
	"encoding/json"
	"fmt"
	"os"

	"./lib/tls"

	"github.com/lucas-clemente/quic-go"
	"github.com/waterloop/wcomms/wjson"
	"github.com/waterloop/wpool"
	"github.com/waterloop/wstream"
)

var podAddr string
var dataAddr string
var commandAddr string

const packetBufferSize = 1000

func main() {
	// Get port from environment variables
	podAddr = ":" + os.Getenv("POD_PORT")
	dataAddr = "255.255.255.255:" + os.Getenv("DATA_PORT")
	commandAddr = ":" + os.Getenv("COMMAND_PORT")

	config := quic.Config{IdleTimeout: 0}
	listener, err := quic.ListenAddr(podAddr, tls.GenerateConfig(), &config)
	CheckError(err)
	data := make(chan *wjson.CommPacketJson, packetBufferSize)

	fmt.Println("Server started")
	go func() {
		for {
			session, err := listener.Accept() // Wait for call and return a Conn
			if err != nil {
				break
			}
			go HandleClient(session, data)
		}
	}()
	go InitPool(data)
}

// InitPool creates a connection pool to relay data to the clients
func InitPool(dataChan <-chan *wjson.CommPacketJson) {
	wpool := wpool.CreateWPool(dataAddr, commandAddr)
	fmt.Println("Pool created")
	go wpool.Serve()
	for {
		packet := <-dataChan
		wpool.BroadcastPacket(packet)
	}
}

// HandleClient accepts a wstream connection from the pod
func HandleClient(session quic.Session, dataChan chan<- *wjson.CommPacketJson) {
	//defer session.Close(nil)
	wconn := wstream.AcceptConn(&session, []string{"sensor1", "sensor2", "sensor3", "command", "log"})
	fmt.Printf("%s %+v\n", "sss", wconn.Streams())
	for k, v := range wconn.Streams() {
		go HandleStream(k, v, dataChan)
	}
}

// HandleStream takes each stream and reads the packets being sent
func HandleStream(channel string, wstream wstream.Stream, dataChan chan<- *wjson.CommPacketJson) {
	defer wstream.Close()
	for {
		AcknowledgeMessage(wstream, 123)
		packet, err := wstream.ReadCommPacketSync()
		if err != nil {
			fmt.Println(err)
			continue
		}
		dataChan <- packet
		fmt.Printf("%s %+v\n", channel, packet)
		p, err := json.Marshal(packet)
		CheckError(err)
		LogPacket(p)
	}
}

// AcknowledgeMessage lets the client know a message was recieved
func AcknowledgeMessage(wstream wstream.Stream, id uint8) {
	packet := &wjson.CommPacketJson{
		Time: 1323,
		Type: "MessageRecieved",
		Id:   id,
		Data: []float32{32.2323, 1222.22, 2323.11},
	}
	wstream.WriteCommPacketSync(packet)
}

// CheckError checks and print errors
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
	}
}

// LogPacket logs the data in json format
func LogPacket(packet []byte) {
	f, err := os.OpenFile("logs/log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	CheckError(err)
	n, err := f.WriteString(string(packet) + "\n")
	_ = n
	CheckError(err)
}
