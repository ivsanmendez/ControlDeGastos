package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher implements user.PasswordHasher using bcrypt.
type Hasher struct {
	cost int
}

func New() *Hasher {
	return &Hasher{cost: bcrypt.DefaultCost}
}

func (h *Hasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *Hasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
