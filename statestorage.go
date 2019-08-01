package discordgo

type StorageEngine interface {
	GuildAdd(guild *Guild) error
	GuildRemove(guild *Guild) error
	Guild(guildID string) (*Guild, error)
	PresenceAdd(guildID string, presence *Presence) error
	PresenceRemove(guildID string, presence *Presence) error
	Presence(guildID string, userID string) (*Presence, error)
	MemberAdd(member *Member) error
	MemberRemove(member *Member) error
	Member(guildID, userID string) (*Member, error)
	RoleAdd(guildID string, role *Role)
	RoleRemove(guildID, roleID string) error
	Role(guildID, roleID string) (*Role, error)
	ChannelAdd(channel *Channel) error
	ChannelRemove(channel *Channel) error
	Channel(channelID string) (*Channel, error)
	Emoji(guildID, emojiID string) (*Emoji, error)
	EmojiAdd(guildID string, emoji *Emoji) error
	EmojisAdd(guildID string, emojis []*Emoji) error
	MessageAdd(message *Message) error
	MessageRemove(message *Message) error
	Message(channelID, messageID string) (*Message, error)
}
