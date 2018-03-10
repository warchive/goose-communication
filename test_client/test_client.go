package main

import (
	"fmt"
	"net"
)

const dataAddr = "localhost:12345"
const commandAddr = ":12346"

func main() {
	// Choose port to listen from
	ServerAddr, err := net.ResolveUDPAddr("udp", dataAddr)
	dataConn, err := net.DialUDP("udp", nil, ServerAddr)
	//commandConn, err := net.Dial("tcp", commandAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dataConn.Close()
	//defer commandConn.Close()

	buffer := make([]byte, 1024)
	for {
		fmt.Println("Hello")
		fmt.Println(dataConn)
		n, _, err := dataConn.ReadFromUDP(buffer)
		fmt.Println(n)
		if err != nil {
			fmt.Println(string(buffer[0:n]))
		} else {
			fmt.Println(err)
		}
	}
}
