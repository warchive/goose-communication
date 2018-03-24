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
var dataChan chan *wjson.CommPacketJson
var commandChan chan *wjson.CommPacketJson

const packetBufferSize = 1024

func main() {
	// Get port from environment variables
	podAddr = ":" + os.Getenv("POD_PORT")
	dataAddr = "255.255.255.255:" + os.Getenv("DATA_PORT")
	commandAddr = ":" + os.Getenv("COMMAND_PORT")

	config := quic.Config{IdleTimeout: 0}
	listener, err := quic.ListenAddr(podAddr, tls.GenerateConfig(), &config)
	CheckError(err)

	dataChan = make(chan *wjson.CommPacketJson, packetBufferSize)
	commandChan = make(chan *wjson.CommPacketJson, packetBufferSize)

	fmt.Println("Server started")
	go func() {
		for {
			session, err := listener.Accept() // Wait for call and return a Conn
			if err != nil {
				break
			}
			go HandlePodConn(session)
		}
	}()
	InitPool()
}

// InitPool creates a connection pool to relay data to the clients
func InitPool() {
	wpool := wpool.CreateWPool(dataAddr, commandAddr)
	fmt.Println("Pool created")
	go wpool.Serve()
	go func() {
		for {
			packet := <-dataChan
			wpool.BroadcastPacket(packet)
		}
	}()

	for {
		command := wpool.GetNextCommand()
		commandChan <- command
	}
}

// HandlePodConn accepts a wstream connection from the pod
func HandlePodConn(session quic.Session) {
	//defer session.Close(nil)
	streams := []string{"sensor1", "sensor2", "sensor3", "command", "log"}
	wconn := wstream.AcceptConn(&session, []string{"sensor1", "sensor2", "sensor3", "command", "log"})
	for _, k := range streams {
		go HandleStream(k, wconn.Streams()[k])
	}
}

// HandleStream takes each stream and reads the packets being sent
func HandleStream(channel string, wstream wstream.Stream) {
	defer wstream.Close()
	if channel == "command" {
		for {
			command := <-commandChan
			wstream.WriteCommPacketSync(command)
		}
	} else {
		for {
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
