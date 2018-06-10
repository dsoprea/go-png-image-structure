package pngstructure

import (
	"bytes"
	"fmt"
	"path"
	"reflect"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-logging"
)

func TestChunk_Bytes(t *testing.T) {
	c := Chunk{
		Offset: 0,
		Length: 5,
		Type:   "ABCD",
		Data:   []byte{0x11, 0x22, 0x33, 0x44, 0x55},
		Crc:    0x5678,
	}

	actual := c.Bytes()

	expected := []byte{
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
		Type:   "ABCD",
		Data:   []byte{0x11, 0x22, 0x33, 0x44, 0x55},
		Crc:    0x5678,
	}

	data := c.Bytes()
	data = data

    // Output:
}

func TestChunk_Write(t *testing.T) {
	c := Chunk{
		Offset: 0,
		Length: 5,
		Type:   "ABCD",
		Data:   []byte{0x11, 0x22, 0x33, 0x44, 0x55},
		Crc:    0x5678,
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
		Type:   "ABCD",
		Data:   []byte{0x11, 0x22, 0x33, 0x44, 0x55},
		Crc:    0x5678,
	}

	b := new(bytes.Buffer)
	_, err := c.Write(b)
	log.PanicIf(err)

	data := c.Bytes()
	data = data

    // Output:
}

func TestChunkSlice_Index(t *testing.T) {
	filepath := path.Join(assetsPath, "Selection_058.png")

    pmp := NewPngMediaParser()

	cs, err := pmp.ParseFile(filepath)
	log.PanicIf(err)

	index := cs.Index()

	tallies := make(map[string]int)
	for key, chunks := range index {
		tallies[key] = len(chunks)
	}

	expected := map[string]int{
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

    pmp := NewPngMediaParser()

	cs, err := pmp.ParseFile(filepath)
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

    pmp := NewPngMediaParser()

	cs, err := pmp.ParseFile(testBasicFilepath)
	log.PanicIf(err)

	exifChunk, err := cs.FindExif()
	log.PanicIf(err)

	exifFilepath := fmt.Sprintf("%s.exif", testBasicFilepath)

	expectedExifData, err := ioutil.ReadFile(exifFilepath)
	log.PanicIf(err)

	if bytes.Compare(exifChunk.Data, expectedExifData) != 0 {
		t.Fatalf("Exif not extract correctly.")
	}
}

func TestChunkSlice_Exif(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    pmp := NewPngMediaParser()

    cs, err := pmp.ParseFile(testExifFilepath)
    log.PanicIf(err)

    rootIfd, _, err := cs.Exif()
    log.PanicIf(err)

    tags := rootIfd.Entries

    if rootIfd.Ii != exif.RootIi {
        t.Fatalf("root-IFD not parsed correctly")
    } else if len(tags) != 2 {
        t.Fatalf("incorrect number of encoded tags")
    } else if tags[0].TagId != 0x0100 {
        t.Fatalf("first tag is not correct")
    } else if tags[1].TagId != 0x0101 {
        t.Fatalf("second tag is not correct")
    }
}

func TestChunkSlice_SetExif_Existing(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()


    // Build EXIF.

    ib := exif.NewIfdBuilder(exif.RootIi, TestDefaultByteOrder)

    err := ib.AddStandardWithName("ImageWidth", []uint32{11})
    log.PanicIf(err)

    err = ib.AddStandardWithName("ImageLength", []uint32{22})
    log.PanicIf(err)


    // Replace into PNG.

    pmp := NewPngMediaParser()

    cs, err := pmp.ParseFile(testBasicFilepath)
    log.PanicIf(err)

    err = cs.SetExif(ib)
    log.PanicIf(err)

    b := new(bytes.Buffer)

    err = cs.Write(b)
    log.PanicIf(err)

    updatedImageData := b.Bytes()


    // Re-parse.

    cs, err = pmp.ParseBytes(updatedImageData)
    log.PanicIf(err)

    exifChunk, err := cs.FindExif()
    log.PanicIf(err)

    chunkData := exifChunk.Bytes()

    // Chunk data length minus length, type, and CRC data.
    expectedExifLen := len(chunkData) - 4 - 4 - 4

    if int(exifChunk.Length) != expectedExifLen {
        t.Fatalf("actual chunk data length does not match prescribed chunk data length: (%d) != (%d)", exifChunk.Length, len(exifChunk.Data))
    } else if len(exifChunk.Data) != expectedExifLen {
        t.Fatalf("chunk data length not correct")
    }

    // The first eight bytes belong to the PNG chunk structure.
    offset := 8
    _, index, err := exif.Collect(chunkData[offset : offset+expectedExifLen])
    log.PanicIf(err)

    tags := index.RootIfd.Entries

    if len(tags) != 2 {
        t.Fatalf("incorrect number of encoded tags")
    } else if tags[0].TagId != 0x0100 {
        t.Fatalf("first tag is not correct")
    } else if tags[1].TagId != 0x0101 {
        t.Fatalf("second tag is not correct")
    }
}

func TestChunkSlice_SetExif_Chunk(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
		}
	}()

	// Build EXIF.

	ib := exif.NewIfdBuilder(exif.RootIi, TestDefaultByteOrder)

	err := ib.AddStandardWithName("ImageWidth", []uint32{11})
	log.PanicIf(err)

	err = ib.AddStandardWithName("ImageLength", []uint32{22})
	log.PanicIf(err)

	// Create PNG.

	cs := NewPngChunkSlice()

	err = cs.SetExif(ib)
	log.PanicIf(err)

	exifChunk, err := cs.FindExif()
	log.PanicIf(err)

	chunkData := exifChunk.Bytes()

	// Chunk data length minus length, type, and CRC data.
	expectedExifLen := len(chunkData) - 4 - 4 - 4

	if int(exifChunk.Length) != expectedExifLen {
		t.Fatalf("actual chunk data length does not match prescribed chunk data length: (%d) != (%d)", exifChunk.Length, len(exifChunk.Data))
	} else if len(exifChunk.Data) != expectedExifLen {
		t.Fatalf("chunk data length not correct")
	}

	// The first eight bytes belong to the PNG chunk structure.
	offset := 8
	_, index, err := exif.Collect(chunkData[offset : offset+expectedExifLen])
	log.PanicIf(err)

	tags := index.RootIfd.Entries

	if len(tags) != 2 {
		t.Fatalf("incorrect number of encoded tags")
	} else if tags[0].TagId != 0x0100 {
		t.Fatalf("first tag is not correct")
	} else if tags[1].TagId != 0x0101 {
		t.Fatalf("second tag is not correct")
	}
}

