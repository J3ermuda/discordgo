package discordgo

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNotACustomEmoji gets thrown when a method gets called on an unicode emoji
	// that can only be called on custom emojis
	ErrNotACustomEmoji = errors.New("you can't do this to a custom emoji")

	// ErrUnknownEmojiGuild gets thrown when the method requires the emoji to be from
	// a guild that is cached, but it isn't
	ErrUnknownEmojiGuild = errors.New("the guild that this emoji comes from is not in the cache")
)

// Emoji struct holds data related to Emoji's
type Emoji struct {
	// The ID of the emoji, this is empty if the emoji is not custom
	ID string `json:"id"`

	// The name of the emoji, this is the unicode character of the emoji if it's not custom
	Name string `json:"name"`

	// A list of roles that is allowed to use this emoji, if it is empty, the emoji is unrestricted.
	Roles []string `json:"roles"`

	// if the emoji is managed by an external service
	Managed bool `json:"managed"`

	// If colons are required to use this emoji in the client
	RequireColons bool `json:"require_colons"`

	// If the emoji is animated
	Animated bool `json:"animated"`

	// The user that created the emoji, his can only be retrieved when fetching the emoji
	User *User `json:"user,omitempty"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"-"`

	// the guild this emoji belongs to
	Guild *Guild `json:"-"`
}

// GetID returns the emoji's ID, this will be an empty string if the emoji is not custom
func (e *Emoji) GetID() string {
	return e.ID
}

// CreatedAt returns the emoji's creation time in UTC,
// will return an error if the emoji is not custom
func (e *Emoji) CreatedAt() (creation time.Time, err error) {
	if !e.IsCustom() {
		err = ErrNotACustomEmoji
		return
	}

	return SnowflakeToTime(e.ID)
}

// IsEqual returns true if the other emoji has the same ID if they are both custom
// and else true if they have the same name
func (e *Emoji) IsEqual(other *Emoji) bool {
	if e.IsCustom() != other.IsCustom() {
		return false
	} else if e.IsCustom() {
		return e.ID == other.ID
	}
	return e.Name == other.Name
}

// IsCustom returns true if the emoji is a custom emoji
func (e *Emoji) IsCustom() bool {
	return e.ID != ""
}

// String renders the string needed to display the emoji correctly in discord
func (e *Emoji) String() string {
	if !e.IsCustom() {
		return e.Name
	} else if e.Animated {
		return fmt.Sprintf("<a:%s:%s>", e.Name, e.ID)
	}
	return fmt.Sprintf("<:%s:%s>", e.Name, e.ID)
}

// MessageFormat returns a correctly formatted Emoji for use in Message content and embeds
func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

// APIName returns an correctly formatted API name for use in the MessageReactions endpoints.
func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

// RoleObjects returns a slice of role objects,
// formed from the slice of strings that is the Roles attribute of Emoji
func (e *Emoji) RoleObjects() (roles []*Role) {
	for _, r := range e.Guild.Roles {
		if Contains(e.Roles, r.ID) {
			roles = append(roles, r)
		}
	}
	return
}

// Delete deletes the emoji
func (e *Emoji) Delete() error {
	if e.ID == "" {
		return ErrNotACustomEmoji
	}

	if e.Guild == nil {
		return ErrUnknownEmojiGuild
	}

	return e.Session.GuildEmojiDelete(e.Guild.ID, e.ID)
}

// EditName edits the name of the custom emoji
// name :  the new name for the custom emoji
func (e *Emoji) EditName(name string) (edited *Emoji, err error) {
	if e.ID == "" {
		err = ErrNotACustomEmoji
		return
	}

	if e.Guild == nil {
		err = ErrUnknownEmojiGuild
		return
	}

	return e.Session.GuildEmojiEdit(e.Guild.ID, e.ID, name, e.Roles)
}

// LimitRoles limits the use of the emoji to the roles given here,
// leave empty to make the emoji unrestricted
// roles :  the list of roles to make the emoji exclusive to
func (e *Emoji) LimitRoles(roles []*Role) (edited *Emoji, err error) {
	if e.ID == "" {
		err = ErrNotACustomEmoji
		return
	}

	if e.Guild == nil {
		err = ErrUnknownEmojiGuild
		return
	}

	var roleIDs []string
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return e.Session.GuildEmojiEdit(e.Guild.ID, e.ID, e.Name, roleIDs)
}
