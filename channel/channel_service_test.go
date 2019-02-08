package channel

import (
	"github.com/drejca/gochat"
	"github.com/drejca/gochat/auth"
	"github.com/drejca/gochat/postgres"
	"testing"
)

func TestCreateChannel(t *testing.T) {
	conn := postgresConnection(t)

	postgres.Migrate(conn)

	channelService := newChannelService(conn)
	authService := auth.NewService(postgres.NewUserRepository(conn))

	owner := setupUser(t, authService, "John", "longpassword123456")

	channelName := "general discord"

	channel, err := channelService.CreateChannel(owner.Id, channelName)
	if err != nil {
		t.Fatal(err)
	}

	if channel.Name != channelName {
		t.Errorf("expected channel name %q got %q", channelName, channel.Name)
	}

	if channel.OwnerId != owner.Id {
		t.Errorf("expected channel ownerId to be %d got %d", owner.Id, channel.OwnerId)
	}

	users, err := channelService.GetChannelUsers(channel.Id)

	if len(users) == 0 || users[0].Username != owner.Username {
		t.Errorf("missing user %q in channel %q", owner.Username, channel.Name)
	}
}

func TestAddUserToChannel(t *testing.T) {
	conn := postgresConnection(t)

	postgres.Migrate(conn)

	channelService := newChannelService(conn)
	authService := auth.NewService(postgres.NewUserRepository(conn))

	owner := setupUser(t, authService, "John", "longpassword123456")

	username := "Mary"
	user := setupUser(t, authService, username, "longpassword123456")

	channel, err := channelService.CreateChannel(owner.Id, "general discord")
	if err != nil {
		t.Fatal(err)
	}

	err = channelService.AddUserToChannel(user.Id, channel.Id)
	if err != nil {
		t.Fatal(err)
	}

	users, err := channelService.GetChannelUsers(channel.Id)

	if len(users) != 2 || users[1].Username != username {
		t.Errorf("missing user %q in channel %q", username, channel.Name)
	}
}

func TestSendMessageToChannel(t *testing.T) {
	conn := postgresConnection(t)

	postgres.Migrate(conn)

	channelService := newChannelService(conn)
	authService := auth.NewService(postgres.NewUserRepository(conn))

	owner := setupUser(t, authService, "John", "longpassword123456")

	channel, err := channelService.CreateChannel(owner.Id, "general discord")
	if err != nil {
		t.Fatal(err)
	}

	content := `<p>Hi</p>`

	err = channelService.SendMessage(owner.Id, channel.Id, content)
	if err != nil {
		t.Error(err)
	}
}

func postgresConnection(t *testing.T) *postgres.Conn {
	conn, err := postgres.New()
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func newChannelService(conn *postgres.Conn) *Service {
	channelRepository := postgres.NewChannelRepository(conn)
	userRepository := postgres.NewUserRepository(conn)
	channelUsersRepository := postgres.NewChannelUsersRepository(conn)
	messageRepository := postgres.NewMessageRepository(conn)

	return NewService(channelRepository, userRepository, channelUsersRepository, messageRepository)
}

func setupUser(t *testing.T, authService *auth.Service, username string, password string) gochat.User {
	user, err := authService.Register(username, password)
	if err != nil {
		t.Fatal(err)
	}
	return user
}