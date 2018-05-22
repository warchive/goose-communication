package main

import (
    "os"
    "github.com/lucas-clemente/quic-go"
    "fmt"

    "./tls"
    "github.com/waterloop/wcomms/wjson"
)

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
