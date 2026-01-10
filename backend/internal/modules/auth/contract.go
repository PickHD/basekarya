package auth

type Hasher interface {
	CheckPasswordHash(password, hash string) bool
}

type TokenProvider interface {
	GenerateToken(userID uint, role string) (string, error)
}
