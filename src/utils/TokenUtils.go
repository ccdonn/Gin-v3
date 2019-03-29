package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"regexp"
	"strconv"
	"strings"

	"../constant"

	"github.com/google/uuid"
)

func CreateToken(agentID int32) (encryptedToken string, err error) {
	uuid, err := uuid.NewRandom()
	token := strings.Join([]string{strings.Replace(uuid.String(), "-", "", -1), strconv.Itoa(int(agentID))}, ".")
	// fmt.Println("token before encrypt", token)
	return EncryptToken(token)
}

func ExtractAgentID(token string) int32 {
	info, err := DecryptToken(token)
	if err != nil {
		log.Println(err)
		return -1
	}

	if ok, _ := regexp.MatchString("[0-9a-fA-f]{32}.\\d", info); !ok {
		return -1
	}

	uid := ""
	if index := strings.Index(info, "."); index > 0 {
		uid = info[index+1:]
	}

	i, err := strconv.Atoi(uid)
	if err != nil {
		log.Println(err)
		return -1
	}
	return int32(i)
}

func EncryptToken(token string) (string, error) {
	key := []byte(constant.TokenKey)
	plainText := []byte(token)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := []byte(constant.TokenIV)
	blockSize := block.BlockSize()
	origData := PKCS5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.URLEncoding.EncodeToString(crypted), nil
}

func DecryptToken(token string) (string, error) {
	key := []byte(constant.TokenKey)
	rtoken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	crypted := []byte(rtoken)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := []byte(constant.TokenIV)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
