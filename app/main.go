package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type DNSHeader struct {
	ID      uint
	QR      uint
	Opcode  uint
	AA      uint
	TC      uint
	RD      uint
	RA      uint
	Z       uint
	RCODE   uint
	QDCOUNT uint
	ANCOUNT uint
	NSCOUNT uint
	ARCOUNT uint
}

func NewDNSHeader() *DNSHeader {
	return &DNSHeader{}
}

func (h *DNSHeader) Parse(data []byte) {
	fmt.Printf("Data: %08b\n", data)
}

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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		h := NewDNSHeader()
		h.Parse(buf[:12])

		// Create an empty response
		response := make([]byte, 12)
		binary.BigEndian.PutUint16(response[0:2], 1234)
		response[2] = 1 << 7

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
