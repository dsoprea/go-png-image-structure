package pngstructure

import (
    "os"
    "path"

    "github.com/dsoprea/go-logging"
)

var (
    assetsPath = ""
)

func init() {
    goPath := os.Getenv("GOPATH")
    if goPath == "" {
        log.Panicf("GOPATH is empty")
    }

    assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-png-image-structure", "assets")
}
