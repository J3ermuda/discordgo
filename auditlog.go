package discordgo

import (
	"encoding/json"
	"time"
)

// AuditLogChanges contains the data before or after the change occurred that made the audit log entry
// Most data will not be filled all the time and you should always be aware of what is being changed before accessing this object
type AuditLogChanges struct {
	VerificationLevel           VerificationLevel
	ExplicitContentFilterLevel  ExplicitContentFilterLevel
	Allow                       Permissions
	Deny                        Permissions
	Permissions                 Permissions
	ID                          string
	Color                       Color
	Owner                       *User
	Inviter                     *User
	Channel                     *Channel
	AFKChannel                  *Channel
	WidgetChannel               *Channel
	PermissionOverwrites        []*PermissionOverwrite
	Splash                      string
	Icon                        string
	Avatar                      string
	RateLimitPerUser            int
	DefaultMessageNotifications int
	GuildName                   string
	GuildRegion                 string
	AFKTimeout                  int
	MfaLevel                    MfaLevel
	VanityURL                   string
	PruneDeleteDays             int
	WidgetEnabled               bool
	Position                    int
	Topic                       string
	Bitrate                     int
	NSFW                        bool
	ApplicationID               string
	Hoist                       bool
	Mentionable                 bool
	InviteCode                  string
	MaxUses                     int
	Uses                        int
	MaxAge                      int
	Temporary                   bool
	Deaf                        bool
	Mute                        bool
	Nickname                    string
	Type                        string
	RolesAdded                  []*Role
	RolesRemoved                []*Role
}

// AuditLogEntry describes an entry in the guild audit log
type AuditLogEntry struct {
	TargetID   string `json:"target_id"`
	RawChanges []struct {
		NewValue interface{} `json:"new_value"`
		OldValue interface{} `json:"old_value"`
		Key      string      `json:"key"`
	} `json:"changes,omitempty"`
	UserID     string `json:"user_id"`
	ID         string `json:"id"`
	ActionType int    `json:"action_type"`
	Options    struct {
		DeleteMembersDay string `json:"delete_member_days"`
		MembersRemoved   string `json:"members_removed"`
		ChannelID        string `json:"channel_id"`
		Count            string `json:"count"`
		ID               string `json:"id"`
		Type             string `json:"type"`
		RoleName         string `json:"role_name"`
	} `json:"options,omitempty"`
	Reason  string   `json:"reason"`
	Session *Session `json:"-"`
	GuildID string   `json:"-"`
}

