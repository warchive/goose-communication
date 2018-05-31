package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    "log"
    "time"
    "strconv"
)

var upgrader = websocket.Upgrader{}

func main() {
    timeChan := time.NewTimer(time.Minute * 60).C
    tickChan := time.NewTicker(time.Millisecond * 1000).C

    num := 0


    http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
        var conn, _ = upgrader.Upgrade(w, r, nil)


        go func() {
            for {
                select {
                case <-timeChan:
                    log.Println("All packets sents")
                case <-tickChan:
                    conn.WriteMessage(1, []byte("Packet number "+strconv.Itoa(num)+" sent!!"))
                    num++
                }
            }
        }()

        // read data from the client
        go func(conn *websocket.Conn) {
             for  {
                  _, msg, _ := conn.ReadMessage()
                  log.Println("Received: " + string(msg))
             }
        }(conn)
    })

    http.ListenAndServe(":3000", nil)
}
