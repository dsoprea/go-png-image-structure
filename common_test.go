package pngstructure

import (
    "os"
    "path"

    "encoding/binary"

    "github.com/dsoprea/go-logging"
)

var (
    assetsPath = ""
    testBasicFilepath = ""
    testExifFilepath = ""

    TestDefaultByteOrder = binary.BigEndian
)

func init() {
    goPath := os.Getenv("GOPATH")
    if goPath == "" {
        log.Panicf("GOPATH is empty")
    }

    assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-png-image-structure", "assets")

    testBasicFilepath = path.Join(assetsPath, "libpng.png")
    testExifFilepath = path.Join(assetsPath, "exif.png")
}
