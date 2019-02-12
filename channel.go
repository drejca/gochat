package gochat

type Channel struct {
	Id int
	Name string
	OwnerId int
}

type ChannelRepository interface {
	Store(channelName string, owner User) (Channel, error)
	Find(channelId int) (Channel, error)
	JoinToChannel(userId int, channelId int) error
	ChannelUsers(channelId int) ([]User, error)
}
