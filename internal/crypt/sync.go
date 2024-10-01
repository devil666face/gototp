package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"math/rand"
	"unsafe"

	crand "crypto/rand"
)

const aesLen = 32
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func randString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func AES32RandomKey() []byte {
	return []byte(randString(aesLen))
}

type Sync struct {
	AesKey []byte
}

func New(key []byte) (*Sync, error) {
	_sync := &Sync{AesKey: AES32RandomKey()}
	b64key := base64.StdEncoding.EncodeToString(key)
	if err := _sync.WithB64Key(b64key); err != nil {
		return nil, err
	}
	return _sync, nil
}

func (s *Sync) WithB64Key(b64key string) error {
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}
	if len(key) < aesLen {
		return fmt.Errorf("key is too short, it must more then %d symbols", aesLen)
	}
	s.AesKey = key[:aesLen]
	return nil
}

func (s *Sync) B64Key() string {
	return base64.StdEncoding.EncodeToString(s.AesKey)
}

func (s *Sync) Encrypt(data []byte) ([]byte, error) {
	aes, err := aes.NewCipher(s.AesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = crand.Read(nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (s *Sync) Decrypt(data []byte) ([]byte, error) {
	aes, err := aes.NewCipher(s.AesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
