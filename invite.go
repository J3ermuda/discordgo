package discordgo

// A Invite stores all data related to a specific Discord Guild or Channel invite.
type Invite struct {
	Guild     *Guild    `json:"guild"`
	Channel   *Channel  `json:"channel"`
	Inviter   *User     `json:"inviter"`
	Code      string    `json:"code"`
	CreatedAt Timestamp `json:"created_at"`
	MaxAge    int       `json:"max_age"`
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	Revoked   bool      `json:"revoked"`
	Temporary bool      `json:"temporary"`
	Unique    bool      `json:"unique"`

	TargetUser     *User `json:"target_user"`
	TargetUserType int   `json:"target_user_type"`

	// will only be filled when using InviteWithCounts
	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`
}

func (i *Invite) build(s *Session) {
	guild, GErr := s.State.Guild(i.Guild.ID)
	if GErr == nil {
		i.Guild = guild
	} else {
		i.Guild.Session = s
	}

	user, UErr := s.FetchUser(i.Inviter.ID)
	if UErr == nil {
		i.Inviter = user
	} else {
		i.Inviter.Session = s
	}

	channel, CErr := s.State.Channel(i.Channel.ID)
	if CErr == nil {
		i.Channel = channel
	} else {
		i.Channel.Session = s
	}
}

// Delete deletes the invite
func (i *Invite) Delete() (err error) {
	_, err = i.Guild.DeleteInvite(i.Code)
	return
}

// InviteBuilder is an object used to create an invite
type InviteBuilder struct {
	MaxAge    int  `json:"max_age"`
	MaxUses   int  `json:"max_uses"`
	Temporary bool `json:"temporary"`
	Unique    bool `json:"unique"`
}

// NewInviteBuilder creates a new InviteBuilder to chain with and use as data for creating an invite
func NewInviteBuilder() *InviteBuilder {
	return &InviteBuilder{
		MaxAge: 86400,
	}
}

// SetMaxAge can be used to set the invite max age in a chain
func (i *InviteBuilder) SetMaxAge(age int) *InviteBuilder {
	i.MaxAge = age
	return i
}

// SetMaxUses can be used to set the invite max uses in a chain
func (i *InviteBuilder) SetMaxUses(uses int) *InviteBuilder {
	i.MaxUses = uses
	return i
}

// IsTemporary can be used to set the invite to being temporary in a chain
func (i *InviteBuilder) IsTemporary() *InviteBuilder {
	i.Temporary = true
	return i
}

// IsUnique can be used to set the invite to being unique  in a chain
func (i *InviteBuilder) IsUnique() *InviteBuilder {
	i.Unique = true
	return i
}
