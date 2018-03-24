package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/waterloop/wcomms/wjson"
	"github.com/waterloop/wstream"
)

// CheckError Simple error verification
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

var addr string

func main() {
	addr = "localhost:" + os.Getenv("POD_PORT")
	config := quic.Config{RequestConnectionIDOmission: false}

	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, &config)
	CheckError(err)

	// Open multiple streams with waitgroup so main doesn't close before the streams finish sending
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()

	streams := []string{"sensor1", "sensor2", "sensor3", "command", "log"}
	wconn := wstream.OpenConn(&session, []string{"sensor1", "sensor2", "sensor3", "command", "log"})
	for _, k := range streams {
		go HandleStream(k, wconn.Streams()[k])
	}
	wg.Wait()
}

// HandleStream opens a new stream to send data over
func HandleStream(channel string, wstream wstream.Stream) {
	fmt.Println(channel)
	if (channel == "sensor1") || (channel == "sensor2") || (channel == "sensor3") {
		for {
			SendPacket(channel, 123, wstream)
			time.Sleep(time.Second)
		}
	} else {
		// SOMETHING IS HORRIBLY WRONG HERE
		for {
			SendPacket(channel, 123, wstream)
			time.Sleep(time.Second)
		}
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
