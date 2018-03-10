package main

import (
	"fmt"
	"net"
)

const dataAddr = ":12345"
const commandAddr = ":12346"

func main() {
	addr, err := net.ResolveUDPAddr("udp", dataAddr)
	dataConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer dataConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := dataConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
