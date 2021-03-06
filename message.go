// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains code related to the Message struct

package discordgo

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// MessageType is the type of Message
type MessageType int

// Block contains the valid known MessageType values
const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
	MessageTypeUserPremiumGuildSubscription
	MessageTypeUserPremiumGuildSubscriptionTierOne
	MessageTypeUserPremiumGuildSubscriptionTierTwo
	MessageTypeUserPremiumGuildSubscriptionTierThree
	MessageTypeChannelFollowAdd
)

// A Message stores all data related to a specific Discord message.
type Message struct {
	// The ID of the message.
	ID string `json:"id"`

	// The ID of the channel in which the message was sent.
	ChannelID string `json:"channel_id"`

	// The ID of the guild in which the message was sent.
	GuildID string `json:"guild_id,omitempty"`

	// The content of the message.
	Content string `json:"content"`

	// The time at which the messsage was sent.
	// CAUTION: this field may be removed in a
	// future API version; it is safer to calculate
	// the creation time via the ID.
	Timestamp Timestamp `json:"timestamp"`

	// The time at which the last edit of the message
	// occurred, if it has been edited.
	EditedTimestamp Timestamp `json:"edited_timestamp"`

	// The roles mentioned in the message.
	MentionRoles []string `json:"mention_roles"`

	// Whether the message is text-to-speech.
	Tts bool `json:"tts"`

	// Whether the message mentions everyone.
	MentionEveryone bool `json:"mention_everyone"`

	// The author of the message. This is not guaranteed to be a
	// valid user (webhook-sent messages do not possess a full author).
	Author *User `json:"author"`

	// A list of attachments present in the message.
	Attachments []*MessageAttachment `json:"attachments"`

	// A list of embeds present in the message. Multiple
	// embeds can currently only be sent by webhooks.
	Embeds []*MessageEmbed `json:"embeds"`

	// A list of users mentioned in the message.
	Mentions []*User `json:"mentions"`

	// A list of reactions to the message.
	Reactions []*MessageReactions `json:"reactions"`

	// If the message has been pinned to the channel
	Pinned bool `json:"pinned"`

	// The type of the message.
	Type MessageType `json:"type"`

	// The webhook ID of the message, if it was generated by a webhook
	WebhookID string `json:"webhook_id"`

	// Member properties for this message's author,
	// contains only partial information
	Member *Member `json:"member"`

	// Channels specifically mentioned in this message
	// Not all channel mentions in a message will appear in mention_channels.
	// Only textual channels that are visible to everyone in a lurkable guild will ever be included.
	// Only crossposted messages (via Channel Following) currently include mention_channels at all.
	// If no mentions in the message meet these requirements, this field will not be sent.
	MentionChannels []*Channel `json:"mention_channels"`

	// Is sent with Rich Presence-related chat embeds
	Activity *MessageActivity `json:"activity"`

	// Is sent with Rich Presence-related chat embeds
	Application *MessageApplication `json:"application"`

	// MessageReference contains reference data sent with crossposted messages
	MessageReference *MessageReference `json:"message_reference"`

	// The flags of the message, which describe extra features of a message.
	// This is a combination of bit masks; the presence of a certain permission can
	// be checked by performing a bitwise AND between this int and the flag.
	Flags MessageFlag `json:"flags"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"-"`
}

// File stores info about files you e.g. send in messages.
type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
}

// MessageSend stores all parameters you can send with ChannelMessageSendComplex.
type MessageSend struct {
	Content string        `json:"content,omitempty"`
	Embed   *MessageEmbed `json:"embed,omitempty"`
	Tts     bool          `json:"tts"`
	Files   []*File       `json:"-"`

	// TODO: Remove this when compatibility is not required.
	File *File `json:"-"`

	// to make us sanitize the text so no everyone/here pings
	noEveryoneSanitization bool `json:"-"`
}

