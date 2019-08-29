package discordgo

import (
	"errors"
	"fmt"
	"time"
)

// ErrNotATextChannel gets returned when a method gets called on a channel
// that does not support sending messages to them
var ErrNotATextChannel = errors.New("not a text or dm channel")

// ErrNotAVoiceChannel gets thrown when a method gets called on a channel
// that is not a Guild Voice channel but does need to be for the method to work
var ErrNotAVoiceChannel = errors.New("not a voice channel")

// ErrNotAGuildChannel gets thrown when a method gets called on a channel
// that is not inside of a Guild but does need to be for the method to work
var ErrNotAGuildChannel = errors.New("not a channel in a guild")

// ChannelType is the type of a Channel
type ChannelType int

// Block contains known ChannelType values
const (
	ChannelTypeGuildText ChannelType = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildNews
	ChannelTypeGuildStore
)

// A Channel holds all data related to an individual Discord channel.
type Channel struct {
	// The ID of the channel.
	ID string `json:"id"`

	// The ID of the guild to which the channel belongs, if it is in a guild.
	// Else, this ID is empty (e.g. DM channels).
	GuildID string `json:"guild_id"`

	// The name of the channel.
	Name string `json:"name"`

	// The topic of the channel.
	Topic string `json:"topic"`

	// The type of the channel.
	Type ChannelType `json:"type"`

	// The ID of the last message sent in the channel. This is not
	// guaranteed to be an ID of a valid message.
	LastMessageID string `json:"last_message_id"`

	// The timestamp of the last pinned message in the channel.
	// Empty if the channel has no pinned messages.
	LastPinTimestamp Timestamp `json:"last_pin_timestamp"`

	// Whether the channel is marked as NSFW.
	NSFW bool `json:"nsfw"`

	// Icon of the group DM channel.
	Icon string `json:"icon"`

	// The position of the channel, used for sorting in client.
	Position int `json:"position"`

	// The bitrate of the channel, if it is a voice channel.
	Bitrate int `json:"bitrate"`

	// The recipients of the channel. This is only populated in DM channels.
	Recipients []*User `json:"recipients"`

	// The messages in the channel. This is only present in state-cached channels,
	// and State.MaxMessageCount must be non-zero.
	Messages []*Message `json:"-"`

	// A list of permission overwrites present for the channel.
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`

	// The user limit of the voice channel.
	UserLimit int `json:"user_limit"`

	// The ID of the parent channel, if the channel is under a category
	ParentID string `json:"parent_id"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"-"`
}

// Mention returns a string which mentions the channel
func (c Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

// GetID returns the ID of the channel
func (c Channel) GetID() string {
	return c.ID
}

// CreatedAt returns the channels creation time in UTC
func (c Channel) CreatedAt() (creation time.Time, err error) {
	return SnowflakeToTime(c.ID)
}

// Guild retrieves the guild belonging to the channel
func (c *Channel) Guild() (g *Guild, err error) {
	if c.GuildID == "" {
		err = ErrNotAGuildChannel
		return
	}

	return c.Session.State.Guild(c.GuildID)
}

// A ChannelEdit holds Channel Field data for a channel edit.
type ChannelEdit struct {
	Name                 string                 `json:"name,omitempty"`
	Topic                string                 `json:"topic,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
	Position             int                    `json:"position"`
	Bitrate              int                    `json:"bitrate,omitempty"`
	UserLimit            int                    `json:"user_limit,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
}

// SendMessage sends a message to the channel
// content         : message content to send if provided
// embed           : embed to attach to the message if provided
// files           : files to attach to the message if provided
func (c Channel) SendMessage(content string, embed *MessageEmbed, files []*File) (message *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	data := &MessageSend{
		Content: content,
		Embed:   embed,
		Files:   files,
	}

	return c.SendMessageComplex(data)
}

// SendMessageComplex sends a message to the channel
// data          : MessageSend object with the data to send
func (c Channel) SendMessageComplex(data *MessageSend) (message *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	return c.Session.ChannelMessageSendComplex(c.ID, data)
}

// EditMessage edits an existing message, replacing it entirely with
// the given MessageEdit struct
func (c Channel) EditMessage(data *MessageEdit) (edited *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	data.Channel = c.ID
	return c.Session.ChannelMessageEditComplex(data)
}

// FetchMessage fetches a message with the given ID from the channel
// ID        : ID of the message to fetch
func (c Channel) FetchMessage(ID string) (message *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	return c.Session.ChannelMessage(c.ID, ID)
}

// GetHistory fetches up to limit messages from the channel
// limit     : The number messages that can be returned. (max 100)
// beforeID  : If provided all messages returned will be before given ID.
// afterID   : If provided all messages returned will be after given ID.
// aroundID  : If provided all messages returned will be around given ID.
func (c Channel) GetHistory(limit int, beforeID, afterID, aroundID string) (st []*Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	return c.Session.ChannelMessages(c.ID, limit, beforeID, afterID, aroundID)
}

// HasPins returns a bool indicating if a channel has pinned messages
func (c *Channel) HasPins() bool {
	return c.LastPinTimestamp != ""
}

// FetchPins fetches all pinned messages in the channel from the discord api
func (c *Channel) FetchPins() ([]*Message, error) {
	return c.Session.ChannelMessagesPinned(c.ID)
}

// DeleteMessage deletes a message from the channel
// message        : message to delete
func (c *Channel) DeleteMessage(message *Message) (err error) {
	return c.Session.ChannelMessageDelete(c.ID, message.ID)
}

// DeleteMessageByID deletes a message with the given ID from the channel
// ID        : ID of the message to delete
func (c *Channel) DeleteMessageByID(ID string) (err error) {
	return c.Session.ChannelMessageDelete(c.ID, ID)
}

// MessagesBulkDelete bulk deletes the messages from the channel for the provided message objects.
// messages  : The messages to be deleted. A slice of message objects. A maximum of 100 messages.
func (c *Channel) MessagesBulkDelete(messages []*Message) (err error) {
	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		err = messages[0].Delete()
		return
	}

	if len(messages) > 100 {
		messages = messages[:100]
	}

	twoWeeks := time.Now().Add(-(time.Hour * 24 * 14))
	var toDelete []string
	var tooOld []*Message

	for _, message := range messages {
		age, _ := message.CreatedAt()
		if age.Before(twoWeeks) {
			tooOld = append(tooOld, message)
		} else {
			toDelete = append(toDelete, message.ID)
		}
	}

	err = c.MessagesBulkDeleteByID(toDelete)
	if err != nil {
		return
	}

	for _, oldMessage := range tooOld {
		err = oldMessage.Delete()
		if err != nil {
			return
		}
	}
	return
}

// MessagesBulkDeleteByID bulk deletes the messages from the channel for the provided messageIDs.
// If only one messageID is in the slice call channelMessageDelete function.
// If the slice is empty do nothing.
// messages  : The IDs of the messages to be deleted. A slice of string IDs. A maximum of 100 messages.
func (c *Channel) MessagesBulkDeleteByID(messages []string) (err error) {

	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		err = c.Session.ChannelMessageDelete(c.ID, messages[0])
		return
	}

	if len(messages) > 100 {
		messages = messages[:100]
	}

	data := struct {
		Messages []string `json:"messages"`
	}{messages}

	_, err = c.Session.RequestWithBucketID("POST", EndpointChannelMessagesBulkDelete(c.ID), data, EndpointChannelMessagesBulkDelete(c.ID))
	return
}

// PermissionsFor calculates the permissions that the member has in the channel,
// taking into account the following cases:
// - Guild owner
// - Guild roles
// - Channel overrides
// - Member overrides
func (c *Channel) PermissionsFor(m *Member) (perms Permissions, err error) {
	g, err := c.Guild()
	if err != nil {
		return
	}

	if g.OwnerID == m.GetID() {
		perms = NewAllPermissions()
		return
	}

	perms = *g.GetDefaultRole().Permissions
	roles, err := m.GetRoles()
	if err != nil {
		return
	}

	for _, role := range roles {
		perms = perms | *role.Permissions
	}

	if perms.Has(PermissionAdministrator) {
		perms = NewAllPermissions()
		return
	}

	var remainingOverwrites []*PermissionOverwrite
	if len(c.PermissionOverwrites) > 0 {
		maybeEveryone := c.PermissionOverwrites[0]
		if maybeEveryone.ID == g.ID {
			perms.HandleOverwrite(*maybeEveryone.Allow, *maybeEveryone.Deny)
			if len(c.PermissionOverwrites) > 1 {
				remainingOverwrites = c.PermissionOverwrites[1:]
			}
		} else {
			remainingOverwrites = c.PermissionOverwrites
		}
	}

	denies := Permissions(0)
	allows := Permissions(0)

	for _, overwrite := range remainingOverwrites {
		if overwrite.Type == "role" && Contains(m.Roles, overwrite.ID) {
			denies = denies | *overwrite.Deny
			allows = allows | *overwrite.Allow
		}
	}

	perms.HandleOverwrite(allows, denies)

	for _, overwrite := range remainingOverwrites {
		if overwrite.Type == "member" && overwrite.ID == m.GetID() {
			perms.HandleOverwrite(*overwrite.Allow, *overwrite.Deny)
			break
		}
	}

	if !perms.Has(PermissionSendMessages) {
		perms.Set(PermissionSendTTSMessages, false)
		perms.Set(PermissionMentionEveryone, false)
		perms.Set(PermissionEmbedLinks, false)
		perms.Set(PermissionAttachFiles, false)
	}

	if !perms.Has(PermissionReadMessages) {
		perms = perms &^ NewAllChannelPermissions()
	}

	return
}

func (c *Channel) SetPermissions(target IDGettable, overwrite *PermissionOverwrite) (err error) {
	var permType string
	switch target.(type) {
	case User:
		permType = "member"
	case Member:
		permType = "member"
	case Role:
		permType = "role"
	default:
		err = errors.New("target parameter must be either a user, member or a role")
	}

	return c.Session.ChannelPermissionSet(c.ID, target.GetID(), permType, int(*overwrite.Allow), int(*overwrite.Deny))
}

// Delete deletes the channel
func (c *Channel) Delete() (err error) {
	_, err = c.Session.ChannelDelete(c.ID)
	return
}

// CreateInvite creates an invite
// TODO: make a special object to create invites with
func (c *Channel) CreateInvite(data Invite) (i *Invite, err error) {
	return c.Session.ChannelInviteCreate(c.ID, data)
}
