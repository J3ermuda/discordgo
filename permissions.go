package discordgo

// PermissionOffset is used to denote how large the bit offset is for a permission
type PermissionOffset uint

// Constants for the different bit offsets of text channel permissions
const (
	PermissionReadMessages = PermissionOffset(1 << (iota + 10))
	PermissionSendMessages
	PermissionSendTTSMessages
	PermissionManageMessages
	PermissionEmbedLinks
	PermissionAttachFiles
	PermissionReadMessageHistory
	PermissionMentionEveryone
	PermissionUseExternalEmojis
)

// Constants for the different bit offsets of voice permissions
const (
	PermissionVoiceConnect = PermissionOffset(1 << (iota + 20))
	PermissionVoiceSpeak
	PermissionVoiceMuteMembers
	PermissionVoiceDeafenMembers
	PermissionVoiceMoveMembers
	PermissionVoiceUseVAD
	PermissionVoicePrioritySpeaker = 1 << (iota + 2)
	PermissionVoiceStream
)

// Constants for general management.
const (
	PermissionChangeNickname = PermissionOffset(1 << (iota + 26))
	PermissionManageNicknames
	PermissionManageRoles
	PermissionManageWebhooks
	PermissionManageEmojis
)

// Constants for the different bit offsets of general permissions
const (
	PermissionCreateInstantInvite = PermissionOffset(1 << iota)
	PermissionKickMembers
	PermissionBanMembers
	PermissionAdministrator
	PermissionManageChannels
	PermissionManageServer
	PermissionAddReactions
	PermissionViewAuditLogs

	permissionAllGeneral = PermissionAdministrator |
		PermissionViewAuditLogs |
		PermissionManageServer |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionCreateInstantInvite |
		PermissionChangeNickname |
		PermissionManageNicknames |
		PermissionManageEmojis |
		PermissionManageWebhooks
	permissionAllText = PermissionReadMessages |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone |
		PermissionUseExternalEmojis
	permissionAllVoice = PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker |
		PermissionVoiceStream
	permissionAllChannel = permissionAllText |
		permissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	permissionAll = permissionAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis |
		PermissionManageNicknames |
		PermissionChangeNickname
)

// Constants for aliases for other permissions
const (
	PermissionManageChannel     = PermissionManageChannels
	PermissionManagePermissions = PermissionManageRoles
)

// Permissions is a type around the int value of a discord permission
type Permissions int

// IsSuperset returns true if the permissions object has the same or fewer permissions as other.
func (p *Permissions) IsSuperset(other *Permissions) bool {
	return (*p | *other) == *p
}

// IsSubset returns true if the permissions object has the same or more permissions as other.
func (p *Permissions) IsSubset(other *Permissions) bool {
	return (*p & *other) == *p
}

// HandleOverwrite returns a Permissions object that has taken the overwrites into account
func (p *Permissions) HandleOverwrite(allow, deny Permissions) Permissions {
	return (*p &^ deny) | allow
}

// Has returns if the Permissions object has a certain permission
func (p *Permissions) Has(perm PermissionOffset) bool {
	return uint(*p)&uint(perm) == uint(perm)
}

// Set sets a permission on the Permissions object
func (p *Permissions) Set(perm PermissionOffset, value bool) {
	if value {
		*p = Permissions(uint(*p) | uint(perm))
	} else {
		*p = Permissions(uint(*p) &^ uint(perm))
	}
}