// MessageEdit is used to chain parameters via ChannelMessageEditComplex, which
// is also where you should get the instance from.
type MessageEdit struct {
	// The content of the message.
	Content *string `json:"content,omitempty"`

	// The embed attached to the message
	Embed *MessageEmbed `json:"embed,omitempty"`

	// The flags of the message, which describe extra features of a message.
	// This is a combination of bit masks; the presence of a certain permission can
	// be checked by performing a bitwise AND between this int and the flag. C
	Flags MessageFlag `json:"flags"`

	ID      string
	Channel string
	message *Message
}

// GetID returns the ID of the message
func (m Message) GetID() string {
	return m.ID
}

// CreatedAt returns the messages creation time in UTC
func (m Message) CreatedAt() (creation time.Time, err error) {
	return SnowflakeToTime(m.ID)
}

// NewMessageEdit returns a MessageEdit struct,
// given a Channel ID and message ID.
func NewMessageEdit(channelID string, messageID string) *MessageEdit {
	return &MessageEdit{
		Channel: channelID,
		ID:      messageID,
	}
}

// NewMessageEdit returns a MessageEdit struct, initialized
// with the Channel and message ID.
func (m *Message) NewMessageEdit() *MessageEdit {
	return &MessageEdit{
		Channel: m.ChannelID,
		ID:      m.ID,
		message: m,
	}
}

// SetContent is the same as setting the variable Content,
// except it doesn't take a pointer.
func (m *MessageEdit) SetContent(str string) *MessageEdit {
	m.Content = &str
	return m
}

// SetEmbed is a convenience function for setting the embed,
// so you can chain commands.
func (m *MessageEdit) SetEmbed(embed *MessageEmbed) *MessageEdit {
	m.Embed = embed
	return m
}

// ToggleEmbedSuppression toggles if the embeds in the message have been suppressed or not
func (m *MessageEdit) ToggleEmbedSuppression() *MessageEdit {
	m.Flags ^= MessageFlagSuppressEmbeds
	return m
}

// Edit takes the MessageEdit and edits the message,
// this only works when the MessageEdit was created with Message.NewMessageEdit()
func (m *MessageEdit) Edit() (res *Message, err error) {
	if m.message == nil {
		err = ErrObjectNotFound
		return
	}

	return m.message.Edit(m)
}

// Channel returns the channel object that the message was posted in,
// this should only ever be nil right after a RESUME or the initial connection
func (m *Message) Channel() *Channel {
	c, err := m.Session.State.Channel(m.ChannelID)
	if err != nil {
		m.Session.log(LogError, "Unable to get channel from state: %s", err)
	}
	return c
}

// Guild returns the guild a message was posted in if applicable, else guild is nil
func (m *Message) Guild() *Guild {
	g, _ := m.Session.State.Guild(m.GuildID)
	return g
}

// JumpURL returns the url to jump to the message
func (m *Message) JumpURL() string {
	g := m.GuildID
	if g == "" {
		g = "@me"
	}
	return fmt.Sprintf("https://discordapp.com/channels/%s/%s/%s", g, m.ChannelID, m.ID)
}

// EmbedSuppressed returns true if the embed(s) of the message have been suppressed
func (m *Message) EmbedSuppressed() bool {
	return m.Flags&1<<2 == 1<<2
}

// IsCrosspost returns true if the message is a crosspost through channel following
func (m *Message) IsCrosspost() bool {
	return m.Flags&1<<1 == 1<<1
}

// HasBeenCrossposted returns true if the message has been crossposted (published) through channel following
func (m *Message) HasBeenCrossposted() bool {
	return m.Flags&1<<0 == 1<<0
}

// A MessageAttachment stores data for message attachments.
type MessageAttachment struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
}

// MessageReactions holds a reactions object for a message.
type MessageReactions struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

// MessageActivity is sent with Rich Presence-related chat embeds
type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id"`
}

// MessageActivityType is the type of message activity
type MessageActivityType int

// Constants for the different types of Message Activity
const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)

// MessageFlag describes an extra feature of the message
type MessageFlag int