func ExampleChunkSlice_SetExif() {
	// Build EXIF.

	ib := exif.NewIfdBuilder(exif.RootIi, TestDefaultByteOrder)

	err := ib.AddStandardWithName("ImageWidth", []uint32{11})
	log.PanicIf(err)

	err = ib.AddStandardWithName("ImageLength", []uint32{22})
	log.PanicIf(err)

	// Add/replace EXIF into PNG (overwrite existing).

    pmp := NewPngMediaParser()

	cs, err := pmp.ParseFile(testBasicFilepath)
	log.PanicIf(err)

	err = cs.SetExif(ib)
	log.PanicIf(err)

	b := new(bytes.Buffer)

	// Write to a `bytes.Buffer`.
	err = cs.Write(b)
	log.PanicIf(err)

    // Output:
}

func ExampleChunkSlice_Exif() {
    pmp := NewPngMediaParser()

	cs, err := pmp.ParseFile(testExifFilepath)
	log.PanicIf(err)

	rootIfd, _, err := cs.Exif()
	log.PanicIf(err)

	rootIfd = rootIfd

    // Output:
}

func ExampleChunkSlice_FindExif() {
    pmp := NewPngMediaParser()

	cs, err := pmp.ParseFile(testBasicFilepath)
	log.PanicIf(err)

	exifChunk, err := cs.FindExif()
	log.PanicIf(err)

	exifChunk = exifChunk

    // Output:
}

func ExampleChunkSlice_Index() {
	filepath := path.Join(assetsPath, "Selection_058.png")

    pmp := NewPngMediaParser()

	cs, err := pmp.ParseFile(filepath)
	log.PanicIf(err)

	index := cs.Index()
	index = index

    // Output:
}

func TestChunk_Crc32_Cycle(t *testing.T) {
    c := &Chunk{
        Type: "pHYs",
        Data: []byte { 0x00, 0x00, 0x0b, 0x13, 0x00, 0x00, 0x0b, 0x13, 0x01 },
    }

    c.UpdateCrc32()

    if c.Crc != calculateCrc32(c) {
        t.Fatalf("CRC value not consistently calculated")
    } else if c.Crc != 0x9a9c18 {
        t.Fatalf("CRC (1) not correct")
    } else if c.CheckCrc32() != true {
        t.Fatalf("CRC (1) check failed")
    }

    c.Type = "tIME"
    c.Data = []byte { 0x07, 0xcc, 0x06, 0x07, 0x11, 0x3a, 0x08 }

    c.UpdateCrc32()

    if c.Crc != 0x8eff267a {
        t.Fatalf("CRC (2) not correct")
    } else if c.CheckCrc32() != true {
        t.Fatalf("CRC (2) check failed")
    }

    c.Data = []byte { 0x99, 0x99, 0x99, 0x99 }

    if c.CheckCrc32() != false {
        t.Fatalf("CRC check didn't fail but should've")
    }
}

