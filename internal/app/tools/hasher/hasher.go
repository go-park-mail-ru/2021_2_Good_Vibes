package hasher

//go:generate mockgen -source=hasher.go -destination=mock/hasher_mock.go
type Hasher interface {
	CompareHashAndPassword(hashPassword []byte, password []byte) error
	GenerateFromPassword(password []byte) ([]byte, error)



}
