package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"./lib/tls"
	"./lib/wstream"

	"github.com/buger/jsonparser"
	"github.com/lucas-clemente/quic-go"
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
	start := time.Now()
	var wstream wstream.Stream = new(wstream.OrderedStream)
	wstream.Open(stream)
	i := 0
	for {
		bytes := wstream.ReadSync()
		i++
		if i%100 == 0 {
			fmt.Println(time.Duration(int64(time.Since(start)) / int64(i)))
			fmt.Printf("%s\n", string(bytes))
		}
		id, iderr := jsonparser.GetString(bytes, "id")
		if iderr == nil {
			acknowledgeMessage(wstream, id, true)
		}
		if i%100 == 99 {
			fmt.Printf("%d\n", i)
		}
	}
}

// Let client know message was recieved
func acknowledgeMessage(wstream wstream.Stream, id string, success bool) {
	msg := map[string]interface{}{"id": id, "type": "recieved", "success": success}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	wstream.WriteSync(bytes)
}

// Check and print errors
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
	}
}
