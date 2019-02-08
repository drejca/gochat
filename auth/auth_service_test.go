package auth

import (
	"github.com/drejca/gochat/postgres"
	"testing"
	"time"
)

func TestRegisterAndLogin(t *testing.T) {
	username := "John"
	password := "1234567891123456789112345678"

	conn, err := postgres.New()
	if err != nil {
		t.Fatal(err)
	}

	postgres.Migrate(conn)

	userRepository := postgres.NewUserRepository(conn)
	authService := NewService(userRepository)

	user, err := authService.Register(username, password)
	if err != nil {
		t.Fatal(err)
	}

	if user.Username != username {
		t.Errorf("expected username %q got %q", username, user.Username)
	}

	err = authService.Login(username, password)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashPassword(t *testing.T) {
	password := "1234567891123456789112345678"

	start := time.Now()

	hashedPassword, err := hashPassword(password)
	if err != nil {
		t.Error(err)
	}

	elapsed := time.Since(start)
	if !(elapsed.Seconds() > 0.2 && elapsed.Seconds() < 0.3) {
		t.Errorf("expected to take between 0.2 and 0.3 sec. It took %f", elapsed.Seconds())
	}

	if !isPasswordHashValid(password, hashedPassword) {
		t.Error("password and hash do not match")
	}
}
