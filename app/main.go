package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Packet struct {
	Questions []question
	header    header
}

type question struct {
	Labels []string
	Type   uint
	Class  uint
}

type header struct {
	ID      uint16
	QR      uint
	Opcode  uint
	AA      uint
	TC      uint
	RD      uint
	RA      uint
	Z       uint
	RCODE   uint
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func NewHeader() *header {
	return &header{
		ID:      0,
		QR:      0,
		Opcode:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		RA:      0,
		Z:       0,
		RCODE:   0,
		QDCOUNT: 0,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
}

func (h *header) Parse(data []byte) {
	h.ID = binary.BigEndian.Uint16(data[0:2])
	h.QR = uint(data[2] & 15)
	// h.OPCODE
	h.AA = uint(data[2] & 10)
	h.TC = uint(data[2] & 9)
	h.RD = uint(data[2] & 8)
	h.RA = uint(data[2] & 7)
	// h.Z
	// h.RCODE
	h.QDCOUNT = binary.BigEndian.Uint16(data[4:6])
	h.ANCOUNT = binary.BigEndian.Uint16(data[6:8])
	h.NSCOUNT = binary.BigEndian.Uint16(data[8:10])
	h.ARCOUNT = binary.BigEndian.Uint16(data[10:12])
}

func (h header) Bytes() []byte {
	b := make([]byte, 12)

	binary.BigEndian.PutUint16(b[0:2], h.ID)
	binary.BigEndian.PutUint16(b[2:4], uint16(h.QR<<15|h.Opcode<<11|h.AA<<10|h.TC<<9|h.RD<<8|h.RA<<7|h.Z<<4|h.RCODE))
	binary.BigEndian.PutUint16(b[4:6], h.QDCOUNT)
	binary.BigEndian.PutUint16(b[6:8], h.ANCOUNT)
	binary.BigEndian.PutUint16(b[8:10], h.NSCOUNT)
	binary.BigEndian.PutUint16(b[10:12], h.ARCOUNT)

	return b
}

func NewQuestion() *question {
	return &question{
		Labels: []string{},
		Type:   0,
		Class:  0,
	}
}

func (q *question) Parse(data []byte, questions uint16) {
	// start at the first byte of the questions section
	offset := uint(0)

	// use the number of questions defined at the QDCOUNT header
	for i := 0; i < int(questions); i++ {
		// iterate while the current offset location isn't the null byte \x00
		for data[offset] != 0 {
			// use the offset location to determine label length
			labelLen := offset + uint(data[offset])
			// transform the label into a string and append to the labels slice
			q.Labels = append(q.Labels, string(data[offset+1:labelLen+1]))
			// increment the offset by the label length plus one to get the next length or the null byte
			offset = labelLen + 1
		}
	}
	// skip the null byte
	offset += 1
	// get the encoded type
	q.Type = uint(binary.BigEndian.Uint16(data[offset : offset+2]))
	// increment the offset to skip the type
	offset += 2
	// get the encoded class
	q.Class = uint(binary.BigEndian.Uint16(data[offset : offset+2]))
}

func (q question) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, label := range q.Labels {
		if err := buf.WriteByte(uint8(len(label))); err != nil {
			return nil, err
		}
		if _, err := buf.Write([]byte(label)); err != nil {
			return nil, err
		}
	}
	if err := buf.WriteByte(uint8(0)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(q.Type)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(q.Class)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s: \n% x\n%08b\n", size, source, string(receivedData), receivedData)

		h := NewHeader()
		h.Parse(receivedData[:12])
		h.QR = 1
		fmt.Printf("header: \n%08b\n%+v\n", h.Bytes(), h)

		q := NewQuestion()
		q.Parse(receivedData[12:], h.QDCOUNT)
		fmt.Printf("question: \n%+v\n", q)

		headerBytes := h.Bytes()
		questionBytes, err := q.Bytes()
		if err != nil {
			fmt.Println("failed to get bytes from question")
			return
		}

		b := append(headerBytes, questionBytes...)

		_, err = udpConn.WriteToUDP(b, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
