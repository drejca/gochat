package gochat

import "time"

type Message struct {
	Id int
	OwnerId int
	ChannelId int
	Content string
	OnDate time.Time
}

type MessageRepository interface {
	Store(ownerId int, channelId int, message string) error
	ChannelMessages(channelId int, before time.Time) ([]Message, error)
}