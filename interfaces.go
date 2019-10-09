package discordgo

import "time"

type Messageable interface {
	// Sends a message to the channel
	SendMessage(content string, embed *MessageEmbed, files []*File) (message *Message, err error)
	SendMessageComplex(data *MessageSend) (message *Message, err error)
	EditMessage(data *MessageEdit) (edited *Message, err error)

	// gets a single message by ID from the channel
	// ID : the ID of a Message
	FetchMessage(ID string) (message *Message, err error)

	// returns an array of Message structures for messages within
	// a given channel.
	// channelID : The ID of a Channel.
	// limit     : The number messages that can be returned. (max 100)
	// beforeID  : If provided all messages returned will be before given ID.
	// afterID   : If provided all messages returned will be after given ID.
	// aroundID  : If provided all messages returned will be around given ID.
	GetHistory(limit int, beforeID, afterID, aroundID string) (st []*Message, err error)

	// GetHistoryIterator returns a bare HistoryIterator for this Messageable.
	GetHistoryIterator() *HistoryIterator
}

// IDGettable objects are objects where it is possible to get a snowflake ID from
type IDGettable interface {
	GetID() string
	CreatedAt() (creation time.Time, err error)
}

// Mentionable objects are objects that can have their snowflake ID be formatted as a mention in discord
type Mentionable interface {
	IDGettable
	Mention() string
}

// TimeSortable objects are objects which can return their creation time and
// thus can be sorted on when they were created
type TimeSortable interface {
	CreatedAt() (creation time.Time, err error)
}