// Constants for the different bit offsets of Message Flags
const (
	// This message has been published to subscribed channels (via Channel Following)
	MessageFlagCrossposted MessageFlag = 1 << iota
	// This message originated from a message in another channel (via Channel Following)
	MessageFlagIsCrosspost
	// Do not include any embeds when serializing this message
	MessageFlagSuppressEmbeds
)

// MessageApplication is sent with Rich Presence-related chat embeds
type MessageApplication struct {
	ID          string `json:"id"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
}

// MessageReference contains reference data sent with crossposted messages
type MessageReference struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

// ContentWithMentionsReplaced will replace all @<id> mentions with the
// username of the mention.
func (m *Message) ContentWithMentionsReplaced() (content string) {
	content = m.Content

	for _, user := range m.Mentions {
		content = strings.NewReplacer(
			"<@"+user.ID+">", "@"+user.Username,
			"<@!"+user.ID+">", "@"+user.Username,
		).Replace(content)
	}
	return
}

var patternChannels = regexp.MustCompile("<#[^>]*>")

// ContentWithMoreMentionsReplaced will replace all @<id> mentions with the
// username of the mention, but also role IDs and more.
func (m *Message) ContentWithMoreMentionsReplaced() (content string, err error) {
	content = m.Content

	if !m.Session.StateEnabled {
		content = m.ContentWithMentionsReplaced()
		return
	}

	channel := m.Channel()
	if channel == nil {
		content = m.ContentWithMentionsReplaced()
		return
	}

	for _, user := range m.Mentions {
		nick := user.Username

		member, err := m.Session.State.Member(channel.GuildID, user.ID)
		if err == nil && member.Nick != "" {
			nick = member.Nick
		}

		content = strings.NewReplacer(
			"<@"+user.ID+">", "@"+user.Username,
			"<@!"+user.ID+">", "@"+nick,
		).Replace(content)
	}
	for _, roleID := range m.MentionRoles {
		role, err := m.Session.State.Role(channel.GuildID, roleID)
		if err != nil || !role.Mentionable {
			continue
		}

		content = strings.Replace(content, "<@&"+role.ID+">", "@"+role.Name, -1)
	}

	content = patternChannels.ReplaceAllStringFunc(content, func(mention string) string {
		channel, err := m.Session.State.Channel(mention[2 : len(mention)-1])
		if err != nil || channel.Type == ChannelTypeGuildVoice {
			return mention
		}

		return "#" + channel.Name
	})
	return
}

// Edit edits the message, replacing it entirely with
// the given MessageEdit struct
func (m *Message) Edit(data *MessageEdit) (edited *Message, err error) {
	return m.Session.ChannelMessageEditComplex(data)
}

// Delete deletes the message
func (m *Message) Delete() (err error) {
	return m.Session.ChannelMessageDelete(m.ChannelID, m.ID)
}

// Pin pins the message
func (m *Message) Pin() (err error) {
	return m.Session.ChannelMessagePin(m.ChannelID, m.ID)
}

// UnPin unpins the message
func (m *Message) UnPin() (err error) {
	return m.Session.ChannelMessageUnpin(m.ChannelID, m.ID)
}

// AddReaction adds a reaction to the current message
// emoji : the emoji to add
func (m *Message) AddReaction(emoji *Emoji) (err error) {
	return m.Session.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName())
}

// RemoveReaction removes the reaction added by user from the message
// emoji : the emoji to remove
// user : the user or member who added the reaction
func (m *Message) RemoveReaction(emoji *Emoji, user IDGettable) (err error) {
	id := user.GetID()
	if id == m.Session.State.MyUser().ID {
		id = "@me"
	}
	return m.Session.MessageReactionRemove(m.ChannelID, m.ID, emoji.APIName(), user.GetID())
}

// RemoveAllReactions removes all the reactions from the message
func (m *Message) RemoveAllReactions() (err error) {
	return m.Session.MessageReactionsRemoveAll(m.ChannelID, m.ID)
}
