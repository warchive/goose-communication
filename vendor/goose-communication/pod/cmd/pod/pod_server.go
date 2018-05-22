package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/waterloop/wcomms/wjson"
    "github.com/waterloop/wstream"
)

// CheckError Simple error verification
func CheckError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
    }
}

// HandleStream opens a new stream to send data over
func HandleStream(channel string, wstream wstream.Stream) {
    fmt.Println(channel)
    if channel != "command" {
        for {
            SendPacket(channel, 123, wstream)
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
