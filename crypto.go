package cmhp

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"os"
	"runtime"
	"time"
)

func HashSha1(data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func HashSha256(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func UID(size int) string {
	a := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

	var m runtime.MemStats //
	runtime.ReadMemStats(&m)
	out := ""
	t := time.Now().Unix() + int64(os.Getpid()) + int64(m.TotalAlloc+m.Alloc+m.Sys)

	for i := 0; i < size; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(1_000_000_000))
		x := int(num.Int64() + t)
		out += string(a[x%len(a)])
	}

	return out
}
