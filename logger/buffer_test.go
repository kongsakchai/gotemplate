package logger

import "testing"

func TestNewBuffer(t *testing.T) {
	b := newBuffer()
	if b == nil {
		t.Fatal("newBuffer returned nil")
	}
	if cap(*b) != 1024 {
		t.Fatalf("expected buffer capacity of 1024, got %d", cap(*b))
	}
}

func TestBufferWrite(t *testing.T) {
	b := newBuffer()
	defer b.Free()

	data := []byte("test data")
	n, err := b.Write(data)
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n != len(data) {
		t.Fatalf("expected Write to return %d, got %d", len(data), n)
	}

	err = b.WriteByte('!')
	if err != nil {
		t.Fatalf("WriteByte returned error: %v", err)
	}

	dataStr := " more data"
	n, err = b.WriteString(dataStr)
	if err != nil {
		t.Fatalf("WriteString returned error: %v", err)
	}
	if n != len(dataStr) {
		t.Fatalf("expected WriteString to return %d, got %d", len(" more data"), n)
	}

	if string(b.Bytes()) != "test data! more data" {
		t.Fatalf("expected buffer content 'test data', got '%s'", b.Bytes())
	}
}

func TestBufferReset(t *testing.T) {
	b := newBuffer()

	data := []byte("test data")
	_, err := b.Write(data)
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	b.Free()
	if len(*b) != 0 {
		t.Fatalf("expected buffer length to be 0 after Reset, got %d", len(*b))
	}
	if cap(*b) != 1024 {
		t.Fatalf("expected buffer capacity to remain 1024 after Reset, got %d", cap(*b))
	}
}
