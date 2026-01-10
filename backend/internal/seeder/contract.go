package seeder

type Hasher interface {
	HashPassword(password string) (string, error)
}
