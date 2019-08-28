package discordgo

// MessageReaction stores the data for a message reaction.
type MessageReaction struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     Emoji  `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`

	Session *Session `json:"-"`
}

// Remove removes the reaction from the message it was added to
func (r *MessageReaction) Remove() error {
	return r.Session.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
}

// GetGuild retrieves the guild that the message that was reacted to was posted in,
// will return an error if it was in dms
func (r *MessageReaction) GetGuild() (g *Guild, err error) {
	return r.Session.State.Guild(r.GuildID)
}

// GetChannel retrieves the channel that the message that was reacted to was posted in
func (r *MessageReaction) GetChannel() (c *Channel, err error) {
	return r.Session.State.Channel(r.ChannelID)
}

// GetMessage retrieves the message that was reacted to
func (r *MessageReaction) GetMessage() (m *Message, err error) {
	m, err = r.Session.State.Message(r.ChannelID, r.MessageID)
	if err == nil {
		return
	}

	return r.Session.ChannelMessage(r.ChannelID, r.MessageID)
}

// GetMember retrieves the member that added the reaction,
// will return an error if it was not in a guild
func (r *MessageReaction) GetMember() (m *Member, err error) {
	g, err := r.GetGuild()
	if err != nil {
		return
	}

	return g.GetMember(r.UserID)
}

// GetUser retrieves the user that added the reaction
func (r *MessageReaction) GetUser() (u *User, err error) {
	return r.Session.State.GetUser(r.UserID)
}
