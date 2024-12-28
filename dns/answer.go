package dns

import (
	"bytes"
	"encoding/binary"
)

type answerA struct {
	Name   labels
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   uint32
}

func NewAnswerA() *answerA {
	return &answerA{
		Name:   labels{},
		Type:   1,
		Class:  1,
		TTL:    0,
		Length: 0,
		Data:   0,
	}
}

func (a *answerA) Find(name labels) {
	a.Name = name
	a.TTL = 60
	a.Length = 4 // 4 bytes for an IPv4 address
	a.Data = 0x08080808
}

func (a answerA) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	lb, err := a.Name.bytes()
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(lb); err != nil {
		return nil, err
	}
	// write the type and class
	if err := binary.Write(buf, binary.BigEndian, a.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, a.Class); err != nil {
		return nil, err
	}
	// TTL
	if err := binary.Write(buf, binary.BigEndian, a.TTL); err != nil {
		return nil, err
	}
	// Data
	if err := binary.Write(buf, binary.BigEndian, a.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
