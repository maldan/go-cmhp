package cmhp

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashSha256(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
