package pngstructure

import (
    "testing"
    "bytes"
    "path"
    "reflect"
    "fmt"

    "io/ioutil"

    "github.com/dsoprea/go-logging"
    // "github.com/dsoprea/go-exif"
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

func TestChunkSlice_Index(t *testing.T) {
    filepath := path.Join(assetsPath, "Selection_058.png")

    cs, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    index := cs.Index()

    tallies := make(map[string]int)
    for key, chunks := range index {
        tallies[key] = len(chunks)
    }

    expected := map[string]int {
        "IDAT": 222,
        "IEND": 1,
        "IHDR": 1,
        "pHYs": 1,
        "tIME": 1,
    }

    if reflect.DeepEqual(tallies, expected) != true {
        t.Fatalf("index not correct")
    }
}

func TestChunkSlice_FindExif_Miss(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    filepath := path.Join(assetsPath, "Selection_058.png")

    cs, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    _, err = cs.FindExif()

    if err == nil {
        t.Fatalf("expected error for missing EXIF")
    } else if log.Is(err, ErrNoExif) == false {
        log.Panic(err)
    }
}

func TestChunkSlice_FindExif_Hit(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    filepath := path.Join(assetsPath, "pngexif.png")

    cs, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    exifChunk, err := cs.FindExif()
    log.PanicIf(err)

    exifFilepath := fmt.Sprintf("%s.exif", filepath)

    expectedExifData, err := ioutil.ReadFile(exifFilepath)
    log.PanicIf(err)

    if bytes.Compare(exifChunk.Data, expectedExifData) != 0 {
        t.Fatalf("Exif not extract correctly.")
    }
}

// TODO(dustin): !! The test-file from the libpng project is actually broken (no next-IFD uint32 at the bottom of the EXIF IFD).
// func TestChunkSlice_Exif(t *testing.T) {
//     defer func() {
//         if state := recover(); state != nil {
//             err := log.Wrap(state.(error))
//             log.PrintErrorf(err, "Test failure.")
//         }
//     }()

//     filepath := path.Join(assetsPath, "pngexif.png")

//     cs, err := ParseFileStructure(filepath)
//     log.PanicIf(err)

//     rootIfd, err := cs.Exif()
//     log.PanicIf(err)

//     if rootIfd.Ii != exif.RootIi {
//         t.Fatalf("root-IFD not parsed correctly")
//     }
// }

func ExampleChunkSlice_Exif() {
    filepath := path.Join(assetsPath, "pngexif.png")

    cs, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    rootIfd, err := cs.Exif()
    log.PanicIf(err)

    rootIfd = rootIfd
}

func ExampleChunkSlice_FindExif() {
    filepath := path.Join(assetsPath, "pngexif.png")

    cs, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    exifChunk, err := cs.FindExif()
    log.PanicIf(err)

    exifChunk = exifChunk
}

func ExampleChunkSlice_Index() {
    filepath := path.Join(assetsPath, "Selection_058.png")

    cs, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    index := cs.Index()
    index = index
}
