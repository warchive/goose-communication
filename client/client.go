package main

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"../server/lib/wstream"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/mogball/wcomms/wjson"
	// "github.com/buger/jsonparser"
)

// CheckError Simple error verification
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// Different addresses and ports to send data to
const addr = "localhost:10000"

func main() {
	config := quic.Config{RequestConnectionIDOmission: false}

	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, &config)
	CheckError(err)

	// Open multiple streams with waitgroup so main doesn't close before the streams finish sending
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()
	wconn := wstream.OpenConn(&session, []string{"sensor1", "sensor2", "sensor3", "command", "log"})
	for k, v := range wconn.Streams() {
		go handleStream(k, v)
	}
	wg.Wait()
}

// Open a new stream to send data over
func handleStream(channel string, wstream wstream.Stream) {
	defer wstream.Close()
	if (channel == "sensor1") || (channel == "sensor2") || (channel == "sensor3") {
		for {
			packet, err := wstream.ReadCommPacketSync()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("%s %+v\n", channel, packet)
		}
	} else {
		for {
			sendPacket(channel, wstream)
			time.Sleep(time.Second)
		}
	}
}

func sendPacket(channel string, wstream wstream.Stream) {
	packet := &wjson.CommPacketJson{
		Time: 1323,
		Type: channel,
		Name: channel,
		Data: []float32{32.2323, 1222.22, 2323.11},
	}
	wstream.WriteCommPacketSync(packet)
}
