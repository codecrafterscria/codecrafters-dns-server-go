package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/dns"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s: \n% x\n%08b\n", size, source, string(receivedData), receivedData)

		p := dns.NewPacket()
		p.Parse(receivedData)
		p.Resolve()
		b, err := p.Bytes()
		if err != nil {
			fmt.Println("Failed to generate response:", err)
			os.Exit(1)
		}

		fmt.Printf("Sending response: \n%+v\n%08b\n", p, b)

		_, err = udpConn.WriteToUDP(b, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
