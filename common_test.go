package pngstructure

import (
    "path"

    "go/build"

    "github.com/dsoprea/go-logging"
)

var (
    assetsPath = ""
)

// getModuleRootPath returns our source-path when running from source during
// tests.
func getModuleRootPath() string {
    p, err := build.Default.Import(
        "github.com/dsoprea/go-png-image-structure",
        build.Default.GOPATH,
        build.FindOnly)

    log.PanicIf(err)

    packagePath := p.Dir
    return packagePath
}

func getTestAssetsPath() string {
    moduleRootPath := getModuleRootPath()
    assetsPath := path.Join(moduleRootPath, "assets")

    return assetsPath
}

func getTestBasicImageFilepath() string {
    return path.Join(assetsPath, "libpng.png")
}

func getTestExifImageFilepath() string {
    return path.Join(assetsPath, "exif.png")
}

func init() {
    assetsPath = getTestAssetsPath()
}
