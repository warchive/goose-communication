package main

import (
    "encoding/json"
    "fmt"
    "net"
    "time"

    "github.com/waterloop/wcomms/wjson"
)

const dataAddr = ":42001"
const commandAddr = "localhost:42002"

func main() {
    addr, err := net.ResolveUDPAddr("udp", dataAddr)
    checkError(err)
    dataConn, err := net.ListenUDP("udp", addr)
    checkError(err)
    commandConn, err := net.Dial("tcp", commandAddr)
    checkError(err)

    defer dataConn.Close()
    defer commandConn.Close()

    buf := make([]byte, 1024)

    go func() {
        for {
            n, addr, err := dataConn.ReadFromUDP(buf)
            fmt.Println("Received ", string(buf[0:n]), " from ", addr)
            if err != nil {
                fmt.Println("Error: ", err)
            }
        }
    }()

    for {
        time.Sleep(time.Second)
        data := wjson.CommPacketJson{
            Time: 1323,
            Type: "command",
            Id:   111,
            Data: []float32{32.2323, 1222.22, 2323.11},
        }
        packet, _ := json.Marshal(data)
        commandConn.Write(packet)
        fmt.Println("Sending ", string(packet))
    }
}

func checkError(err error, callback ...func(err error)) {
    if err != nil {
        if callback == nil {
            panic(err)
        } else {
            callback[0](err)
        }
    }
}
