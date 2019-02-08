package channel

import "github.com/drejca/gochat"

type Service struct {
	channelRepository gochat.ChannelRepository
	userRepository gochat.UserRepository
	channelUsersRepository gochat.ChannelUsersRepository
	messageRepository gochat.MessageRepository
}

func NewService(channelRepository gochat.ChannelRepository, userRepository gochat.UserRepository,
	channelUsersRepository gochat.ChannelUsersRepository, messageRepository gochat.MessageRepository) *Service {
	return &Service{
		channelRepository: channelRepository,
		userRepository: userRepository,
		channelUsersRepository: channelUsersRepository,
		messageRepository: messageRepository,
	}
}

func (s *Service) CreateChannel(ownerId int, channelName string) (gochat.Channel, error) {
	user, err := s.userRepository.Find(ownerId)
	if err != nil {
		return gochat.Channel{}, err
	}

	channel, err := s.channelRepository.Store(channelName, user)
	if err != nil {
		return channel, err
	}

	err = s.AddUserToChannel(user.Id, channel.Id)

	return channel, err
}

func (s *Service) AddUserToChannel(userId int, channelId int) error {
	channel, err := s.channelRepository.Find(channelId)
	if err != nil {
		return err
	}

	user, err := s.userRepository.Find(userId)
	if err != nil {
		return err
	}

	return s.channelUsersRepository.Store(channel, user)
}

func (s *Service) GetChannelUsers(channelId int) ([]gochat.User, error) {
	return s.channelUsersRepository.GetChannelUsers(channelId)
}

func (s *Service) SendMessage(userId int, channelId int, content string) error {
	user, err := s.userRepository.Find(userId)
	if err != nil {
		return err
	}

	channel, err := s.channelRepository.Find(channelId)
	if err != nil {
		return err
	}

	return s.messageRepository.Store(user.Id, channel.Id, content)
}
