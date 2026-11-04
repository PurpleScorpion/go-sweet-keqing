package keqing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func rsaTestKeyPair(t *testing.T) (string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privateDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})), string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateDER}))
}

func rsaTryLoad(public, private, keyType string) (err error, panicValue any) {
	defer func() { panicValue = recover() }()
	err = RsaLoadKey(public, private, keyType)
	return
}

func TestRsaLoadKeyFailurePreservesPair(t *testing.T) {
	publicA, privateA := rsaTestKeyPair(t)
	publicB, privateB := rsaTestKeyPair(t)
	dir := t.TempDir()
	publicFile := filepath.Join(dir, "public.pem")
	if err := os.WriteFile(publicFile, []byte(publicB), 0600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, public, private, keyType string
		wantPanic                      bool
	}{
		{"invalid public PEM", "invalid public PEM", privateB, RSA_KEY_STRING_TYPE, true},
		{"invalid private PEM", publicB, "invalid private PEM", RSA_KEY_STRING_TYPE, true},
		{"empty public", "", privateB, RSA_KEY_STRING_TYPE, true},
		{"empty private", publicB, "", RSA_KEY_STRING_TYPE, true},
		{"empty key type", publicB, privateB, "", true},
		{"invalid key type", publicB, privateB, "invalid type", false},
		{"mismatched pair", publicB, privateA, RSA_KEY_STRING_TYPE, false},
		{"missing public file", filepath.Join(dir, "missing"), publicFile, RSA_KEY_FILE_TYPE, false},
		{"missing private file", publicFile, filepath.Join(dir, "missing"), RSA_KEY_FILE_TYPE, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := RsaLoadKey(publicA, privateA, RSA_KEY_STRING_TYPE); err != nil {
				t.Fatal(err)
			}
			ciphertext := RsaEncrypt("before failure")
			err, panicValue := rsaTryLoad(tc.public, tc.private, tc.keyType)
			if tc.wantPanic {
				if panicValue == nil {
					t.Error("expected legacy validation panic")
				}
			} else {
				if panicValue != nil || err == nil {
					t.Errorf("load returned error=%v panic=%v; want error", err, panicValue)
				}
			}
			if got := RsaDecrypt(ciphertext); got != "before failure" {
				t.Errorf("old ciphertext no longer decrypts: %q", got)
			}
			if got := RsaDecryptCustom(privateA, RsaEncrypt("after failure")); got != "after failure" {
				t.Errorf("old public key changed: %q", got)
			}
		})
	}
}

func TestRsaLoadKeySuccessFileAndString(t *testing.T) {
	for _, keyType := range []string{RSA_KEY_STRING_TYPE, RSA_KEY_FILE_TYPE} {
		t.Run(keyType, func(t *testing.T) {
			public, private := rsaTestKeyPair(t)
			if keyType == RSA_KEY_FILE_TYPE {
				dir := t.TempDir()
				publicPath := filepath.Join(dir, "public.pem")
				privatePath := filepath.Join(dir, "private.pem")
				if err := os.WriteFile(publicPath, []byte(public), 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(privatePath, []byte(private), 0600); err != nil {
					t.Fatal(err)
				}
				public, private = publicPath, privatePath
			}
			if err := RsaLoadKey(public, private, keyType); err != nil {
				t.Fatal(err)
			}
			if got := RsaDecrypt(RsaEncrypt("loaded pair")); got != "loaded pair" {
				t.Fatalf("round trip = %q", got)
			}
		})
	}
}

func TestRsaLoadKeyConcurrentReload(t *testing.T) {
	public, private := rsaTestKeyPair(t)
	if err := RsaLoadKey(public, private, RSA_KEY_STRING_TYPE); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				if err := RsaLoadKey(public, private, RSA_KEY_STRING_TYPE); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				if got := RsaDecrypt(RsaEncrypt("concurrent")); got != "concurrent" {
					t.Errorf("concurrent round trip = %q", got)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestRsaLoadKeyPublishesMatchingSnapshot(t *testing.T) {
	publicA, privateA := rsaTestKeyPair(t)
	publicB, privateB := rsaTestKeyPair(t)
	if err := RsaLoadKey(publicA, privateA, RSA_KEY_STRING_TYPE); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for _, pair := range [][2]string{{publicA, privateA}, {publicB, privateB}} {
		pair := pair
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				if err := RsaLoadKey(pair[0], pair[1], RSA_KEY_STRING_TYPE); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 2000; n++ {
				snapshot := loadedRSAKeyPair()
				if snapshot.public == nil || snapshot.private == nil || !snapshot.public.Equal(&snapshot.private.PublicKey) {
					t.Error("observed incomplete or mismatched pair")
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestRsaUninitializedGlobalFunctionsStillPanic(t *testing.T) {
	rsaKeyMu.Lock()
	previous := rsaKeys
	rsaKeys = rsaKeyPair{}
	rsaKeyMu.Unlock()
	t.Cleanup(func() { rsaKeyMu.Lock(); rsaKeys = previous; rsaKeyMu.Unlock() })
	for _, call := range []func(){func() { RsaEncrypt("message") }, func() { RsaDecrypt("ciphertext") }} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("uninitialized global key must panic")
				}
			}()
			call()
		}()
	}
}
