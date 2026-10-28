package rpc

import "testing"

func TestRoundtripPlain(t *testing.T) {
	payload := []byte(`{"m":"session/ping"}`)
	frame, err := Encode(payload, false)
	if err != nil {
		t.Fatal(err)
	}
	length, compressed, err := DecodeHeader(frame[:4])
	if err != nil {
		t.Fatal(err)
	}
	if compressed {
		t.Fatal("unexpected compressed flag")
	}
	if int(length) != len(payload) {
		t.Fatalf("length %d, want %d", length, len(payload))
	}
	if string(frame[4:]) != string(payload) {
		t.Fatal("payload corrupted")
	}
}

func TestCompressedFlagSurvives(t *testing.T) {
	frame, err := Encode([]byte("abc"), true)
	if err != nil {
		t.Fatal(err)
	}
	length, compressed, err := DecodeHeader(frame[:4])
	if err != nil {
		t.Fatal(err)
	}
	if !compressed || length != 3 {
		t.Fatalf("got len=%d compressed=%v", length, compressed)
	}
}

func TestRejectsOversized(t *testing.T) {
	if _, err := Encode(make([]byte, MaxFrame+1), false); err == nil {
		t.Fatal("expected error for oversized frame")
	}
}

func TestRejectsBadHeader(t *testing.T) {
	if _, _, err := DecodeHeader([]byte{1, 2}); err == nil {
		t.Fatal("expected error for short header")
	}
}
