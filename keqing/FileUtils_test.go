package keqing

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileExistsInvalidPath(t *testing.T) {
	if FileExists("invalid\x00path") {
		t.Fatal("invalid path reported as existing")
	}
}

func TestFileExistsWithError(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	if err := os.WriteFile(file, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{dir, file} {
		if exists, err := FileExistsWithError(path); !exists || err != nil {
			t.Fatalf("existing %s: %v %v", path, exists, err)
		}
	}
	if exists, err := FileExistsWithError(filepath.Join(dir, "missing")); exists || err != nil {
		t.Fatalf("missing: %v %v", exists, err)
	}
	if exists, err := FileExistsWithError("invalid\x00path"); exists || err == nil || errors.Is(err, os.ErrNotExist) {
		t.Fatalf("invalid: %v %v", exists, err)
	}
}

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
