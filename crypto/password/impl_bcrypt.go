package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher bcrypt 哈希器
type BcryptHasher struct {
	cost int
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	if h.cost == 0 {
		h.cost = bcrypt.DefaultCost
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (h *BcryptHasher) Compare(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrMismatchedHashAndPassword
	}
	return err
}
