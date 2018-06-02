package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    // "github.com/waterloop/wcomms/wjson"
    // "github.com/waterloop/wcomms/wbinary"
    "log"
    "time"
    "strconv"
)

var upgrader = websocket.Upgrader{}

func main() {

    timeChan := time.NewTimer(time.Minute / 5).C
    tickChan := time.NewTicker(time.Second).C

    http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {

        var conn, _ = upgrader.Upgrade(w, r, nil)

        go func() {

            num := 0

                select {

                case <-timeChan:
                    log.Println("All packets sents")
                    break

                case <-tickChan:
                    conn.WriteMessage(1, []byte("Packet [" + strconv.Itoa(num) + "] sent!!"))
                    num++

                }
        }()

        // read data from the client
        go func(conn *websocket.Conn) {

            for {
                _, msg, _ := conn.ReadMessage()
                log.Println("Received: " + string(msg))
            }
        }(conn)
    })

    http.ListenAndServe(":3000", nil)
}
