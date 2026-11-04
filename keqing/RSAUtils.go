package keqing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type rsaKeyPair struct {
	public  *rsa.PublicKey
	private *rsa.PrivateKey
}

var (
	rsaKeyMu            sync.RWMutex
	rsaKeys             rsaKeyPair
	RSA_KEY_FILE_TYPE   = "FILE"
	RSA_KEY_STRING_TYPE = "STRING"
	RSA_PKCS1           = "PKCS1"
	RSA_PKCS8_SHA_256   = "sha256"
)

// Return one immutable pair snapshot; parsing and cryptographic operations run
// outside the lock. Keys published here are never changed in place.
func loadedRSAKeyPair() rsaKeyPair {
	rsaKeyMu.RLock()
	pair := rsaKeys
	rsaKeyMu.RUnlock()
	return pair
}

/*
加载密钥对
试用环境:

	整个项目只有一个RSA密钥对的
	若是有多个密钥对, 请直接使用 EncryptCustom 和 DecryptCustom 方法 ,
	若是想通过已加载好的秘钥进行加密解密 , 请使用 Encrypt4RsaKey 和 Decrypt4RsaKey 方法
*/
func RsaLoadKey(publicKey string, privateKey string, keyType string) error {
	if IsEmpty(keyType) {
		panic("KeyType cannot be empty")
	}
	if IsEmpty(publicKey) {
		panic("PublicKey cannot be empty")
	}
	if IsEmpty(privateKey) {
		panic("PrivateKey cannot be empty")
	}
	if keyType == RSA_KEY_FILE_TYPE {
		publicKeyPEM, err := os.ReadFile(publicKey)
		if err != nil {
			return err
		}
		privateKeyPEM, err := os.ReadFile(privateKey)
		if err != nil {
			return err
		}
		publicKey, privateKey = string(publicKeyPEM), string(privateKeyPEM)
	} else if keyType != RSA_KEY_STRING_TYPE {
		return errors.New("KeyType type error")
	}
	// Keep both parsed keys local until all validation succeeds. In particular,
	// a legacy PEM parsing panic must not change the currently loaded pair.
	private := RsaLoadPrivateKey(privateKey)
	public := RsaLoadPublicKey(publicKey)
	if !public.Equal(&private.PublicKey) {
		return errors.New("RSA public and private keys do not match")
	}
	rsaKeyMu.Lock()
	rsaKeys = rsaKeyPair{public: public, private: private}
	rsaKeyMu.Unlock()
	return nil
}

/*
 * 生成RSA密钥到指定目录下
 */
func RsaGenerateKey(path string) string {
	result, err := RsaGenerateKeyWithError(path)
	if err != nil {
		panic(err)
	}
	return result
}

// RsaGenerateKeyWithError creates a new key pair without overwriting either key file.
// On failure it returns an empty path and removes only files created by this call.
// Private keys use mode 0600 on Unix; Windows permissions follow the parent ACL.
func RsaGenerateKeyWithError(path string) (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("generate RSA key: %w", err)
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("encode RSA private key: %w", err)
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("encode RSA public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err := os.MkdirAll(path, 0700); err != nil {
		return "", fmt.Errorf("create RSA key directory: %w", err)
	}
	if err := createRSAKeyFiles(path, privateKeyPEM, publicKeyPEM, openRSAKeyFile); err != nil {
		return "", err
	}
	return path, nil
}

type rsaKeyFile interface {
	io.Writer
	Sync() error
	Close() error
	Stat() (os.FileInfo, error)
	Chmod(os.FileMode) error
}

func openRSAKeyFile(path string, flag int, mode os.FileMode) (rsaKeyFile, error) {
	return os.OpenFile(path, flag, mode)
}

