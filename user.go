package gochat

type User struct {
	Id int
	Username string
}

type UserRepository interface {
	Store(username string, passwordHash []byte) (User, error)
	Find(id int) (User, error)
	GetUserPasswordHash(username string) (passwordHash []byte, err error)
}
