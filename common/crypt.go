package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/xxtea/xxtea-go/xxtea"
)

// UUIDPath default path for saving machine ID
var UUIDPath = homeDir() + "/.optool-id"
var appended = []byte{
	110, 67, 104, 105, 110, 97, 46, 230, 136, 145, 230, 157, 165, 232, 135, 170, 228, 184, 173, 229, 155, 189, 227, 128, 130,
}

// Encrypt encrypt string with uuid
func Encrypt(s string) string {
	enc := xxtea.Encrypt([]byte(s), getUUID())
	return base64.URLEncoding.EncodeToString(enc)
}

// Decrypt decrypt string with uuid
func Decrypt(s string) []byte {
	dec, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		fmt.Println("Decode string ", s, "err:", err)
		os.Exit(1)
	}
	return xxtea.Decrypt(dec, getUUID())
}

func getUUID() []byte {
	_, err := os.Stat(UUIDPath)
	if err != nil {
		genUUID()
	}
	UUID, err := ioutil.ReadFile(UUIDPath)
	if err != nil || len(UUID) < 10 {
		fmt.Println("ERROR read UUID", err)
	}
	return append(UUID, appended...)
}

func genUUID() {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		fmt.Println("Get Random UUID failed.", err)
		os.Exit(1)
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
	if err := ioutil.WriteFile(UUIDPath, []byte(hex.EncodeToString(h.Sum(nil))), 0700); err != nil {
		fmt.Println("Write UUID failed. ", err)
		os.Exit(2)
	}
}
