package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/GoTools/file"
	"os"
)

const rsaKeyBits = 2048

// TokenGen 生成认证 token
func TokenGen() {

}

func isRsaKeyExits() bool {
	privateKeyFile := file.OpenFile(GetRsaPrivateKeyPath())
	publicKeyFile := file.OpenFile(GetRsaPublicKeyPath())
	if privateKeyFile.Exits() || publicKeyFile.Exits() {
		log.I("rsa key alardey")
		return true
	}
	return false
}

func RSAGenKey(reGen bool) error {
	// 判断是否已经存在文件，如果存在的话就不生成
	if !reGen && isRsaKeyExits() {
		log.I("RSA KEY is exits, path: " + storage.GetConfigDirPath())
		return nil
	}

	// 生成私钥
	//1、使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	if err != nil {
		return err
	}

	//2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串
	privateStream, err := x509.MarshalPKCS8PrivateKey(privateKey)
	//3、将私钥字符串设置到pem格式块中
	block1 := pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	//4、通过pem将设置的数据进行编码，并写入磁盘文件
	fPrivate, err := os.Create(GetRsaPrivateKeyPath())
	if err != nil {
		return err
	}

	err = pem.Encode(fPrivate, &block1)
	if err != nil {
		return err
	}

	// 生成公钥
	publicKey := privateKey.PublicKey
	publicStream, err := x509.MarshalPKIXPublicKey(&publicKey)
	//publicStream:=x509.MarshalPKCS1PublicKey(&publicKey)
	block2 := pem.Block{
		Type:  "public key",
		Bytes: publicStream,
	}
	fPublic, err := os.Create(GetRsaPublicKeyPath())
	if err != nil {
		return err
	}
	//defer fPublic.Close()
	pem.Encode(fPublic, &block2)
	return nil
}

func GetRsaPrivateKeyPath() string {
	return storage.GetConfigDirPath() + "/privateKey.pem"
}

func GetRsaPublicKeyPath() string {
	return storage.GetConfigDirPath() + "/publicKey.pem"
}

// EncryptRSA 对数据进行加密操作
func EncryptRSA(str string, path string) (base64EncryptStr string, err error) {

	//1.获取秘钥（从本地磁盘读取）
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	b := make([]byte, fileInfo.Size())
	f.Read(b)
	// 2、将得到的字符串解码
	block, _ := pem.Decode(b)

	// 使用X509将解码之后的数据 解析出来
	//x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes) //对应于生成秘钥的x509.MarshalPKIXPublicKey(&publicKey)
	//keyInit1,err:=x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	//4.使用公钥加密数据
	pubKey := keyInit.(*rsa.PublicKey)
	res, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(str))

	base64Str := base64.StdEncoding.EncodeToString(res)
	return base64Str, err
}

// DecryptRSA 对数据进行解密操作
func DecryptRSA(base64EncryptStr, path string) (decryptStr string, err error) {
	// 先转成 base64 的 byte
	base64Bytes, err := base64.StdEncoding.DecodeString(base64EncryptStr)
	if err != nil {
		log.E("Decrypt RSA Error, base 64 decode fail")
		log.E(err)
		//fmt.Println(base64Bytes)
		return
	}

	//1.获取秘钥（从本地磁盘读取）
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	b := make([]byte, fileInfo.Size())
	f.Read(b)
	block, _ := pem.Decode(b)                                 //解码
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes) //还原数据
	res, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), base64Bytes)
	return string(res), err
}