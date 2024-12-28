package dns

import "bytes"

const headerOffset = 12

type packet struct {
	questions []question
	answers   []answerA
	header    header
}

type labels []string

func NewPacket() *packet {
	return &packet{}
}

func (p *packet) Parse(b []byte) {
	p.header.Parse(b[:headerOffset])
	p.questions = make([]question, p.header.QDCOUNT)
	// not prepared for multiple questions, offset needs to be adjusted
	for i := 0; i < int(p.header.QDCOUNT); i++ {
		q := NewQuestion()
		q.Parse(b[headerOffset:], p.header.QDCOUNT)
		p.questions[i] = *q
	}
}

func (p *packet) Resolve() {
	// this is a stub for now before implementing multiple answers
	p.header.ANCOUNT = 1
	a := NewAnswerA()
	a.Find(p.questions[0].Labels)
	p.answers = append(p.answers, *a)
}

func (p *packet) Bytes() ([]byte, error) {
	b := make([]byte, 0)
	headerBytes := p.header.Bytes()
	b = append(b, headerBytes...)
	for i := 0; i < int(p.header.QDCOUNT); i++ {
		questionBytes, err := p.questions[i].Bytes()
		if err != nil {
			return nil, err
		}
		b = append(b, questionBytes...)
	}
	for i := 0; i < int(p.header.ANCOUNT); i++ {
		answerBytes, err := p.answers[i].Bytes()
		if err != nil {
			return nil, err
		}
		b = append(b, answerBytes...)
	}

	return b, nil
}

func (l labels) bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, label := range l {
		// write the length of the label
		if err := buf.WriteByte(uint8(len(label))); err != nil {
			return nil, err
		}
		if _, err := buf.Write([]byte(label)); err != nil {
			return nil, err
		}
	}
	// null byte to indicate that the labels section is over
	if err := buf.WriteByte(uint8(0)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