func TestChunkSlice_ConstructExifBuilder(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    pmp := NewPngMediaParser()

    cs, err := pmp.ParseFile(testExifFilepath)
    log.PanicIf(err)


    // Add a new tag to the additional EXIF.

    rootIb, err := cs.ConstructExifBuilder()
    log.PanicIf(err)

    err = rootIb.SetStandardWithName("ImageLength", []uint32{44})
    log.PanicIf(err)

    err = rootIb.AddStandardWithName("BitsPerSample", []uint16{33})
    log.PanicIf(err)


    // Update the image.

    err = cs.SetExif(rootIb)
    log.PanicIf(err)

    b := new(bytes.Buffer)

    err = cs.Write(b)
    log.PanicIf(err)

    updatedImageData := b.Bytes()


    // Re-parse.

    pmp = NewPngMediaParser()

    cs, err = pmp.ParseBytes(updatedImageData)
    log.PanicIf(err)

    rootIfd, _, err := cs.Exif()
    log.PanicIf(err)


    tags := rootIfd.Entries

    v1, err := rootIfd.TagValue(tags[0])
    log.PanicIf(err)

    v2, err := rootIfd.TagValue(tags[1])
    log.PanicIf(err)

    v3, err := rootIfd.TagValue(tags[2])
    log.PanicIf(err)

    if rootIfd.Ii != exif.RootIi {
        t.Fatalf("root-IFD not parsed correctly")
    } else if len(tags) != 3 {
        t.Fatalf("incorrect number of encoded tags")
    } else if tags[0].TagId != 0x0100 || reflect.DeepEqual(v1.([]uint32), []uint32 { 11 }) != true {
        t.Fatalf("first tag is not correct")
    } else if tags[1].TagId != 0x0101 || reflect.DeepEqual(v2.([]uint32), []uint32 { 44 }) != true {
        t.Fatalf("second tag is not correct")
    } else if tags[2].TagId != 0x0102 || reflect.DeepEqual(v3.([]uint16), []uint16 { 33 }) != true {
        t.Fatalf("third tag is not correct")
    }
}

func ExampleChunkSlice_ConstructExifBuilder() {
    pmp := NewPngMediaParser()

    cs, err := pmp.ParseFile(testExifFilepath)
    log.PanicIf(err)


    // Add a new tag to the additional EXIF.

    rootIb, err := cs.ConstructExifBuilder()
    log.PanicIf(err)

    err = rootIb.SetStandardWithName("ImageLength", []uint32{44})
    log.PanicIf(err)

    err = rootIb.AddStandardWithName("BitsPerSample", []uint16{33})
    log.PanicIf(err)


    // Update the image.

    err = cs.SetExif(rootIb)
    log.PanicIf(err)

    b := new(bytes.Buffer)

    err = cs.Write(b)
    log.PanicIf(err)

    updatedImageData := b.Bytes()


    // Re-parse.

    pmp = NewPngMediaParser()

    cs, err = pmp.ParseBytes(updatedImageData)
    log.PanicIf(err)

    rootIfd, _, err := cs.Exif()
    log.PanicIf(err)


    for i, ite := range rootIfd.Entries {
        value, err := rootIfd.TagValue(ite)
        log.PanicIf(err)

        fmt.Printf("%d: (0x%04x) %v\n", i, ite.TagId, value)
    }

    // Output:
    // 0: (0x0100) [11]
    // 1: (0x0101) [44]
    // 2: (0x0102) [33]
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

    pmp := NewPngMediaParser()

    ps, err := pmp.ParseBytes(original)
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

    pmp := NewPngMediaParser()

    cs, err := pmp.ParseFile(filepath)
    log.PanicIf(err)

    b := new(bytes.Buffer)

    err = cs.Write(b)
    log.PanicIf(err)

    // Output:
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

    pmp := NewPngMediaParser()

    cs, err := pmp.Parse(b, len(b.Bytes()))
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
