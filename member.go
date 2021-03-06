package discordgo

import (
	"sort"
	"time"
)

// A Member stores user information for Guild members. A guild
// member represents a certain user's presence in a guild.
type Member struct {
	// The guild ID on which the member exists.
	GuildID string `json:"guild_id"`

	// The time at which the member joined the guild, in ISO8601.
	JoinedAt Timestamp `json:"joined_at"`

	// The nickname of the member, if they have one.
	Nick string `json:"nick"`

	// Whether the member is deafened at a guild level.
	Deaf bool `json:"deaf"`

	// Whether the member is muted at a guild level.
	Mute bool `json:"mute"`

	// The underlying user on which the member is based.
	User *User `json:"user"`

	// A list of IDs of the roles which are possessed by the member.
	Roles []string `json:"roles"`

	// When the user used their Nitro boost on the server
	PremiumSince Timestamp `json:"premium_since"`

	// Guild gets set when GetGuild gets called for the first time as an optimisation technique
	guild *Guild
}

// String returns a unique identifier of the form displayName#discriminator
func (m Member) String() string {
	return m.GetDisplayName() + "#" + m.User.Discriminator
}

// GetID returns the members ID
func (m Member) GetID() string {
	return m.User.ID
}

// CreatedAt returns the members creation time in UTC
func (m Member) CreatedAt() (creation time.Time, err error) {
	return m.User.CreatedAt()
}

// Mention creates a member mention
func (m Member) Mention() string {
	if m.Nick != "" {
		return "<@!" + m.User.ID + ">"
	}
	return m.User.Mention()
}

// IsOwner checks if the member is the owner of the guild they are in
func (m *Member) IsOwner() bool {
	g, err := m.GetGuild()
	if err != nil {
		return false
	}

	return g.OwnerID == m.GetID()
}

// IsMentionedIn checks if the member is mentioned in the given message
// message      : message to check for mentions
func (m *Member) IsMentionedIn(message *Message) bool {
	if m.User.IsMentionedIn(message) {
		return true
	}

	rRoles, err := m.GetRoles()
	if err != nil {
		return false
	}
	roles := Roles(rRoles)

	for _, roleID := range message.MentionRoles {
		if roles.ContainsID(roleID) {
			return true
		}
	}

	return false
}

// AvatarURL returns a URL to the user's avatar.
//    size:    The size of the user's avatar as a power of two
//             if size is an empty string, no size parameter will
//             be added to the URL.
func (m *Member) AvatarURL(size string) string {
	return m.User.AvatarURL(size)
}

// GetDisplayName returns the members nick if one has been set and else their username
func (m *Member) GetDisplayName() string {
	if m.Nick != "" {
		return m.Nick
	}
	return m.User.Username
}

// GetGuild returns the guild object where the Member belongs to
func (m *Member) GetGuild() (g *Guild, err error) {
	if m.guild != nil {
		g = m.guild
		return
	}

	g, err = m.User.Session.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	m.guild = g
	return
}

// GetRoles returns a slice with all roles the Member has, sorted from highest to lowest
func (m *Member) GetRoles() (roles []*Role, err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	var base Roles
	for _, roleID := range m.Roles {
		r, errGR := g.GetRole(roleID)
		if errGR != nil {
			err = errGR
			return
		}
		base = append(base, r)
	}
	sort.Sort(base)
	roles = append(roles, base...)
	return
}

// GetColor returns the hex code of the members color as displayed in the server
func (m *Member) GetColor() (color Color, err error) {
	roles, err := m.GetRoles()
	if err != nil {
		return
	}

	for _, role := range roles {
		if role.Color != 0 {
			return role.Color, nil
		}
	}

	return
}

// GetTopRole returns the members highest role
func (m *Member) GetTopRole() (role *Role, err error) {
	roles, err := m.GetRoles()
	if err != nil {
		return
	}

	role = roles[0]
	return
}

// GetVoiceState returns the members voice state
func (m *Member) GetVoiceState() (voice *VoiceState, err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	return g.GetVoiceState(m.GetID())
}

// Kick kicks the member from their guild
// reason   : reason for the kick
func (m *Member) Kick(reason string) (err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	return g.Kick(m.User, reason)
}

// Ban bans the member from their guild
// reason     : reason for the ban as it will be displayed in the audit log
// days       : days of messages to delete
func (m *Member) Ban(reason string, days int) (err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	return g.Ban(m.User, reason, days)
}

// EditRoles replaces all roles of the user with the provided slice of roles
// roles   : a slice of Role objects
// reason  : the reason for the change in roles
func (m *Member) EditRoles(roles Roles, reason string) (err error) {
	roleIDs := make([]string, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	return m.User.Session.GuildMemberEdit(m.GuildID, m.User.ID, reason, roleIDs)
}

// EditNickname sets the nickname of the member
// nick      : the new nickname the member will have
func (m *Member) EditNickname(nick string) (err error) {
	return m.User.Session.GuildMemberNickname(m.GuildID, m.User.ID, nick)
}

// MoveTo moves the member to a voice channel
// channel   : voice channel to move the user to
// reason    : the reason for the move
func (m *Member) MoveTo(channel *Channel, reason string) (err error) {
	if channel.Type != ChannelTypeGuildVoice {
		return ErrNotAVoiceChannel
	}
	return m.User.Session.GuildMemberMove(m.GuildID, m.User.ID, channel.ID, reason)
}

// DisconnectFromVoice disconnects the member from whatever voice channel they are in
// reason    : the reason for the disconnect
func (m *Member) DisconnectFromVoice(reason string) (err error) {
	return m.User.Session.GuildMemberVoiceDisconnect(m.GuildID, m.User.ID, reason)
}

// AddRole adds a role to the member
// role     : role to add
// reason   : the reason for the role add
func (m *Member) AddRole(role *Role, reason string) (err error) {
	return m.User.Session.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID, reason)
}

// RemoveRole removes a role from the member
// role     : role to remove
// reason   : the reason for the role remove
func (m *Member) RemoveRole(role *Role, reason string) (err error) {
	return m.User.Session.GuildMemberRoleRemove(m.GuildID, m.User.ID, role.ID, reason)
}

// RemoveRoles removes multiple roles from the member
// roles     : roles to remove
// reason   : the reason for the role removes
func (m *Member) RemoveRoles(roles []*Role, reason string) (err error) {
	keepRoles := m.Roles[:0]
	for _, role := range m.Roles {
		var remove bool
		for _, toRemove := range roles {
			if toRemove.ID == role && !toRemove.IsDefault() {
				remove = true
				break
			}
		}
		if !remove {
			keepRoles = append(keepRoles, role)
		}
	}

	return m.User.Session.GuildMemberEdit(m.GuildID, m.GetID(), reason, keepRoles)
}

// AddRoles adds multiple roles to the member
// roles     : roles to add
// reason   : the reason for the role adds
func (m *Member) AddRoles(roles []*Role, reason string) (err error) {
	ids := make([]string, 0, len(roles)+len(m.Roles))
	for _, role := range roles {
		ids = append(ids, role.ID)
	}
	for _, role := range m.Roles {
		if !Contains(ids, role) {
			ids = append(ids, role)
		}
	}
	return m.User.Session.GuildMemberEdit(m.GuildID, m.GetID(), reason, ids)
}
