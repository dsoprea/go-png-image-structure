package pngstructure

import (
    "os"
    "io"
    "bufio"
    "bytes"

    "github.com/dsoprea/go-logging"
)

func ParseSegments(r io.Reader, size int) (chunks *ChunkSlice, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    ps := NewPngSplitter()

    err = ps.readHeader(r)
    log.PanicIf(err)

    s := bufio.NewScanner(r)

    // Since each segment can be any size, our buffer must be allowed to grow
    // as large as the file.
    buffer := []byte {}
    s.Buffer(buffer, size)
    s.Split(ps.Split)

    for ; s.Scan() != false; { }
    log.PanicIf(s.Err())

    return ps.Chunks(), nil
}

func ParseFileStructure(filepath string) (chunks *ChunkSlice, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    f, err := os.Open(filepath)
    log.PanicIf(err)

    stat, err := f.Stat()
    log.PanicIf(err)

    size := stat.Size()

    chunks, err = ParseSegments(f, int(size))
    log.PanicIf(err)

    return chunks, nil
}

func ParseBytesStructure(data []byte) (chunks *ChunkSlice, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    b := bytes.NewBuffer(data)

    chunks, err = ParseSegments(b, len(data))
    log.PanicIf(err)

    return chunks, nil
}
