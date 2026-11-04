package keqing

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFilePreservesSameFile(t *testing.T) {
	for _, alias := range []bool{false, true} {
		t.Run(map[bool]string{false: "same-path", true: "hard-link"}[alias], func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "source")
			if err := os.WriteFile(src, []byte("original"), 0600); err != nil {
				t.Fatal(err)
			}
			dst := src
			if alias {
				dst = filepath.Join(dir, "alias")
				if err := os.Link(src, dst); err != nil {
					t.Skip(err)
				}
			}
			if err := CopyFile(src, dst); err == nil {
				t.Error("same-file copy must report an error")
			}
			data, err := os.ReadFile(src)
			if err != nil || string(data) != "original" {
				t.Fatalf("source changed: %q %v", data, err)
			}
		})
	}
}

func TestCopyFileReplacesDestination(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "source"), filepath.Join(dir, "sub", "target")
	if err := os.WriteFile(src, []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := CopyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(src, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := CopyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dst)
	if err != nil || string(data) != "x" {
		t.Fatalf("copy: %q %v", data, err)
	}
}
