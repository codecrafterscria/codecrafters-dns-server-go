package dns

import "encoding/binary"

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
