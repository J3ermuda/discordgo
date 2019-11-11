package discordgo

import "errors"

var (
	// ErrNotATextChannel gets returned when a method gets called on a channel
	// that does not support sending messages to them.
	ErrNotATextChannel = errors.New("not a text or dm channel")

	// ErrNotAVoiceChannel gets returned when a method gets called on a channel
	// that is not a Guild Voice channel but does need to be for the method to work.
	ErrNotAVoiceChannel = errors.New("not a voice channel")

	// ErrNotAGuildChannel gets returned when a method gets called on a channel
	// that is not inside of a Guild but does need to be for the method to work.
	ErrNotAGuildChannel = errors.New("not a channel in a guild")

	// ErrNotOverwriteSettable gets returned when an IDGettable that isn't an user, member or role,
	// gets passed to a function to set or request the permission overwrite for
	ErrNotOverwriteSettable = errors.New("target parameter must be either a user, member or a role")

	// ErrObjectNotFound gets returned when the requested object wasn't found
	ErrObjectNotFound = errors.New("requested object not found")

	// ErrNilState is returned when the state is nil.
	ErrNilState = errors.New("state not instantiated, please use discordgo.New() or assign Session.State")

	// ErrStateNotFound is returned when the state cache
	// requested is not found
	ErrStateNotFound = errors.New("state cache not found")

	// ErrRolePositionBounds gets returned when trying to set the role position to lower than 1
	ErrRolePositionBounds = errors.New("the position cannot be lower than 1")

	// ErrUnmovableDefaultRole gets returned when trying to move the default role as this is impossible
	ErrUnmovableDefaultRole = errors.New("can't move the default role")

	// ErrWSAlreadyOpen is returned when you attempt to open
	// a websocket that already is open.
	ErrWSAlreadyOpen = errors.New("web socket already opened")

	// ErrWSNotFound is returned when you attempt to use a websocket
	// that doesn't exist
	ErrWSNotFound = errors.New("no websocket connection exists")

	// ErrWSShardBounds is returned when you try to use a shard ID that is
	// less than the total shard count
	ErrWSShardBounds = errors.New("ShardID must be less than ShardCount")

	// ErrNotACustomEmoji gets returned when a method gets called on an unicode emoji
	// that can only be called on custom emojis
	ErrNotACustomEmoji = errors.New("you can't do this to a custom emoji")

	// ErrUnknownEmojiGuild gets returned when the method requires the emoji to be from
	// a guild that is cached, but it isn't
	ErrUnknownEmojiGuild = errors.New("the guild that this emoji comes from is not in the cache")

	// ErrJSONUnmarshal gets returned when the JSON unmarshall fails
	ErrJSONUnmarshal = errors.New("json unmarshal")

	// ErrStatusOffline gets returned when trying to set your status to offline
	ErrStatusOffline = errors.New("you can't set your Status to offline")

	// ErrVerificationLevelBounds gets returned when trying to set a verification level that is not between 0 and 3
	ErrVerificationLevelBounds = errors.New("VerificationLevel out of bounds, should be between 0 and 3")

	// ErrPruneDaysBounds gets returned when a prune gets requested of less than 1 day
	ErrPruneDaysBounds = errors.New("the number of days should be more than or equal to 1")

	// ErrGuildNoIcon gets returned when the guild does not have a icon set while requesting the icon image
	ErrGuildNoIcon = errors.New("guild does not have an icon set")

	// ErrGuildNoSplash gets returned when the guild does not have a splash set while requesting the splash image
	ErrGuildNoSplash = errors.New("guild does not have a splash set")

	// ErrUnauthorized gets returned when the HTTP request was unauthorized
	ErrUnauthorized = errors.New("HTTP request was unauthorized. This could be because the provided token was not a bot token")
)
