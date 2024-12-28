package main

import "testing"

var (
	headerBytes = []byte{
		0b00010011, 0b10111110, 0b00000001, 0b00100000,
		0b00000000, 0b00000001, 0b00000000, 0b00000000,
		0b00000000, 0b00000000, 0b00000000, 0b00000000,
	}
	questionBytes = []byte{
		0b00000110, 0b01100111, 0b01101111, 0b01101111,
		0b01100111, 0b01101100, 0b01100101, 0b00000011,
		0b01100011, 0b01101111, 0b01101101, 0b00000000,
		0b00000000, 0b00000001, 0b00000000, 0b00000001,
	}
)

func TestHeaderParse(t *testing.T) {
	h := NewHeader()
	h.Parse(headerBytes)

	if h.ID != 5054 {
		t.Errorf("expected ID to be 5054, instead, got %d", h.ID)
	}
	if h.QR != 1 {
		t.Errorf("expected QR to be 1, instead, got %d", h.QR)
	}
	if h.QDCOUNT != 1 {
		t.Errorf("expected QDCOUNT to be 1, instead, got %d", h.QDCOUNT)
	}
}

func TestQuestionParse(t *testing.T) {
	q := NewQuestion()
	const qdcount = 1
	q.Parse(questionBytes, qdcount)

	if len(q.Labels) != 2 {
		t.Errorf("expected 2 labels, instead, got %d", len(q.Labels))
	}

	if q.Labels[0] != "google" {
		t.Errorf("expected Labels[0] to be 'google', instead, got %s", q.Labels[0])
	}

	if q.Labels[1] != "com" {
		t.Errorf("expected Labels[1] to be 'com', instead, got %s", q.Labels[1])
	}

	if q.Type != 1 {
		t.Errorf("expected Type to be 1, instead, got %d", q.Type)
	}

	if q.Class != 1 {
		t.Errorf("expected Class to be 1, isntead, got %d", q.Type)
	}
}

func TestQuestionBytes(t *testing.T) {
	q := NewQuestion()
	const qdcount = 1
	q.Parse(questionBytes, qdcount)
	b, err := q.Bytes()
	if err != nil {
		t.Error(err)
	}
	if !compareSlices(questionBytes, b) {
		t.Errorf("generated question byte slice is different from source: \nsource: \n% x\ngenerated: \n% x", questionBytes, b)
	}
}

func compareSlices[K comparable](sl1, sl2 []K) bool {
	if len(sl1) != len(sl2) {
		return false
	}
	for i := 0; i < len(sl1); i++ {
		if sl1[i] != sl2[i] {
			return false
		}
	}
	return true
}