// ToMap converts the Permissions to a key-value based representation
func (p *Permissions) ToMap() map[string]bool {
	return map[string]bool{
		"ReadMessages":       p.Has(PermissionReadMessages),
		"SendMessages":       p.Has(PermissionSendMessages),
		"SendTTSMessages":    p.Has(PermissionSendTTSMessages),
		"ManageMessages":     p.Has(PermissionManageMessages),
		"EmbedLinks":         p.Has(PermissionEmbedLinks),
		"AttachFiles":        p.Has(PermissionAttachFiles),
		"ReadMessageHistory": p.Has(PermissionReadMessageHistory),
		"MentionEveryone":    p.Has(PermissionMentionEveryone),
		"UseExternalEmojis":  p.Has(PermissionUseExternalEmojis),

		"VoiceConnect":         p.Has(PermissionVoiceConnect),
		"VoiceSpeak":           p.Has(PermissionVoiceSpeak),
		"VoiceMuteMembers":     p.Has(PermissionVoiceMuteMembers),
		"VoiceDeafenMembers":   p.Has(PermissionVoiceDeafenMembers),
		"VoiceMoveMembers":     p.Has(PermissionVoiceMoveMembers),
		"VoiceUseVAD":          p.Has(PermissionVoiceUseVAD),
		"VoicePrioritySpeaker": p.Has(PermissionVoicePrioritySpeaker),
		"VoiceStream":          p.Has(PermissionVoiceStream),

		"CreateInstantInvite": p.Has(PermissionCreateInstantInvite),
		"ManageChannels":      p.Has(PermissionManageChannels),
		"AddReactions":        p.Has(PermissionAddReactions),
		"ManageRoles":         p.Has(PermissionManageRoles),

		"KickMembers":     p.Has(PermissionKickMembers),
		"BanMembers":      p.Has(PermissionBanMembers),
		"ManageServer":    p.Has(PermissionManageServer),
		"Administrator":   p.Has(PermissionAdministrator),
		"ManageWebhooks":  p.Has(PermissionManageWebhooks),
		"ManageEmojis":    p.Has(PermissionManageEmojis),
		"ManageNicknames": p.Has(PermissionManageNicknames),
		"ChangeNickname":  p.Has(PermissionChangeNickname),
		"ViewAuditLogs":   p.Has(PermissionViewAuditLogs),
	}
}

// A PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Deny  Permissions `json:"deny"`
	Allow Permissions `json:"allow"`
}

// Has returns the value that has been set for a given permission in this overwrite
// it will return true for explicit allow, false for explicit deny
// and nil for if nothing has been explicitly set
func (p *PermissionOverwrite) Has(perm PermissionOffset) *bool {
	if p.Allow.Has(perm) {
		x := true
		return &x
	} else if p.Deny.Has(perm) {
		x := false
		return &x
	} else {
		return nil
	}
}

// Set sets the permission overwrite value, the default value is not true or false, but nil
func (p *PermissionOverwrite) Set(perm PermissionOffset, value *bool) {
	if value == nil {
		p.Deny.Set(perm, false)
		p.Allow.Set(perm, false)
	} else if *value {
		p.Deny.Set(perm, false)
		p.Allow.Set(perm, true)
	} else {
		p.Deny.Set(perm, true)
		p.Allow.Set(perm, false)
	}
}

// ToMap converts the PermissionOverwrite to a key-value based representation of the overwrite
func (p *PermissionOverwrite) ToMap() map[string]*bool {
	return map[string]*bool{
		"ReadMessages":       p.Has(PermissionReadMessages),
		"SendMessages":       p.Has(PermissionSendMessages),
		"SendTTSMessages":    p.Has(PermissionSendTTSMessages),
		"ManageMessages":     p.Has(PermissionManageMessages),
		"EmbedLinks":         p.Has(PermissionEmbedLinks),
		"AttachFiles":        p.Has(PermissionAttachFiles),
		"ReadMessageHistory": p.Has(PermissionReadMessageHistory),
		"MentionEveryone":    p.Has(PermissionMentionEveryone),
		"UseExternalEmojis":  p.Has(PermissionUseExternalEmojis),

		"VoiceConnect":         p.Has(PermissionVoiceConnect),
		"VoiceSpeak":           p.Has(PermissionVoiceSpeak),
		"VoiceMuteMembers":     p.Has(PermissionVoiceMuteMembers),
		"VoiceDeafenMembers":   p.Has(PermissionVoiceDeafenMembers),
		"VoiceMoveMembers":     p.Has(PermissionVoiceMoveMembers),
		"VoiceUseVAD":          p.Has(PermissionVoiceUseVAD),
		"VoicePrioritySpeaker": p.Has(PermissionVoicePrioritySpeaker),
		"VoiceStream":          p.Has(PermissionVoiceStream),

		"CreateInstantInvite": p.Has(PermissionCreateInstantInvite),
		"ManageChannel":       p.Has(PermissionManageChannel),
		"AddReactions":        p.Has(PermissionAddReactions),
		"ManagePermissions":   p.Has(PermissionManagePermissions),
		"ManageWebhooks":      p.Has(PermissionManageWebhooks),
	}
}

