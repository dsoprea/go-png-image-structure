package pngstructure

import (
    "testing"
    "bytes"

    "github.com/dsoprea/go-logging"
)

func TestChunk_Bytes(t *testing.T) {
    c := Chunk{
        Offset: 0,
        Length: 5,
        Type: "ABCD",
        Data: []byte { 0x11, 0x22, 0x33, 0x44, 0x55 },
        Crc: 0x5678,
    }

    actual := c.Bytes()

    expected := []byte {
        0x00, 0x00, 0x00, 0x05,
        0x41, 0x42, 0x43, 0x44,
        0x11, 0x22, 0x33, 0x44, 0x55,
        0x00, 0x00, 0x56, 0x78,
    }

    if bytes.Compare(actual, expected) != 0 {
        t.Fatalf("bytes not correct")
    }
}

func ExampleChunk_Bytes() {
    c := Chunk{
        Offset: 0,
        Length: 5,
        Type: "ABCD",
        Data: []byte { 0x11, 0x22, 0x33, 0x44, 0x55 },
        Crc: 0x5678,
    }

    data := c.Bytes()
    data = data
}

func TestChunk_Write(t *testing.T) {
    c := Chunk{
        Offset: 0,
        Length: 5,
        Type: "ABCD",
        Data: []byte { 0x11, 0x22, 0x33, 0x44, 0x55 },
        Crc: 0x5678,
    }

    b := new(bytes.Buffer)
    _, err := c.Write(b)
    log.PanicIf(err)

    expected := c.Bytes()

    if bytes.Compare(b.Bytes(), expected) != 0 {
        t.Fatalf("bytes not correct")
    }
}

func ExampleChunk_Write() {
    c := Chunk{
        Offset: 0,
        Length: 5,
        Type: "ABCD",
        Data: []byte { 0x11, 0x22, 0x33, 0x44, 0x55 },
        Crc: 0x5678,
    }

    b := new(bytes.Buffer)
    _, err := c.Write(b)
    log.PanicIf(err)

    data := c.Bytes()
    data = data
}
