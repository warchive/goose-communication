package main

import (
	"fmt"
	"os"

	"./lib/tls"
	"./lib/wstream"

	"github.com/lucas-clemente/quic-go"
	"github.com/mogball/wcomms/wjson"
)

const addr1 = ":10000"
const addr2 = ":12345"

var i int

func main() {
	// Choose port to listen from
	config := quic.Config{IdleTimeout: 0}
	listener1, err := quic.ListenAddr(addr1, tls.GenerateConfig(), &config)
	checkError(err)
	listener2, err := quic.ListenAddr(addr2, tls.GenerateConfig(), &config)
	checkError(err)
	fmt.Println("Server started")
	for {
		session1, err := listener1.Accept() // Wait for call and return a Conn
		session2, err := listener2.Accept() // Wait for call and return a Conn
		if err != nil {
			break
		}
		go handleClient(session1)
		go handleClient(session2)
	}
}

func handleClient(session quic.Session) {
	defer session.Close(nil)
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			fmt.Println(err)
			break
		} else {
			go handleStream(&stream)
		}
	}
}

func handleStream(stream *quic.Stream) {
	var wstream wstream.Stream = new(wstream.OrderedStream)
	wstream.Open(stream)
	defer wstream.Close()
	for {
		packet, err := wstream.ReadCommPacketSync()
		if err != nil {
			fmt.Println(err)
			continue
		}
		acknowledgeMessage(wstream, packet.Name)
		fmt.Printf("%+v\n", packet)
	}
}

// Let client know message was recieved
func acknowledgeMessage(wstream wstream.Stream, name string) {
	packet := &wjson.CommPacketJson{
		Time: 1323,
		Type: "State",
		Name: name,
		Data: []float32{32.2323, 1222.22, 2323.11},
	}
	wstream.WriteCommPacketSync(packet)
}

// Check and print errors
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
	}
}