// A GuildAuditLog stores data for a guild audit log.
type GuildAuditLog struct {
	Webhooks []struct {
		ChannelID string `json:"channel_id"`
		GuildID   string `json:"guild_id"`
		ID        string `json:"id"`
		Avatar    string `json:"avatar"`
		Name      string `json:"name"`
	} `json:"webhooks,omitempty"`
	Users []struct {
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Bot           bool   `json:"bot"`
		ID            string `json:"id"`
		Avatar        string `json:"avatar"`
	} `json:"users,omitempty"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
}

// Block contains Discord Audit Log Action Types
const (
	AuditLogActionGuildUpdate = 1

	AuditLogActionChannelCreate          = 10
	AuditLogActionChannelUpdate          = 11
	AuditLogActionChannelDelete          = 12
	AuditLogActionChannelOverwriteCreate = 13
	AuditLogActionChannelOverwriteUpdate = 14
	AuditLogActionChannelOverwriteDelete = 15

	AuditLogActionMemberKick       = 20
	AuditLogActionMemberPrune      = 21
	AuditLogActionMemberBanAdd     = 22
	AuditLogActionMemberBanRemove  = 23
	AuditLogActionMemberUpdate     = 24
	AuditLogActionMemberRoleUpdate = 25

	AuditLogActionRoleCreate = 30
	AuditLogActionRoleUpdate = 31
	AuditLogActionRoleDelete = 32

	AuditLogActionInviteCreate = 40
	AuditLogActionInviteUpdate = 41
	AuditLogActionInviteDelete = 42

	AuditLogActionWebhookCreate = 50
	AuditLogActionWebhookUpdate = 51
	AuditLogActionWebhookDelete = 52

	AuditLogActionEmojiCreate = 60
	AuditLogActionEmojiUpdate = 61
	AuditLogActionEmojiDelete = 62

	AuditLogActionMessageDelete = 72
)

// GetID returns the audit log entry's ID
func (e AuditLogEntry) GetID() string {
	return e.ID
}

// CreatedAt returns the audit log entry's creation time
func (e AuditLogEntry) CreatedAt() (creation time.Time, err error) {
	return SnowflakeToTime(e.ID)
}

// Changes returns a before and after AuditLogChanges for the changes in the AuditLogEntry
func (e *AuditLogEntry) Changes() (before, after *AuditLogChanges) {
	before = &AuditLogChanges{}
	after = &AuditLogChanges{}

	for _, change := range e.RawChanges {
		switch change.Key {
		case "name":
			{
				before.GuildName = change.OldValue.(string)
				after.GuildName = change.NewValue.(string)
			}
		case "icon_hash":
			{
				before.Icon = change.OldValue.(string)
				after.Icon = change.NewValue.(string)
			}
		case "splash_hash":
			{
				before.Splash = change.OldValue.(string)
				after.Splash = change.NewValue.(string)
			}
		case "owner_id":
			{
				before.Owner, _ = e.Session.State.GetUser(change.OldValue.(string))
				after.Owner, _ = e.Session.State.GetUser(change.NewValue.(string))
			}
		case "region":
			{
				before.GuildRegion = change.OldValue.(string)
				after.GuildRegion = change.NewValue.(string)
			}
		case "afk_channel_id":
			{
				before.AFKChannel, _ = e.Session.State.Channel(change.OldValue.(string))
				after.AFKChannel, _ = e.Session.State.Channel(change.NewValue.(string))
			}
		case "afk_timeout":
			{
				before.AFKTimeout = change.OldValue.(int)
				after.AFKTimeout = change.NewValue.(int)
			}
		case "mfa_level":
			{
				before.MfaLevel = MfaLevel(change.OldValue.(int))
				after.MfaLevel = MfaLevel(change.NewValue.(int))
			}
		case "verification_level":
			{
				before.VerificationLevel = VerificationLevel(change.OldValue.(int))
				after.VerificationLevel = VerificationLevel(change.NewValue.(int))
			}
		case "explicit_content_filter":
			{
				before.ExplicitContentFilterLevel = ExplicitContentFilterLevel(change.OldValue.(int))
				after.ExplicitContentFilterLevel = ExplicitContentFilterLevel(change.NewValue.(int))
			}
		case "default_message_notifications":
			{
				before.DefaultMessageNotifications = change.OldValue.(int)
				after.DefaultMessageNotifications = change.NewValue.(int)
			}
		case "vanity_url_code":
			{
				before.VanityURL = change.OldValue.(string)
				after.VanityURL = change.NewValue.(string)
			}
		case "prune_delete_days":
			{
				before.PruneDeleteDays = change.OldValue.(int)
				after.PruneDeleteDays = change.NewValue.(int)
			}
		case "widget_enabled":
			{
				before.WidgetEnabled = change.OldValue.(bool)
				after.WidgetEnabled = change.NewValue.(bool)
			}
		case "widget_channel_id":
			{
				before.WidgetChannel, _ = e.Session.State.Channel(change.OldValue.(string))
				after.WidgetChannel, _ = e.Session.State.Channel(change.NewValue.(string))
			}
		case "position":
			{
				before.Position = change.OldValue.(int)
				after.Position = change.NewValue.(int)
			}
		case "topic":
			{
				before.Topic = change.OldValue.(string)
				after.Topic = change.NewValue.(string)
			}
		case "bitrate":
			{
				before.Bitrate = change.OldValue.(int)
				after.Bitrate = change.NewValue.(int)
			}
		case "permission_overwrites":
			{
				o, err := json.Marshal(change.OldValue)
				if err != nil {
					continue
				}
				n, err := json.Marshal(change.NewValue)
				if err != nil {
					continue
				}
				_ = json.Unmarshal(o, &before.PermissionOverwrites)
				_ = json.Unmarshal(n, &after.PermissionOverwrites)
			}
		case "nsfw":
			{
				before.NSFW = change.OldValue.(bool)
				after.NSFW = change.NewValue.(bool)
			}
		case "application_id":
			{
				before.ApplicationID = change.OldValue.(string)
				after.ApplicationID = change.NewValue.(string)
			}
		case "permissions":
			{
				before.Permissions = Permissions(change.OldValue.(int))
				after.Permissions = Permissions(change.NewValue.(int))
			}
		case "color":
			{
				before.Color = Color(change.OldValue.(int))
				after.Color = Color(change.NewValue.(int))
			}
		case "hoist":
			{
				before.Hoist = change.OldValue.(bool)
				after.Hoist = change.NewValue.(bool)
			}
		case "mentionable":
			{
				before.Mentionable = change.OldValue.(bool)
				after.Mentionable = change.NewValue.(bool)
			}
		case "allow":
			{
				before.Allow = Permissions(change.OldValue.(int))
				after.Allow = Permissions(change.NewValue.(int))
			}
		case "deny":
			{
				before.Deny = Permissions(change.OldValue.(int))
				after.Deny = Permissions(change.NewValue.(int))
			}
		case "code":
			{
				before.InviteCode = change.OldValue.(string)
				after.InviteCode = change.NewValue.(string)
			}
		case "channel_id":
			{
				before.Channel, _ = e.Session.State.Channel(change.OldValue.(string))
				after.Channel, _ = e.Session.State.Channel(change.NewValue.(string))
			}
		case "inviter_id":
			{
				before.Owner, _ = e.Session.State.GetUser(change.OldValue.(string))
				after.Owner, _ = e.Session.State.GetUser(change.NewValue.(string))
			}
		case "max_uses":
			{
				before.MaxUses = change.OldValue.(int)
				after.MaxUses = change.NewValue.(int)
			}
		case "uses":
			{
				before.Uses = change.OldValue.(int)
				after.Uses = change.NewValue.(int)
			}
		case "max_age":
			{
				before.MaxAge = change.OldValue.(int)
				after.MaxAge = change.NewValue.(int)
			}
		case "temporary":
			{
				before.Temporary = change.OldValue.(bool)
				after.Temporary = change.NewValue.(bool)
			}
		case "deaf":
			{
				before.Deaf = change.OldValue.(bool)
				after.Deaf = change.NewValue.(bool)
			}
		case "mute":
			{
				before.Mute = change.OldValue.(bool)
				after.Mute = change.NewValue.(bool)
			}
		case "nick":
			{
				before.Nickname = change.OldValue.(string)
				after.Nickname = change.NewValue.(string)
			}
		case "avatar_hash":
			{
				before.Avatar = change.OldValue.(string)
				after.Avatar = change.NewValue.(string)
			}
		case "id":
			{
				before.ID = change.OldValue.(string)
				after.ID = change.NewValue.(string)
			}
		case "type":
			{
				before.Type = change.OldValue.(string)
				after.Type = change.NewValue.(string)
			}
		case "$add":
			{
				g, err := e.Session.State.Guild(e.GuildID)
				if err != nil {
					continue
				}
				o, err := json.Marshal(change.OldValue)
				if err != nil {
					continue
				}
				n, err := json.Marshal(change.NewValue)
				if err != nil {
					continue
				}
				_ = json.Unmarshal(o, &before.RolesAdded)
				_ = json.Unmarshal(n, &after.RolesAdded)

				for _, r := range before.RolesAdded {
					r.Session = e.Session
					r.Guild = g
				}
				for _, r := range after.RolesAdded {
					r.Session = e.Session
					r.Guild = g
				}
			}
		case "$remove":
			{
				g, err := e.Session.State.Guild(e.GuildID)
				if err != nil {
					continue
				}
				o, err := json.Marshal(change.OldValue)
				if err != nil {
					continue
				}
				n, err := json.Marshal(change.NewValue)
				if err != nil {
					continue
				}
				_ = json.Unmarshal(o, &before.RolesRemoved)
				_ = json.Unmarshal(n, &after.RolesRemoved)

				for _, r := range before.RolesRemoved {
					r.Session = e.Session
					r.Guild = g
				}
				for _, r := range after.RolesRemoved {
					r.Session = e.Session
					r.Guild = g
				}
			}
		default:
			continue
		}
	}
	return
}
