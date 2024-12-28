package dns

import "testing"

var headerBytes = []byte{
	0b00010011, 0b10111110, 0b00000001, 0b00100000,
	0b00000000, 0b00000001, 0b00000000, 0b00000000,
	0b00000000, 0b00000000, 0b00000000, 0b00000000,
}

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
