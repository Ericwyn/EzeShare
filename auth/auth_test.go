package auth

import (
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"testing"
)

func TestRSAGenKey(t *testing.T) {
	RSAGenKey(true)
}

// 公钥加密
func TestEncryptRSA(t *testing.T) {
	rsa, err := EncryptRSA("hello", GetRsaPublicKeyPath())
	if err != nil {
		panic(err)
	}

	// base64 密文
	log.I("encrypt: \n", rsa)

	decryptRSA, err := DecryptRSA(rsa, GetRsaPrivateKeyPath())
	if err != nil {
		panic(err)
	}
	log.I("decrypt: \n", decryptRSA)
}

func TestTokenGen(t *testing.T) {
	storage.InitDb(true)
	TokenGen(false)
}
