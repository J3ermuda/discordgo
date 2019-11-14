package discordgo

import (
	"strings"
	"time"
)

// A Guild holds all data related to a specific Discord Guild.  Guilds are also
// sometimes referred to as Servers in the Discord client.
type Guild struct {
	// The ID of the guild.
	ID string `json:"id"`

	// The name of the guild. (2–100 characters)
	Name string `json:"name"`

	// The hash of the guild's icon. Use Session.GuildIcon
	// to retrieve the icon itself.
	Icon string `json:"icon"`

	// The voice region of the guild.
	Region string `json:"region"`

	// The ID of the AFK voice channel.
	AfkChannelID string `json:"afk_channel_id"`

	// The ID of the embed channel ID, used for embed widgets.
	EmbedChannelID string `json:"embed_channel_id"`

	// The user ID of the owner of the guild.
	OwnerID string `json:"owner_id"`

	// The time at which the current user joined the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	JoinedAt Timestamp `json:"joined_at"`

	// The hash of the guild's splash.
	Splash string `json:"splash"`

	// The timeout, in seconds, before a user is considered AFK in voice.
	AfkTimeout int `json:"afk_timeout"`

	// The number of members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	MemberCount int `json:"member_count"`

	// The verification level required for the guild.
	VerificationLevel VerificationLevel `json:"verification_level"`

	// Whether the guild has embedding enabled.
	EmbedEnabled bool `json:"embed_enabled"`

	// Whether the guild is considered large. This is
	// determined by a member threshold in the identify packet,
	// and is currently hard-coded at 250 members in the library.
	Large bool `json:"large"`

	// The default message notification setting for the guild.
	// 0 == all messages, 1 == mentions only.
	DefaultMessageNotifications int `json:"default_message_notifications"`

	// A list of roles in the guild.
	Roles []*Role `json:"roles"`

	// A list of the custom emojis present in the guild.
	Emojis []*Emoji `json:"emojis"`

	// A list of the members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Members []*Member `json:"members"`

	// A list of partial presence objects for members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Presences []*Presence `json:"presences"`

	// A list of channels in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Channels []*Channel `json:"channels"`

	// A list of voice states for the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	VoiceStates []*VoiceState `json:"voice_states"`

	// Whether this guild is currently unavailable (most likely due to outage).
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Unavailable bool `json:"unavailable"`

	// The explicit content filter level
	ExplicitContentFilter ExplicitContentFilterLevel `json:"explicit_content_filter"`

	// The list of enabled guild features
	Features []string `json:"features"`

	// Required MFA level for the guild
	MfaLevel MfaLevel `json:"mfa_level"`

	// Whether or not the Server Widget is enabled
	WidgetEnabled bool `json:"widget_enabled"`

	// The Channel ID for the Server Widget
	WidgetChannelID string `json:"widget_channel_id"`

	// The Channel ID to which system messages are sent (eg join and leave messages)
	SystemChannelID string `json:"system_channel_id"`

	// the vanity url code for the guild
	VanityURLCode string `json:"vanity_url_code"`

	// the description for the guild
	Description string `json:"description"`

	// The hash of the guild's banner
	Banner string `json:"banner"`

	// The premium tier of the guild
	PremiumTier PremiumTier `json:"premium_tier"`

	// The total number of users currently boosting this server
	PremiumSubscriptionCount int `json:"premium_subscription_count"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"-"`
}

// A UserGuild holds a brief version of a Guild
type UserGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int    `json:"permissions"`
}

// A GuildParams stores all the data needed to update discord guild settings
type GuildParams struct {
	Name                        string             `json:"name,omitempty"`
	Region                      string             `json:"region,omitempty"`
	VerificationLevel           *VerificationLevel `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                `json:"default_message_notifications,omitempty"` // TODO: Separate type?
	AfkChannelID                string             `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                `json:"afk_timeout,omitempty"`
	Icon                        string             `json:"icon,omitempty"`
	OwnerID                     string             `json:"owner_id,omitempty"`
	Splash                      string             `json:"splash,omitempty"`
}

// GuildLimit stores the information about the limits for the guild that can be boosted
type GuildLimit struct {
	Emoji    Emoji
	Bitrate  int
	FileSize int
}

// String returns the guild's name for easy formatting
func (g Guild) String() string {
	return g.Name
}

// GetID returns the guilds ID
func (g Guild) GetID() string {
	return g.ID
}

// CreatedAt returns the guilds creation time in UTC
func (g Guild) CreatedAt() (creation time.Time, err error) {
	return SnowflakeToTime(g.ID)
}

// GetIconURL returns the URL to the guild icon
func (g *Guild) GetIconURL() string {
	return EndpointGuildIcon(g.ID, g.Icon)
}

