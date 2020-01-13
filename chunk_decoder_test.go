package pngstructure

import (
    "path"
    "testing"

    "github.com/dsoprea/go-logging"
)

func TestChunkDecoder_decodeIHDR(t *testing.T) {
    assetsPath := getTestAssetsPath()
    filepath := path.Join(assetsPath, "Selection_058.png")

    pmp := NewPngMediaParser()

    cs, err := pmp.ParseFile(filepath)
    log.PanicIf(err)

    index := cs.Index()
    ihdrRawSlice, found := index["IHDR"]

    if found != true {
        t.Fatalf("Could not find IHDR chunk.")
    }

    cd := NewChunkDecoder()

    ihdrRaw, err := cd.Decode(ihdrRawSlice[0])
    log.PanicIf(err)

    ihdr := ihdrRaw.(*ChunkIHDR)

    expected := &ChunkIHDR{
        Width:             1472,
        Height:            598,
        BitDepth:          8,
        ColorType:         2,
        CompressionMethod: 0,
        FilterMethod:      0,
        InterlaceMethod:   0,
    }

    if *ihdr != *expected {
        t.Fatalf("ihdr not correct")
    }
}

func ExampleChunkDecoder_Decode() {
    filepath := path.Join(assetsPath, "Selection_058.png")

    pmp := NewPngMediaParser()

    cs, err := pmp.ParseFile(filepath)
    log.PanicIf(err)

    index := cs.Index()
    ihdrRawSlice, found := index["IHDR"]

    if found != true {
        log.Panicf("IHDR chunk not found")
    }

    cd := NewChunkDecoder()

    ihdrRaw, err := cd.Decode(ihdrRawSlice[0])
    log.PanicIf(err)

    ihdr := ihdrRaw.(*ChunkIHDR)
    ihdr = ihdr
}
