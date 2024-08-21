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
	"log"
	"os"
)

var (
	publicKeyBlock      *rsa.PublicKey  = nil
	privateKeyBlock     *rsa.PrivateKey = nil
	RSA_KEY_FILE_TYPE                   = "FILE"
	RSA_KEY_STRING_TYPE                 = "STRING"
)

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
		privateKeyBlock = RsaLoadPrivateKey(string(privateKeyPEM))
		publicKeyBlock = RsaLoadPublicKey(string(publicKeyPEM))
		return nil
	} else if keyType == RSA_KEY_STRING_TYPE {
		privateKeyBlock = RsaLoadPrivateKey(privateKey)
		publicKeyBlock = RsaLoadPublicKey(publicKey)
		return nil
	}
	return errors.New("KeyType type error")
}

/*
 * 生成RSA密钥到指定目录下
 */
func RsaGenerateKey(path string) string {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	// 将私钥编码为PKCS#8格式
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	// 提取公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	MkdirAll(path)
	err = writeToFile(path+"/publicKey.pem", publicKeyPEM)
	err = writeToFile(path+"/privateKey.pem", privateKeyPEM)
	return path
}

func writeToFile(filePath string, pemBlock *pem.Block) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(file, pemBlock)
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

/*
加密文本
*/
func RsaEncrypt(data string) string {
	if publicKeyBlock == nil {
		panic("请先加载公钥")
	}
	if IsEmpty(data) {
		panic("The data to be encrypted cannot be empty")
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKeyBlock, GetBytes(data), nil)
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

/*
解密文本
*/
func RsaDecrypt(ciphertextString string) string {
	if privateKeyBlock == nil {
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
	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKeyBlock, data, nil)
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
