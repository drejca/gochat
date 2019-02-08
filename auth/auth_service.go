package auth

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/drejca/gochat"
)

type Service struct {
	userRepository gochat.UserRepository
}

func NewService(userRepository gochat.UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) Register(username string, password string) (gochat.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return gochat.User{}, err
	}
	return s.userRepository.Store(username, passwordHash)
}

func (s *Service) Login(username string, password string) error {
	passwordHash, err := s.userRepository.GetUserPasswordHash(username)
	if err != nil {
		return err
	}

	if !isPasswordHashValid(password, passwordHash) {
		return errors.New("Wrong username or password")
	}
	return nil
}

func hashPassword(password string) (passwordHash []byte, err error) {
	passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return []byte{}, err
	}
	return passwordHash, nil
}

func isPasswordHashValid(password string, passwordHash []byte) bool {
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	return err == nil
}
