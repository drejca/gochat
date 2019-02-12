package channel

import (
	"github.com/drejca/gochat"
	"github.com/drejca/gochat/auth"
	"github.com/drejca/gochat/postgres"
	"testing"
	"time"
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

func TestJoinToChannel(t *testing.T) {
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

	err = channelService.JoinToChannel(user.Id, channel.Id)
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

	content := `Hi`

	err = channelService.SendMessage(owner.Id, channel.Id, content)
	if err != nil {
		t.Error(err)
	}
}

func TestChannelMessages(t *testing.T) {
	conn := postgresConnection(t)

	postgres.Migrate(conn)

	channelService := newChannelService(conn)
	authService := auth.NewService(postgres.NewUserRepository(conn))

	owner := setupUser(t, authService, "John", "randompassword123")

	channel, err := channelService.CreateChannel(owner.Id, "general discord")
	if err != nil {
		t.Fatal(err)
	}

	content := `Hi`

	err = channelService.SendMessage(owner.Id, channel.Id, content)
	if err != nil {
		t.Fatal(err)
	}

	messages, err := channelService.ChannelMessages(channel.Id, time.Now().UTC().Add(time.Microsecond))
	if err != nil {
		t.Fatal(err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected one message got %d", len(messages))
	}

	if messages[0].Content != content {
		t.Errorf("expected content to be %q, got %q", content, messages[0].Content)
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
	messageRepository := postgres.NewMessageRepository(conn)

	return NewService(channelRepository, userRepository, messageRepository)
}

func setupUser(t *testing.T, authService *auth.Service, username string, password string) gochat.User {
	user, err := authService.Register(username, password)
	if err != nil {
		t.Fatal(err)
	}
	return user
}
