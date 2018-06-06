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

// TODO(dustin): !! Write test for ConstructExifBuilder

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

// TODO(dustin): !! Add example to update existing.

func ExampleChunkSlice_SetExif() {
	// Build EXIF.

	ib := exif.NewIfdBuilder(exif.RootIi, TestDefaultByteOrder)

	err := ib.AddStandardWithName("ImageWidth", []uint32{11})
	log.PanicIf(err)

	err = ib.AddStandardWithName("ImageLength", []uint32{22})
	log.PanicIf(err)

	// Add/replace EXIF into PNG (overwrite existing).

	filepath := path.Join(assetsPath, "pngexif.png")

	cs, err := ParseFileStructure(filepath)
	log.PanicIf(err)

	err = cs.SetExif(ib)
	log.PanicIf(err)

	b := new(bytes.Buffer)

	// Write to a `bytes.Buffer`.
	err = cs.Write(b)
	log.PanicIf(err)
}

func ExampleChunkSlice_Exif() {
	filepath := path.Join(assetsPath, "pngexif.png")

	cs, err := ParseFileStructure(filepath)
	log.PanicIf(err)

	_, rootIfd, err := cs.Exif()
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
