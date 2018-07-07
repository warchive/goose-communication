package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/waterloop/wcomms/wbinary"
    "log"
    "time"
    "strconv"
    "encoding/json"
    "os"
    "github.com/waterloop/wcomms/wjson"
)



var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func PrintPackate(packet *wbinary.CommPacket) {
    log.Println("Type: " + string(packet.PacketType) + " ID: " + string(packet.PacketId) +
        "Data: [" + strconv.FormatFloat(float64(packet.Data1), 'f', 2, 32) + "," +
        strconv.FormatFloat(float64(packet.Data2), 'f', 2, 32) + "," +
        strconv.FormatFloat(float64(packet.Data3), 'f', 2, 32)  + "]")
}

func main() {
    // create initial mag wheels packet
    var magPacket = &wjson.CommPacketJson{
        Time: int64(time.Now().Second()),
        Type: "sensor",
        Data: []float32{0},
    }


    // create initial friction packet
    var frictionPacket = &wjson.CommPacketJson{
        Time: int64(time.Now().Second()),
        Type: "sensor",
        Data: []float32{0},
    }

    // create initial levitation packet
    var levPacket = &wjson.CommPacketJson{
        Time: int64(time.Now().Second()),
        Type: "sensor",
        Data: []float32{0},
    }

    // which packet to send
    packetType := 0
    var err error

    http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {

        var conn, _ = upgrader.Upgrade(w, r, nil)

        go func() {
            for {
                var bytesToSend []byte

                if packetType == 0 {
                    bytesToSend, err = json.Marshal(magPacket)
                    if err != nil {
                        log.Println(err.Error())
                        os.Exit(2)
                    }
                    magPacket.Data[0] += 1
                    magPacket.Time = int64(time.Now().Second())

                    if magPacket.Data[0]> 100 {
                        magPacket.Data[0] = 0
                    }
                } else if packetType == 1 {
                    bytesToSend, err = json.Marshal(frictionPacket)
                    if err != nil {
                        log.Println(err.Error())
                        os.Exit(2)
                    }
                    frictionPacket.Data[0]+= 1
                    frictionPacket.Time = int64(time.Now().Second())

                    if frictionPacket.Data[0] > 100 {
                        frictionPacket.Data[0] = 0
                    }
                } else if packetType == 2 {
                    bytesToSend, err = json.Marshal(levPacket)
                    if err != nil {
                        log.Println(err.Error())
                        os.Exit(2)
                    }
                    levPacket.Data[0] += 1
                    levPacket.Time = int64(time.Now().Second())

                    if levPacket.Data[0]> 100 {
                        levPacket.Data[0] = 0
                    }
                } else {
                    packetType = 0
                    continue
                }

                log.Println(packetType)

                // Send the binary bytes of the original packet to the connection.
                conn.WriteMessage(websocket.TextMessage, bytesToSend)

                log.Printf("Sent this json %s, \n", string(bytesToSend))

                // Sleep before sending the next packet.
                time.Sleep(time.Second * 1)
                packetType += 1
            }
        }()

        // Read data from the client and print to the console
        go func(conn *websocket.Conn) {
            for {
                // Try to receive/read the bytes (that were sent).
                receivedType, receivedBytes, _ := conn.ReadMessage()

                // we only care if the message is string
                if receivedType == websocket.TextMessage {
                    log.Printf("Received json %s\n", string(receivedBytes))
                }
            } // End of for-loop.
        }(conn)
    })

    http.ListenAndServe(":6500", nil)
}
