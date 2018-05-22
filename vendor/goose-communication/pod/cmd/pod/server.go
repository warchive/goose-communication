package main

import (
    "os"
    "github.com/lucas-clemente/quic-go"
    "crypto/tls"
    "sync"
    "github.com/waterloop/wstream"
)

func main() {
    addr := "localhost:" + os.Getenv("POD_PORT")
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
