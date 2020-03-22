package pngstructure

import (
    "bufio"
    "bytes"
    "io"
    "os"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-utility/image"
)

type PngMediaParser struct {
}

func NewPngMediaParser() *PngMediaParser {
    return new(PngMediaParser)
}

func (pmp *PngMediaParser) Parse(r io.Reader, size int) (ec riimage.MediaContext, err error) {
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
    buffer := []byte{}
    s.Buffer(buffer, size)
    s.Split(ps.Split)

    for s.Scan() != false {
    }
    log.PanicIf(s.Err())

    return ps.Chunks(), nil
}

func (pmp *PngMediaParser) ParseFile(filepath string) (ec riimage.MediaContext, err error) {
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

    chunks, err := pmp.Parse(f, int(size))
    log.PanicIf(err)

    return chunks, nil
}

func (pmp *PngMediaParser) ParseBytes(data []byte) (ec riimage.MediaContext, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    b := bytes.NewBuffer(data)

    chunks, err := pmp.Parse(b, len(data))
    log.PanicIf(err)

    return chunks, nil
}

func (pmp *PngMediaParser) LooksLikeFormat(data []byte) bool {
    return bytes.Compare(data[:len(PngSignature)], PngSignature[:]) == 0
}
