package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/waterloop/wcomms/wjson"
    "github.com/waterloop/wcomms/wbinary"
    "log"
    "time"
    // "strconv"
)

var upgrader = websocket.Upgrader{}

func main() {

    expectedPacket := &wbinary.CommPacket {
        PacketType: wbinary.State,
        PacketId:   54,
        Data1:      -724.875,
        Data2:      846.5,
        Data3:      442.5625,
    }

    toSend := []byte{178, 157, 26, 78, 167, 88, 234, 94}

    timeChan := time.NewTimer(time.Minute / 5).C
    tickChan := time.NewTicker(time.Second).C

    http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {

        var conn, _ = upgrader.Upgrade(w, r, nil)

        go func() {

                select {

                case <-timeChan:
                    log.Println("Packet Sent!")

                case <-tickChan:
                    // Read the 8-byte buffer and make a packet out of it.
                    packetToSend := wbinary.ReadPacket(toSend)

                    // Encode that packet, so we can send it.
                    encodedPacket, thereWasAnError := wjson.PacketEncodeJson(packetToSend)
                    if thereWasAnError != nil {
                        panic(thereWasAnError)
                    }

                    // Send the encoded packet to the connection.
                    conn.WriteMessage(1, encodedPacket)

                    log.Println("Sent an encoded packet: ")
                    log.Printf("%s", encodedPacket)
                }
        }()

        // Read data from the client
        go func(conn *websocket.Conn) {

            for {
                // Try to receive/read the encoded packet (that was sent).
                _, receivedPacket, _ := conn.ReadMessage()

                log.Println("Received an encoded packet: ")
                log.Printf("%s", receivedPacket)

                // Decode the recieved packet (encoded).
                decodedPacket, thereWasAnError := wjson.PacketDecodeJson(receivedPacket)
                if thereWasAnError != nil {
                    panic(thereWasAnError)
                }

                if *decodedPacket == *expectedPacket {
                    log.Println("Success!!!!!!!!!!!!!!!!")
                } else {
                    log.Println("Expected a different packet.")
                }
            }

        }(conn)
    })

    http.ListenAndServe(":3000", nil)
}
