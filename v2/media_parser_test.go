package pngstructure

import (
	"fmt"
	"path"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

func TestPngMediaParser_ParseFile(t *testing.T) {
	filepath := path.Join(assetsPath, "Selection_058.png")

	pmp := NewPngMediaParser()

	_, err := pmp.ParseFile(filepath)
	log.PanicIf(err)
}

func TestPngMediaParser_LooksLikeFormat(t *testing.T) {
	filepath := path.Join(assetsPath, "libpng.png")

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	pmp := NewPngMediaParser()

	if pmp.LooksLikeFormat(data) != true {
		t.Fatalf("not detected as png")
	}
}

func ExamplePngMediaParser_LooksLikeFormat() {
	filepath := path.Join(assetsPath, "libpng.png")

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	pmp := NewPngMediaParser()

	isPng := pmp.LooksLikeFormat(data)
	fmt.Printf("%v\n", isPng)

	// Output:
	// true
}
