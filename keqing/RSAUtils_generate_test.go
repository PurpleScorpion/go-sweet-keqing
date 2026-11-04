package keqing

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func rsaGeneratePanic(path string) (value any) {
	defer func() { value = recover() }()
	RsaGenerateKey(path)
	return nil
}

func TestRsaGenerateKeyWithErrorAndRepeat(t *testing.T) {
	dir := t.TempDir()
	if got, err := RsaGenerateKeyWithError(dir); err != nil || got != dir {
		t.Fatalf("generate = %q, %v", got, err)
	}
	before := make(map[string]string)
	for _, name := range []string{"publicKey.pem", "privateKey.pem"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		before[name] = string(data)
	}
	if got, err := RsaGenerateKeyWithError(dir); got != "" || !errors.Is(err, os.ErrExist) {
		t.Fatalf("repeat = %q, %v", got, err)
	}
	for name, want := range before {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || string(data) != want {
			t.Errorf("existing %s changed: %v", name, err)
		}
	}
}

func TestRsaGenerateKeyRejectsExistingSymlinks(t *testing.T) {
	for _, name := range []string{"publicKey.pem", "privateKey.pem"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, "original")
			if err := os.WriteFile(target, []byte("original"), 0600); err != nil {
				t.Fatal(err)
			}
			link := filepath.Join(dir, name)
			if err := os.Symlink(target, link); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}
			if _, err := RsaGenerateKeyWithError(dir); !errors.Is(err, os.ErrExist) {
				t.Fatalf("expected exist error, got %v", err)
			}
			data, err := os.ReadFile(target)
			if err != nil || string(data) != "original" {
				t.Fatalf("symlink target changed: %v", err)
			}
			if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
				t.Fatalf("symlink changed: %v", err)
			}
		})
	}
}

type rsaFaultFile struct {
	*os.File
	writeErr, syncErr, closeErr error
	short                       bool
	closed                      *int
	afterClose                  func()
}

func (f *rsaFaultFile) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	if f.short {
		return len(data) - 1, nil
	}
	return f.File.Write(data)
}

func (f *rsaFaultFile) Sync() error {
	if f.syncErr != nil {
		return f.syncErr
	}
	return f.File.Sync()
}

func (f *rsaFaultFile) Close() error {
	err := f.File.Close()
	if f.closed != nil {
		*f.closed++
	}
	if f.afterClose != nil {
		f.afterClose()
	}
	return errors.Join(err, f.closeErr)
}

func TestRsaKeyFileErrorsAreReportedAndCleaned(t *testing.T) {
	for _, name := range []string{"privateKey.pem", "publicKey.pem"} {
		for _, operation := range []string{"write", "short write", "sync", "close"} {
			t.Run(name+"/"+operation, func(t *testing.T) {
				dir := t.TempDir()
				want := errors.New("injected " + operation)
				if operation == "short write" {
					want = io.ErrShortWrite
				}
				closed := 0
				opener := func(path string, flag int, mode os.FileMode) (rsaKeyFile, error) {
					file, err := os.OpenFile(path, flag, mode)
					if err != nil {
						return nil, err
					}
					fault := &rsaFaultFile{File: file, closed: &closed}
					if filepath.Base(path) == name {
						switch operation {
						case "write":
							fault.writeErr = want
						case "short write":
							fault.short = true
						case "sync":
							fault.syncErr = want
						case "close":
							fault.closeErr = want
						}
					}
					return fault, nil
				}
				err := createRSAKeyFiles(dir, []byte("private"), []byte("public"), opener)
				if !errors.Is(err, want) {
					t.Fatalf("error = %v; want %v", err, want)
				}
				if closed != 2 {
					t.Fatalf("closed %d handles, want 2", closed)
				}
				for _, file := range []string{"privateKey.pem", "publicKey.pem"} {
					if _, err := os.Lstat(filepath.Join(dir, file)); !os.IsNotExist(err) {
						t.Errorf("incomplete %s remains: %v", file, err)
					}
				}
			})
		}
	}
}