// GetRole gets the role with the given ID as it is stored in Guild.Roles
func (g *Guild) GetRole(roleID string) (role *Role, err error) {
	for _, role = range g.Roles {
		if role.ID == roleID {
			return role, nil
		}
	}

	return nil, ErrObjectNotFound
}

// GetRoleNamed gets the role with the given name as it is stored in Guild.Roles
// It is semi-case-sensitive; if a name matches full, the first role with that name gets returned
// if a name matches but with different capitalization, the last role with that name gets returned
func (g *Guild) GetRoleNamed(name string) (role *Role, err error) {
	var savedRole *Role
	lowerCaseName := strings.ToLower(name)

	for _, role = range g.Roles {
		if role.Name == name {
			return
		} else if role.Name == lowerCaseName {
			savedRole = role
		}
	}

	if savedRole != nil {
		role = savedRole
		return
	}

	return nil, ErrObjectNotFound
}

// GetChannel gets the channel with the given ID as it is stored in Guild.Channels
// channelID    : The ID of the channel to search for
func (g *Guild) GetChannel(channelID string) (channel *Channel, err error) {
	for _, channel = range g.Channels {
		if channel.ID == channelID {
			return
		}
	}

	return nil, ErrObjectNotFound
}

// GetChannelNamed gets the channel with the given name as it is stored in Guild.Channels
// It is semi-case-sensitive; if a name matches full, the first channel with that name gets returned
// if a name matches but with different capitalization, the last channel with that name gets returned
// name    : The name of the channel to search for
func (g *Guild) GetChannelNamed(name string) (channel *Channel, err error) {
	var savedChannel *Channel
	lowerCaseName := strings.ToLower(name)

	for _, channel = range g.Channels {
		if channel.Name == name {
			return
		} else if strings.ToLower(channel.Name) == lowerCaseName {
			savedChannel = channel
		}
	}

	if savedChannel != nil {
		channel = savedChannel
		return
	}

	return nil, ErrObjectNotFound
}

// GetMember gets the member with the given ID from the guild.
// userID   : The ID of the member to search for
func (g *Guild) GetMember(userID string) (member *Member, err error) {
	for _, member = range g.Members {
		if member.User.ID == userID {
			return member, nil
		}
	}

	return nil, ErrObjectNotFound
}

// Me retrieves the member object representing this client in the guild
func (g *Guild) Me() (member *Member, err error) {
	return g.GetMember(g.Session.State.MyUser().ID)
}

// FetchMembers fetches count members of this guild from discord and adds them to the state.
// limit   : The max amount of members to fetch (max 1000)
// after   : The id of the member to return members after
// TODO: Make this use the websocket instead of the API
func (g *Guild) FetchMembers(max int, after string) (err error) {
	members, err := g.Session.GuildMembers(g.ID, after, max)
	if err != nil {
		return
	}

	for _, m := range members {
		err = g.Session.State.MemberAdd(m, g.Session)
		if err != nil {
			return
		}
	}

	return nil
}

// GetVoiceState returns the the VoiceState of the member with the given userID
// userID  : userID of the member whose VoiceState should be returned
func (g *Guild) GetVoiceState(userID string) (voice *VoiceState, err error) {
	for _, voice = range g.VoiceStates {
		if voice.UserID == userID {
			return voice, nil
		}
	}

	err = ErrObjectNotFound
	return
}

// NewEdit creates a new GuildParams to chain an edit with
func (g *Guild) NewEdit() *GuildParams {
	return &GuildParams{}
}

// SetName can be used to set the guild name in a chain
func (g *GuildParams) SetName(name string) *GuildParams {
	g.Name = name
	return g
}

// SetRegion can be used to set the guild region in a chain
func (g *GuildParams) SetRegion(region string) *GuildParams {
	g.Region = region
	return g
}

// SetVerificationLevel can be used to set the guild VerificationLevel in a chain
func (g *GuildParams) SetVerificationLevel(level VerificationLevel) *GuildParams {
	g.VerificationLevel = &level
	return g
}

// SetDefaultMessageNotifications can be used to set the guild DefaultMessageNotifications in a chain
func (g *GuildParams) SetDefaultMessageNotifications(setting int) *GuildParams {
	g.DefaultMessageNotifications = setting
	return g
}

// SetAfkChannel can be used to set the guild region in a chain
func (g *GuildParams) SetAfkChannel(channel *Channel) *GuildParams {
	g.AfkChannelID = channel.ID
	return g
}

// SetAfkTimeout can be used to set the guild AfkTimeout in a chain
func (g *GuildParams) SetAfkTimeout(timeout int) *GuildParams {
	g.AfkTimeout = timeout
	return g
}

