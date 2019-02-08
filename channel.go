package gochat

type Channel struct {
	Id int
	Name string
	OwnerId int
}

type ChannelRepository interface {
	Store(channelName string, owner User) (Channel, error)
	Find(channelId int) (Channel, error)
}

type ChannelUsers struct {
	Id int
	ChannelId int
	UserId int
}

type ChannelUsersRepository interface {
	Store(channel Channel, user User) error
	GetChannelUsers(channelId int) ([]User, error)
}