func TestRsaKeyFileKeepsWriteAndCloseErrors(t *testing.T) {
	dir := t.TempDir()
	writeErr, privateCloseErr, publicCloseErr := errors.New("write"), errors.New("private close"), errors.New("public close")
	opener := func(path string, flag int, mode os.FileMode) (rsaKeyFile, error) {
		file, err := os.OpenFile(path, flag, mode)
		if err != nil {
			return nil, err
		}
		fault := &rsaFaultFile{File: file, closeErr: publicCloseErr}
		if filepath.Base(path) == "privateKey.pem" {
			fault.writeErr = writeErr
			fault.closeErr = privateCloseErr
		}
		return fault, nil
	}
	err := createRSAKeyFiles(dir, []byte("private"), []byte("public"), opener)
	for _, want := range []error{writeErr, privateCloseErr, publicCloseErr} {
		if !errors.Is(err, want) {
			t.Errorf("missing error %v from %v", want, err)
		}
	}
}

func TestRsaKeyFileCleanupPreservesReplacement(t *testing.T) {
	dir := t.TempDir()
	want := errors.New("write failed")
	opener := func(path string, flag int, mode os.FileMode) (rsaKeyFile, error) {
		file, err := os.OpenFile(path, flag, mode)
		if err != nil {
			return nil, err
		}
		fault := &rsaFaultFile{File: file}
		if filepath.Base(path) == "privateKey.pem" {
			fault.writeErr = want
			fault.afterClose = func() {
				if err := os.Rename(path, path+".original"); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("replacement"), 0600); err != nil {
					t.Fatal(err)
				}
			}
		}
		return fault, nil
	}
	err := createRSAKeyFiles(dir, []byte("private"), []byte("public"), opener)
	if !errors.Is(err, want) {
		t.Fatalf("error = %v", err)
	}
	data, readErr := os.ReadFile(filepath.Join(dir, "privateKey.pem"))
	if readErr != nil || string(data) != "replacement" {
		t.Fatalf("replacement removed: %v", readErr)
	}
}

func TestRsaKeyFileConcurrentCreation(t *testing.T) {
	dir := t.TempDir()
	var wg sync.WaitGroup
	results := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- createRSAKeyFiles(dir, []byte("private"), []byte("public"), openRSAKeyFile)
		}()
	}
	wg.Wait()
	close(results)
	succeeded := 0
	for err := range results {
		if err == nil {
			succeeded++
		} else if !errors.Is(err, os.ErrExist) {
			t.Errorf("unexpected error: %v", err)
		}
	}
	if succeeded != 1 {
		t.Fatalf("successful writers = %d, want 1", succeeded)
	}
	for _, name := range []string{"privateKey.pem", "publicKey.pem"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Error(err)
		}
	}
}

func TestRsaGenerateKeyPreservesExistingEntries(t *testing.T) {
	for _, name := range []string{"publicKey.pem", "privateKey.pem"} {
		for _, directory := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "-file", true: "-directory"}[directory], func(t *testing.T) {
				dir := t.TempDir()
				path := filepath.Join(dir, name)
				if directory {
					if err := os.Mkdir(path, 0700); err != nil {
						t.Fatal(err)
					}
				} else if err := os.WriteFile(path, []byte("existing key bytes"), 0644); err != nil {
					t.Fatal(err)
				}
				if got := rsaGeneratePanic(dir); got == nil {
					t.Error("expected a reported failure for an existing entry")
				}
				if !directory {
					data, err := os.ReadFile(path)
					if err != nil || string(data) != "existing key bytes" {
						t.Errorf("existing file changed: size=%d, err=%v", len(data), err)
					}
				}
				other := "publicKey.pem"
				if name == other {
					other = "privateKey.pem"
				}
				if _, err := os.Lstat(filepath.Join(dir, other)); !os.IsNotExist(err) {
					t.Errorf("unexpected counterpart after failure: %v", err)
				}
			})
		}
	}
}

func TestRsaGenerateKeyRoundTripAndPermissions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "keys")
	if got := RsaGenerateKey(dir); got != dir {
		t.Fatalf("path = %q", got)
	}
	privatePEM, err := os.ReadFile(filepath.Join(dir, "privateKey.pem"))
	if err != nil {
		t.Fatal(err)
	}
	publicPEM, err := os.ReadFile(filepath.Join(dir, "publicKey.pem"))
	if err != nil {
		t.Fatal(err)
	}
	privateKey := RsaLoadPrivateKey(string(privatePEM))
	publicKey := RsaLoadPublicKey(string(publicPEM))
	ciphertext := RsaEncrypt4RsaKey(publicKey, "round trip")
	if got := RsaDecrypt4RsaKey(privateKey, ciphertext); got != "round trip" {
		t.Fatalf("plaintext = %q", got)
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(filepath.Join(dir, "privateKey.pem"))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0600 {
			t.Errorf("private key mode = %o, want 600", info.Mode().Perm())
		}
	}
}
