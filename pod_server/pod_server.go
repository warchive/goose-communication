package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"../desktop_client/lib/wstream"
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
		go HandleStream(k, v)
	}
	wg.Wait()
}

// HandleStream opens a new stream to send data over
func HandleStream(channel string, wstream wstream.Stream) {
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
			SendPacket(channel, 123, wstream)
			time.Sleep(time.Second)
		}
	}
}

// SendPacket takes a CommPacketJson to send back to server and log the data
func SendPacket(channel string, id uint8, wstream wstream.Stream) {
	packet := &wjson.CommPacketJson{
		Time: 1323,
		Type: channel,
		Id:   id,
		Data: []float32{32.2323, 1222.22, 2323.11},
	}
	wstream.WriteCommPacketSync(packet)
	p, err := json.Marshal(packet)
	CheckError(err)
	LogPacket(p)
}

// LogPacket logs the data in json format
func LogPacket(packet []byte) {
	f, err := os.OpenFile("logs/log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	CheckError(err)
	n, err := f.WriteString(string(packet) + "\n")
	_ = n
	CheckError(err)
}
