package jsluice

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// Every file in the repository goes into the Go module zip that each
// consumer downloads, so build output must not be committed.
func TestNoCompiledBinaries(t *testing.T) {
	magics := [][]byte{
		{0x7f, 'E', 'L', 'F'},    // ELF
		{0xcf, 0xfa, 0xed, 0xfe}, // Mach-O 64-bit
		{0xfe, 0xed, 0xfa, 0xcf},
		{'M', 'Z'}, // PE
	}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		head := make([]byte, 4)
		n, _ := f.Read(head)
		_ = f.Close()
		for _, m := range magics {
			if n >= len(m) && bytes.Equal(head[:len(m)], m) {
				t.Errorf("%s is a compiled binary", path)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
