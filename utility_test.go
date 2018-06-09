package pngstructure

import (
    "path"
    "testing"
    "fmt"

    "io/ioutil"

    "github.com/dsoprea/go-logging"
)

func TestIsPng(t *testing.T) {
    filepath := path.Join(assetsPath, "libpng.png")

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    if IsPng(data) != true {
        t.Fatalf("not detected as png")
    }
}

func ExampleIsPng() {
    filepath := path.Join(assetsPath, "libpng.png")

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    isPng := IsPng(data)
    fmt.Printf("%v\n", isPng)

    // Output:
    // true
}
