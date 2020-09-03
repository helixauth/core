package utils

import (
	"crypto/sha256"

	"github.com/google/uuid"
)

func Hash(key string) string {
	return uuid.NewHash(sha256.New(), uuid.NameSpaceOID, []byte(key), 5).String()
}
