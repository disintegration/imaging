package imaging_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// ConvertToJXL converts an image file supported by imaging to JPEG XL.
//
// imaging itself has no JPEG XL encoder (neither does the Go standard
// library or golang.org/x/image), so conversion uses the first available
// external libjxl backend: cjxl (the only one that can also transcode
// JPEG losslessly, bit-exact), then ffmpeg with libjxl, then ImageMagick.
// The output is only reported as success if it starts with a JPEG XL
// codestream (FF 0A) or container signature.
func ConvertToJXL(src, dst string) error {
	if _, err := imaging.Open(src); err != nil { // validate decodable input
		return err
	}
	if filepath.Ext(dst) != ".jxl" {
		return fmt.Errorf("output must end in .jxl: %s", dst)
	}
	attempts := [][]string{
		{"cjxl", src, dst, "--effort=7", "--distance=1.0"},
		{"ffmpeg", "-y", "-i", src, "-c:v", "libjxl", "-effort", "7", "-distance", "1.0", dst},
		{"magick", src, dst},
	}
	for _, cmd := range attempts {
		if _, err := exec.LookPath(cmd[0]); err != nil {
			continue
		}
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			os.Remove(dst)
			continue
		}
		head, err := os.ReadFile(dst)
		if err == nil && (bytes.HasPrefix(head, []byte{0xFF, 0x0A}) ||
			bytes.HasPrefix(head, []byte{0, 0, 0, 0x0C, 'J', 'X', 'L', ' ', 0x0D, 0x0A, 0x87, 0x0A})) {
			return nil
		}
		os.Remove(dst)
	}
	return fmt.Errorf("no working JPEG XL backend (install libjxl cjxl, ffmpeg --enable-libjxl, or ImageMagick with jxl delegate)")
}

func ExampleConvertToJXL() {
	src := filepath.Join("testdata", "lena_512.png")
	dst := filepath.Join(os.TempDir(), "imaging-example.jxl")
	defer os.Remove(dst)
	// No Output comment: this example is compile-checked only, because it
	// converts only when a libjxl backend (cjxl/ffmpeg/magick) is installed.
	_ = ConvertToJXL(src, dst)
}
