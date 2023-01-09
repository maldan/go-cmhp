package cmhp_crypto

import (
	"crypto/rand"
	"math/big"
	"os"
	"runtime"
	"time"
)

var m runtime.MemStats
var a = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
var a2 = "0123456789"
var al = len(a)
var al2 = len(a2)
var b = big.NewInt(1_000_000_000)

func init() {
	runtime.ReadMemStats(&m)
}

func UID(size int) string {
	out := make([]byte, size)
	t := time.Now().Unix() + int64(os.Getpid()) + int64(m.TotalAlloc+m.Alloc+m.Sys)

	for i := 0; i < size; i++ {
		num, _ := rand.Int(rand.Reader, b)
		out[i] = a[int(num.Int64()+t)%al]
	}

	return string(out)
}

func RandomCode(size int) string {
	out := make([]byte, size)
	t := time.Now().Unix() + int64(os.Getpid()) + int64(m.TotalAlloc+m.Alloc+m.Sys)

	for i := 0; i < size; i++ {
		num, _ := rand.Int(rand.Reader, b)
		out[i] = a2[int(num.Int64()+t)%al2]
	}

	return string(out)
}
