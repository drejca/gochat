package gochat

type Message struct {
	Id int
	OwnerId int
	ChannelId int
	Message string
}

type MessageRepository interface {
	Store(ownerId int, channelId int, message string) error
}