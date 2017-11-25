package main

import (
	"fmt"
	"os"
	"time"
	// "bufio"

	"./lib/tls"
	"./lib/wstream"

	"github.com/lucas-clemente/quic-go"
	"github.com/mogball/wcomms/wjson"
	// "github.com/tarm/serial"
)

const addr = ":10000"

var i int

func main() {
	// Choose port to listen from
	config := quic.Config{IdleTimeout: 0}
	listener, err := quic.ListenAddr(addr, tls.GenerateConfig(), &config)
	checkError(err)

	fmt.Println("Server started")

	/*
		c := &serial.Config{Name: "COM3", Baud: 9600}
		s, err := serial.OpenPort(c)
		checkError(err)
		reader := bufio.NewReader(s)
	*/
	for {
		/*
			r, err := reader.ReadBytes(255)
			checkError(err)
			fmt.Println(r)
		*/
		session, err := listener.Accept() // Wait for call and return a Conn
		if err != nil {
			break
		}

		go handleClient(session)

	}
}

func handleClient(session quic.Session) {
	//defer session.Close(nil)
	wconn := wstream.AcceptConn(&session, []string{"sensor1", "sensor2", "sensor3", "command", "log"})
	fmt.Printf("%s %+v\n", "sss", wconn.Streams())
	for k, v := range wconn.Streams() {
		go handleStream(k, v)
	}
}

func handleStream(channel string, wstream wstream.Stream) {
	defer wstream.Close()
	if (channel == "sensor1") || (channel == "sensor2") || (channel == "sensor3") {
		for {
			acknowledgeMessage(wstream, "sensor data")
			time.Sleep(time.Second)
		}
	} else {
		for {
			packet, err := wstream.ReadCommPacketSync()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("%s %+v\n", channel, packet)
		}
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