// NewAllPermissions is a factory function to create a Permissions object
// with all permissions set to true
func NewAllPermissions() Permissions {
	return Permissions(permissionAll)
}

// NewAllChannelPermissions is a factory function to create a Permissions object
// with all channel-specific permissions set to true
func NewAllChannelPermissions() Permissions {
	return Permissions(permissionAllChannel)
}

// NewTextPermissions is a factory function to create a Permissions object
// with all text permissions from the official Discord UI set to true
func NewTextPermissions() Permissions {
	return Permissions(permissionAllText)
}

// NewVoicePermissions is a factory function to create a Permissions object
// with all voice permissions from the official Discord UI set to true
func NewVoicePermissions() Permissions {
	return Permissions(permissionAllVoice)
}

// NewGeneralPermissions is a factory function to create a Permissions object
// with all general permissions from the official Discord UI set to true
func NewGeneralPermissions() Permissions {
	return Permissions(permissionAllGeneral)
}

func (p PermissionOffset) String() string {
	switch p {
	case permissionAll:
		return "All permissions"
	case permissionAllText:
		return "All Text Permissions"
	case permissionAllVoice:
		return "All Voice Permissions"
	case permissionAllChannel:
		return "All Channel Permissions"
	case permissionAllGeneral:
		return "All General Permissions"

	case PermissionReadMessages:
		return "Read Messages"
	case PermissionSendMessages:
		return "Send Messages"
	case PermissionSendTTSMessages:
		return "Send TTS Messages"
	case PermissionManageMessages:
		return "Manage Messages"
	case PermissionEmbedLinks:
		return "Embed Links"
	case PermissionAttachFiles:
		return "Attach Files"
	case PermissionReadMessageHistory:
		return "Read Message History"
	case PermissionMentionEveryone:
		return "Mention Everyone"
	case PermissionUseExternalEmojis:
		return "Use External Emojis"
	case PermissionAddReactions:
		return "Add Reactions"

	case PermissionVoiceConnect:
		return "Connect"
	case PermissionVoiceSpeak:
		return "Speak"
	case PermissionVoiceMuteMembers:
		return "Mute Members"
	case PermissionVoiceDeafenMembers:
		return "Deafen Members"
	case PermissionVoiceMoveMembers:
		return "Move Members"
	case PermissionVoiceUseVAD:
		return "Use Voice Activity"
	case PermissionVoicePrioritySpeaker:
		return "Priority Speaker"
	case PermissionVoiceStream:
		return "Go Live"

	case PermissionCreateInstantInvite:
		return "Create Invite"
	case PermissionManageChannels:
		return "Manage Channels"
	case PermissionManageChannel:
		return "Manage Channel"
	case PermissionManageRoles:
		return "Manage Roles"
	case PermissionManagePermissions:
		return "Manage Permissions"
	case PermissionKickMembers:
		return "Kick Members"
	case PermissionBanMembers:
		return "Ban Members"
	case PermissionManageServer:
		return "Manage Server"
	case PermissionAdministrator:
		return "Administrator"
	case PermissionManageWebhooks:
		return "Manage Webhooks"
	case PermissionManageEmojis:
		return "Manage Emojis"
	case PermissionManageNicknames:
		return "Manage Nicknames"
	case PermissionChangeNickname:
		return "Change Nickname"
	case PermissionViewAuditLogs:
		return "View Audit Logs"
	default:
		return "Unknown"
	}
}
