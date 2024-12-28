package dns

import (
	"bytes"
	"encoding/binary"
)

type question struct {
	Labels labels
	Type   uint16
	Class  uint16
}

func NewQuestion() *question {
	return &question{
		Labels: labels{},
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
	q.Type = binary.BigEndian.Uint16(data[offset : offset+2])
	// increment the offset to skip the type
	offset += 2
	// get the encoded class
	q.Class = binary.BigEndian.Uint16(data[offset : offset+2])
}

func (q question) Bytes() ([]byte, error) {
	// create an empty buffer to write the question
	buf := new(bytes.Buffer)

	lb, err := q.Labels.bytes()
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(lb); err != nil {
		return nil, err
	}
	// write the type and class
	if err := binary.Write(buf, binary.BigEndian, q.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, q.Class); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
