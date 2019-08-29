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
		PermissionVoiceUseVAD
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

// Permission is a type around the int value of a discord permission
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
	return (int(*p)>>uint(perm))&1 == 1
}

// Set sets a permission on the Permissions object
func (p *Permissions) Set(perm PermissionOffset, value bool) {
	if value {
		*p = Permissions(int(*p) | (1 << uint(perm)))
	} else {
		*p = Permissions(int(*p) &^ (1 << uint(perm)))
	}
}

// A PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string       `json:"id"`
	Type  string       `json:"type"`
	Deny  *Permissions `json:"deny"`
	Allow *Permissions `json:"allow"`
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
