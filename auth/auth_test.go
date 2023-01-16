package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan"
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
	tokenSelf := GetSelfToken()
	log.I("token self: ", tokenSelf)
}

func TestGetRsaPublicKey(t *testing.T) {
	fmt.Println(GetRsaPublicKey())
}

func TestWaitGroup(t *testing.T) {
	a := func() {
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Add(-1)
		wg.Add(-1)
	}

	defer func() {
		fmt.Println("c")
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
		fmt.Println("d")
	}()

	a()
}

func TestJsonMarshal(t *testing.T) {
	msg := scan.BroadcastMsg{
		Name:       "nameee",
		Address:    "addressss",
		DeviceId:   "deviceIdddd",
		DeviceType: "deviceTypeeeee",
	}

	marshal, _ := json.Marshal(msg)
	fmt.Println(string(marshal))
}
