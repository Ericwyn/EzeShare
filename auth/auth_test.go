package auth

import (
	"fmt"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"sync"
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

func TestGetRsaPublicKey(t *testing.T) {
	fmt.Println(GetRsaPublicKey())
}

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Add(-1)
	wg.Add(-1)
}
