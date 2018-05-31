package pngstructure

import (
    "testing"
    "path"
    "bytes"
    "fmt"

    "io/ioutil"

    "github.com/dsoprea/go-logging"
)

func TestParseFileStructure(t *testing.T) {
    filepath := path.Join(assetsPath, "Selection_058.png")

    _, err := ParseFileStructure(filepath)
    log.PanicIf(err)
}

func TestPngSplitter_Write(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintError(err)
        }
    }()

    filepath := path.Join(assetsPath, "Selection_058.png")

    original, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    ps, err := ParseBytesStructure(original)
    log.PanicIf(err)

    b := new(bytes.Buffer)

    err = ps.Write(b)
    log.PanicIf(err)

    written := b.Bytes()

    if bytes.Compare(written, original) != 0 {
        t.Fatalf("written bytes (%d) do not equal read bytes (%d)", len(written), len(original))
    }
}

func ExampleChunkSlice_Write() {
    filepath := path.Join(assetsPath, "Selection_058.png")

    cs, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    b := new(bytes.Buffer)

    err = cs.Write(b)
    log.PanicIf(err)
}

func TestChunkSlice_Write(t *testing.T) {
    chunkData := []byte {
        0x00, 0x00, 0x00, 0x0d,
        0x49, 0x48, 0x44, 0x52,
        0x00, 0x00, 0x05, 0xc0, 0x00, 0x00, 0x02, 0x56, 0x08, 0x02, 0x00, 0x00, 0x00,
        0xf0, 0x49, 0xb3, 0x65,

        0x00, 0x00, 0x00, 0x09,
        0x70, 0x48, 0x59, 0x73,
        0x00, 0x00, 0x0b, 0x13, 0x00, 0x00, 0x0b, 0x13, 0x01,
        0x00, 0x9a, 0x9c, 0x18,
    }

    b := new(bytes.Buffer)

    _, err := b.Write(PngSignature[:])
    log.PanicIf(err)

    _, err = b.Write(chunkData)
    log.PanicIf(err)

    originalFull := make([]byte, len(b.Bytes()))
    copy(originalFull, b.Bytes())

    cs, err := ParseSegments(b, len(b.Bytes()))
    log.PanicIf(err)

    chunks := cs.Chunks()
    if len(chunks) != 2 {
        t.Fatalf("number of chunks not correct")
    }

    b2 := new(bytes.Buffer)

    err = cs.Write(b2)
    log.PanicIf(err)


    actual := b2.Bytes()

    if bytes.Compare(actual, originalFull) != 0 {
        fmt.Printf("ACTUAL:\n")
        DumpBytesClause(actual)

        fmt.Printf("EXPECTED:\n")
        DumpBytesClause(originalFull)

        t.Fatalf("did not write correctly")
    }
}