// SetIcon can be used to set the guild icon in a chain
func (g *GuildParams) SetIcon(icon string) *GuildParams {
	g.Icon = icon
	return g
}

// SetOwner can be used to set the guild owner in a chain
// (this is only possible if the bot is the owner)
func (g *GuildParams) SetOwner(user *User) *GuildParams {
	g.OwnerID = user.ID
	return g
}

// SetSplash can be used to set the guild Splash in a chain
func (g *GuildParams) SetSplash(splash string) *GuildParams {
	g.Splash = splash
	return g
}

// Edit edits the guild with the given GuildParams
// edit  : the GuildParams saying what to edit
func (g *Guild) Edit(edit *GuildParams) (st *Guild, err error) {
	return g.Session.GuildEdit(g.ID, edit)
}

// Ban bans the given user from the guild.
// user      : The User
// reason    : The reason for this ban
// days      : The number of days of previous comments to delete.
func (g *Guild) Ban(user *User, reason string, days int) error {
	return g.Session.GuildBanCreateWithReason(g.ID, user.ID, reason, days)
}

// UnBan unbans the given user
// user    : The User
func (g *Guild) UnBan(user *User) error {
	return g.Session.GuildBanDelete(g.ID, user.ID)
}

// Kick kicks the given user from the guild.
// user      : The User to kick
// reason    : The reason for the kick
func (g *Guild) Kick(user *User, reason string) error {
	return g.Session.GuildMemberDeleteWithReason(g.ID, user.ID, reason)
}

// GetBan returns the ban object of the given user
// user      : The User to get the ban object for
func (g *Guild) GetBan(user *User) (ban *GuildBan, err error) {
	return g.Session.GuildBan(g.ID, user.ID)
}

// GetBans returns an array of GuildBan structures for all bans of the guild
func (g *Guild) GetBans() (bans []*GuildBan, err error) {
	return g.Session.GuildBans(g.ID)
}

// GetInvites returns an array of Invite structures for the guild
func (g *Guild) GetInvites() (invites []*Invite, err error) {
	return g.Session.GuildInvites(g.ID)
}

// AuditLogs returns the audit log of the Guild.
// userID      : If provided the log will be filtered for the given ID.
// beforeID    : If provided all log entries returned will be before the given ID.
// actionType  : If provided the log will be filtered for the given Action Type.
// limit       : The number messages that can be returned. (default 50, min 1, max 100)
func (g *Guild) AuditLogs(userID, beforeID string, actionType, limit int) (log *GuildAuditLog, err error) {
	return g.Session.GuildAuditLog(g.ID, userID, beforeID, actionType, limit)
}

// CreateRole creates and then returns a new Guild Role.
func (g *Guild) CreateRole() (role *Role, err error) {
	return g.Session.GuildRoleCreate(g.ID)
}

// CreateChannel creates and returns a new channel in the guild
// name            : Name of the channel (2-100 chars length)
// channelType     : Type of the channel
func (g *Guild) CreateChannel(name string, channelType ChannelType) (channel *Channel, err error) {
	return g.Session.GuildChannelCreate(g.ID, name, channelType)
}

// EmojiLimit returns the maximum amount of emojis that this guild can have
func (g *Guild) EmojiLimit() int {
	moreEmojis := Contains(g.Features, "MORE_EMOJI")

	switch g.PremiumTier {
	case PremiumTier1:
		if moreEmojis {
			return 200
		}
		return 100
	case PremiumTier2:
		if moreEmojis {
			return 200
		}
		return 150
	case PremiumTier3:
		return 250
	default:
		if moreEmojis {
			return 200
		}
		return 50
	}
}

// BitrateLimit returns the maximum bitrate for voice channels this guild can have
func (g *Guild) BitrateLimit() int {
	switch g.PremiumTier {
	case PremiumTier1:
		return 128000
	case PremiumTier2:
		return 256000
	case PremiumTier3:
		return 384000
	default:
		return 96000
	}
}

// FileSizeLimit returns the maximum number of bytes files can have when uploaded to this guild
func (g *Guild) FileSizeLimit() int {
	switch g.PremiumTier {
	case PremiumTier1:
		return 8388608
	case PremiumTier2:
		return 52428800
	case PremiumTier3:
		return 104857600
	default:
		return 8388608
	}
}

// DeleteInvite deletes the given invite from the guild
func (g *Guild) DeleteInvite(inviteID string) (*Invite, error) {
	return g.Session.InviteDelete(inviteID)
}

// GetDefaultRole gets the @everyone role that everyone has by default
func (g *Guild) GetDefaultRole() *Role {
	for _, role := range g.Roles {
		if role.IsDefault() {
			return role
		}
	}
	return nil
}