// Opening both files exclusively before writing also prevents concurrent calls
// from replacing an existing pair. Cleanup never includes a failed open.
func createRSAKeyFiles(path string, privatePEM, publicPEM []byte, openFile func(string, int, os.FileMode) (rsaKeyFile, error)) (err error) {
	outputs := []struct {
		path string
		data []byte
		mode os.FileMode
		file rsaKeyFile
		info os.FileInfo
	}{
		{path: filepath.Join(path, "privateKey.pem"), data: privatePEM, mode: 0600},
		{path: filepath.Join(path, "publicKey.pem"), data: publicPEM, mode: 0644},
	}
	defer func() {
		// Close every handle before cleanup, including on Windows.
		for _, output := range outputs {
			if output.file != nil {
				if closeErr := output.file.Close(); closeErr != nil {
					err = errors.Join(err, fmt.Errorf("close RSA key %q: %w", output.path, closeErr))
				}
			}
		}
		if err == nil {
			return
		}
		for _, output := range outputs {
			if output.file == nil {
				continue
			}
			current, statErr := os.Lstat(output.path)
			if os.IsNotExist(statErr) {
				continue
			}
			if statErr != nil {
				err = errors.Join(err, fmt.Errorf("inspect RSA key for cleanup %q: %w", output.path, statErr))
				continue
			}
			// Do not remove a replacement installed by another actor.
			if output.info == nil || !os.SameFile(output.info, current) {
				err = errors.Join(err, fmt.Errorf("RSA key %q changed during cleanup; left untouched", output.path))
				continue
			}
			if removeErr := os.Remove(output.path); removeErr != nil {
				err = errors.Join(err, fmt.Errorf("remove incomplete RSA key %q: %w", output.path, removeErr))
			}
		}
	}()
	for i := range outputs {
		output := &outputs[i]
		file, openErr := openFile(output.path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, output.mode)
		if openErr != nil {
			return fmt.Errorf("create RSA key %q: %w", output.path, openErr)
		}
		output.file = file
		output.info, err = file.Stat()
		if err != nil {
			return fmt.Errorf("inspect new RSA key %q: %w", output.path, err)
		}
		if output.mode == 0600 {
			if err := file.Chmod(0600); err != nil {
				return fmt.Errorf("set RSA private key permissions: %w", err)
			}
		}
	}
	for _, output := range outputs {
		n, writeErr := output.file.Write(output.data)
		if writeErr == nil && n != len(output.data) {
			writeErr = io.ErrShortWrite
		}
		if writeErr != nil {
			return fmt.Errorf("write RSA key %q: %w", output.path, writeErr)
		}
		if syncErr := output.file.Sync(); syncErr != nil {
			return fmt.Errorf("sync RSA key %q: %w", output.path, syncErr)
		}
	}
	return nil
}

func RsaLoadPrivateKey(privateKeyPEM string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		panic("Invalid PEM data")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic("Failed to parse PKCS # 8 private key")
	}

	switch key := privKey.(type) {
	case *rsa.PrivateKey:
		return key
	default:
		panic("Unsupported key type")
	}
}

func RsaLoadPublicKey(publicKeyPEM string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("Invalid public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("Public key parsing failed")
	}

	switch pk := pubKey.(type) {
	case *rsa.PublicKey:
		return pk
	default:
		panic("unknown type of public key")
	}
}

// 使用PKCS1 进行加密, 这个通常与Java中的默认配置一致
func RsaEncrypt4PKCS1(publicKey *rsa.PublicKey, data string) string {
	if publicKey == nil {
		panic("请先加载公钥")
	}
	if IsEmpty(data) {
		panic("The data to be encrypted cannot be empty")
	}
	plaintext := []byte(data)

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext
}

/*
加密文本
*/
func RsaEncrypt(data string) string {
	publicKey := loadedRSAKeyPair().public
	if publicKey == nil {
		panic("请先加载公钥")
	}
	if IsEmpty(data) {
		panic("The data to be encrypted cannot be empty")
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, GetBytes(data), nil)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

/*
加密文本
*/
func RsaEncryptCustom(publicKey string, data string) string {
	if IsEmpty(publicKey) {
		panic("PublicKey cannot be empty")
	}
	if IsEmpty(data) {
		panic("The data to be encrypted cannot be empty")
	}
	pk := RsaLoadPublicKey(publicKey)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, GetBytes(data), nil)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

/*
加密文本
*/
func RsaEncrypt4RsaKey(publicKey *rsa.PublicKey, data string) string {
	if publicKey == nil {
		panic("PublicKey cannot be empty")
	}
	if IsEmpty(data) {
		panic("The data to be encrypted cannot be empty")
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, GetBytes(data), nil)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func RsaDecrypt4PKCS1(key *rsa.PrivateKey, ciphertextString string) string {
	if key == nil {
		panic("Please load the private key first")
	}
	if IsEmpty(ciphertextString) {
		panic("The ciphertext cannot be empty")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		log.Println("Password parsing failed:", err)
		return ""
	}
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, key, ciphertext)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return ""
	}
	return string(plaintext)
}

/*
解密文本
*/
func RsaDecrypt(ciphertextString string) string {
	privateKey := loadedRSAKeyPair().private
	if privateKey == nil {
		panic("Please load the private key first")
	}
	if IsEmpty(ciphertextString) {
		panic("The ciphertext cannot be empty")
	}
	data, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		log.Println("Password parsing failed:", err)
		return ""
	}
	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return ""
	}
	return string(decryptedText)
}

/*
解密文本
*/
func RsaDecryptCustom(privateKey string, ciphertextString string) string {
	if IsEmpty(privateKey) {
		panic("Private Key cannot be empty")
	}
	if IsEmpty(ciphertextString) {
		panic("The ciphertext cannot be empty")
	}
	pk := RsaLoadPrivateKey(privateKey)
	data, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		log.Println("Decryption failed:", err)
		return ""
	}
	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pk, data, nil)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return ""
	}
	return string(decryptedText)
}

/*
解密文本
*/
func RsaDecrypt4RsaKey(key *rsa.PrivateKey, ciphertextString string) string {
	if key == nil {
		panic("Private Key cannot be empty")
	}
	if IsEmpty(ciphertextString) {
		panic("The ciphertext cannot be empty")
	}
	data, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		log.Println("Decryption failed:", err)
		return ""
	}
	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, data, nil)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return ""
	}
	return string(decryptedText)
}
