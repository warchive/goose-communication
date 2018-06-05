package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/waterloop/wcomms/wjson"
    "github.com/waterloop/wcomms/wbinary"
    "log"
    "time"
)

// originalPacket converts to := []byte{178, 157, 26, 78, 167, 88, 234, 94}
var originalPacket = &wbinary.CommPacket{
    PacketType: wbinary.State,
    PacketId:   54,
    Data1:      -724.875,
    Data2:      846.5,
    Data3:      442.5625,
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func main() {
    http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {

        var conn, _ = upgrader.Upgrade(w, r, nil)

        go func() {

            for {
                // Write the packet after converting it to the 8-byte buffer.
                bytesToSend := wbinary.WritePacket(originalPacket)

                // Send the binary bytes of the original packet to the connection.
                conn.WriteMessage(websocket.BinaryMessage, bytesToSend)
                log.Printf("Sent these bytes: [%x]\n", bytesToSend)

                // Sleep before sending the next packet.
                time.Sleep(time.Second * 5)

            }
        }()

        // Read data from the client
        go func(conn *websocket.Conn) {
            for {
                // Try to receive/read the bytes (that were sent).
                receivedType, receivedBytes, _ := conn.ReadMessage()

                // We only care if the message is binary (discard the rest).
                if receivedType == websocket.BinaryMessage {
                    log.Printf("Recieved these bytes: [%x]\n", receivedBytes)

                    // Read the received bytes and make a packet out of them.
                    receivedPacket := wbinary.ReadPacket(receivedBytes)

                    // Encode that packet, to make it pretty (JSON).
                    encodedPacket, thereWasAnError := wjson.PacketEncodeJson(receivedPacket)
                    if thereWasAnError != nil {
                        panic(thereWasAnError)
                    } else {
                        log.Println("Encoded the received packet to JSON: ")
                        log.Printf("%s\n", encodedPacket)
                    }
                } // End of check to see if the message we received was indeed binary.
            } // End of for-loop.
        }(conn)
    })

    http.ListenAndServe(":3000", nil)
}
