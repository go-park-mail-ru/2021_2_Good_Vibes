package impl

import "golang.org/x/crypto/bcrypt"

type HasherBCrypt struct {
	cost int
}

func NewHasherBCrypt(cost int) *HasherBCrypt {
	return &HasherBCrypt{cost: cost}
}

func (h *HasherBCrypt) CompareHashAndPassword(hashPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashPassword, password)
}

func (h *HasherBCrypt) GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, h.cost)
}
