package pngstructure

import (
    "errors"
    "io"
    "bytes"
    "fmt"

    "encoding/binary"

    "github.com/dsoprea/go-logging"
)

var (
    PngSignature = [8]byte { 137, 'P', 'N', 'G', '\r', '\n', 26, '\n' }
)

var (
    ErrNotPng = errors.New("not png data")
)


// ChunkSlice encapsulates a slice of chunks.
type ChunkSlice struct {
    chunks []*Chunk
}

func NewChunkSlice(chunks []*Chunk) *ChunkSlice {
    return &ChunkSlice{
        chunks: chunks,
    }
}

func (cs *ChunkSlice) String() string {
    return fmt.Sprintf("ChunkSlize<LEN=(%d)>", len(cs.chunks))
}

// Chunks exposes the actual slice.
func (cs *ChunkSlice) Chunks() []*Chunk {
    return cs.chunks
}

// Write encodes and writes all chunks.
func (cs *ChunkSlice) Write(w io.Writer) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    _, err = w.Write(PngSignature[:])
    log.PanicIf(err)

    for _, c := range cs.chunks {
        _, err := c.Write(w)
        log.PanicIf(err)
    }

    return nil
}

// Index returns a map of chunk types to chunk slices, grouping all like chunks.
func (cs *ChunkSlice) Index() (index map[string][]*Chunk) {
    index = make(map[string][]*Chunk)
    for _, c := range cs.chunks {
        if grouped, found := index[c.Type]; found == true {
            index[c.Type] = append(grouped, c)
        } else {
            index[c.Type] = []*Chunk { c }
        }
    }

    return index
}


// PngSplitter hosts the princpal `Split()` method uses by `bufio.Scanner`.
type PngSplitter struct {
    chunks []*Chunk
    currentOffset int
}

func (ps *PngSplitter) Chunks() *ChunkSlice {
    return NewChunkSlice(ps.chunks)
}

func NewPngSplitter() *PngSplitter {
    return &PngSplitter{
        chunks: make([]*Chunk, 0),
    }
}


// Chunk describes a single chunk.
type Chunk struct {
    Offset int
    Length uint32
    Type string
    Data []byte
    Crc uint32
}

func (c Chunk) String() string {
    return fmt.Sprintf("Chunk<OFFSET=(%d) LENGTH=(%d) TYPE=[%s] CRC=(%d)>", c.Offset, c.Length, c.Type, c.Crc)
}

// Bytes encodes and returns the bytes for this chunk.
func (c Chunk) Bytes() []byte {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.Panic(err)
        }
    }()

    if len(c.Data) != int(c.Length) {
        log.Panicf("length of data not correct")
    }

    preallocated := make([]byte, 0, 4 + 4 + c.Length + 4)
    b := bytes.NewBuffer(preallocated)

    err := binary.Write(b, binary.BigEndian, c.Length)
    log.PanicIf(err)

    _, err = b.Write([]byte(c.Type))
    log.PanicIf(err)

    _, err = b.Write(c.Data)
    log.PanicIf(err)

    err = binary.Write(b, binary.BigEndian, c.Crc)
    log.PanicIf(err)

    return b.Bytes()
}

// Write encodes and writes the bytes for this chunk.
func (c Chunk) Write(w io.Writer) (count int, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if len(c.Data) != int(c.Length) {
        log.Panicf("length of data not correct")
    }

    err = binary.Write(w, binary.BigEndian, c.Length)
    log.PanicIf(err)

    _, err = w.Write([]byte(c.Type))
    log.PanicIf(err)

    _, err = w.Write(c.Data)
    log.PanicIf(err)

    err = binary.Write(w, binary.BigEndian, c.Crc)
    log.PanicIf(err)

    return 4 + len(c.Type) + len(c.Data) + 4, nil
}

// readHeader verifies that the PNG header bytes appear next.
func (ps *PngSplitter) readHeader(r io.Reader) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    len_ := len(PngSignature)
    header := make([]byte, len_)

    _, err = r.Read(header)
    log.PanicIf(err)

    ps.currentOffset += len_

    if bytes.Compare(header, PngSignature[:]) != 0 {
        log.Panic(ErrNotPng)
    }

    return nil
}

// Split fulfills the `bufio.SplitFunc` function definition for
// `bufio.Scanner`.
func (ps *PngSplitter) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    // We might have more than one chunk's worth, and, if `atEOF` is true, we
    // won't be called again. We'll repeatedly try to read additional chunks,
    // but, when we run out of the data we were given then we'll return the
    // number of bytes fo rthe chunks we've already completely read. Then,
    // we'll be called again from theend ofthose bytes, at which point we'll
    // indicate that we don't yet have enough for another chunk, and we should
    // be then called with more.
    for {
        len_ := len(data)
        if len_ < 8 {
            return advance, nil, nil
        }

        length := binary.BigEndian.Uint32(data[:4])
        type_ := string(data[4:8])
        chunkSize := (8 + int(length) + 4)

        if len_ < chunkSize {
            return advance, nil, nil
        }

        crcIndex := 8 + length
        crc := binary.BigEndian.Uint32(data[crcIndex:crcIndex + 4])

        content := make([]byte, length)
        copy(content, data[8:8+length])

        c := &Chunk{
            Length: length,
            Type: type_,
            Data: content,
            Crc: crc,
            Offset: ps.currentOffset,
        }

        ps.chunks = append(ps.chunks, c)

        advance += chunkSize
        ps.currentOffset += chunkSize

        data = data[chunkSize:]
    }

    return advance, nil, nil
}
