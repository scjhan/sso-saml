package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func pkcs5Padding(input []byte, blockSize int) []byte {
	padding := blockSize - len(input)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(input, padtext...)
}

func pkcs5UnPadding(input []byte) []byte {
	length := len(input)
	unpadding := int(input[length-1])
	return input[:(length - unpadding)]
}

// DesEncrypt ...
func DesEncrypt(input, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bm := cipher.NewCBCEncrypter(block, key)
	input = pkcs5Padding(input, bm.BlockSize())

	dest := make([]byte, len(input))
	bm.CryptBlocks(dest, input)

	return dest, nil
}

// DesDecrypt ...
func DesDecrypt(input, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bm := cipher.NewCBCDecrypter(block, key)
	dest := make([]byte, len(input))
	bm.CryptBlocks(dest, input)

	dest = pkcs5UnPadding(dest)

	return dest, nil
}

// RsaEncrypt ...
func RsaEncrypt(input, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("invalid public key")
	}
	pubi, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pubi.(*rsa.PublicKey), input)
}

// RsaDecrypt ...
func RsaDecrypt(input, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("invalid private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, input)
}
