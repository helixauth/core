package utils

import (
	"crypto/sha256"
	"log"

	"github.com/helixauth/core/src/lib/mapper"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Hash(key string) string {
	return uuid.NewHash(sha256.New(), uuid.NameSpaceOID, []byte(key), 5).String()
}

func HashPassword(password *string) (*string, error) {
	if password == nil {
		return nil, nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return mapper.StringPtr(string(hash)), nil
}
