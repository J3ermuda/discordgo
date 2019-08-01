package discordgo

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

type CircularBuffer struct {
	Key    string
	Length int
}

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(redisOptions interface{}) (RedisAdapter, error) {
	options := redisOptions.(redis.Options)

	return RedisAdapter{client: redis.NewClient(&options)}, nil
}

func (a *RedisAdapter) GuildAdd(guild *Guild) error {
	guild.Session = nil
	data, err := json.Marshal(guild)
	if err != nil {
		return err
	}

	_, err = a.client.HSet("cache:guild", guild.ID, string(data)).Result()
	return err
}

func (a *RedisAdapter) GuildRemove(guild *Guild) error {
	_, err := a.client.HDel("cache:guild", guild.ID).Result()
	return err
}

func (a *RedisAdapter) Guild(guildID string) (*Guild, error) {
	data, err := a.client.HGet("cache:guild", guildID).Result()
	if err != nil {
		return nil, err
	}
	g := new(Guild)
	err = json.Unmarshal([]byte(data), g)
	return g, err
}

func (a *RedisAdapter) MemberAdd(member *Member) error {
	data, err := json.Marshal(member)
	if err != nil {
		return err
	}

	_, err = a.client.HSet(fmt.Sprintf("cache:member:%s", member.GuildID), member.User.ID, []byte(data)).Result()
	return err
}

func (a *RedisAdapter) MemberRemove(member *Member) error {
	_, err := a.client.HDel(fmt.Sprintf("cache:member:%s", member.GuildID), member.User.ID).Result()
	return err
}

func (a *RedisAdapter) Member(guildID string, userID string) (*Member, error) {
	data, err := a.client.HGet(fmt.Sprintf("cache:member:%s", guildID), userID).Result()
	if err != nil {
		return nil, err
	}

	m := new(Member)
	err = json.Unmarshal([]byte(data), m)
	return m, err
}

func (a *RedisAdapter) RoleAdd(guildID string, role *Role) error {
	data, err := json.Marshal(role)
	if err != nil {
		return err
	}

	_, err = a.client.HSet(fmt.Sprintf("cache:role:%s", guildID), role.ID, []byte(data)).Result()
	return err
}

func (a *RedisAdapter) RoleRemove(guildID string, role *Role) error {
	_, err := a.client.HDel(fmt.Sprintf("cache:role:%s", guildID), role.ID).Result()
	return err
}

func (a *RedisAdapter) Role(guildID, roleID string) (*Role, error) {
	data, err := a.client.HGet(fmt.Sprintf("cache:role:%s", guildID), roleID).Result()
	if err != nil {
		return nil, err
	}

	r := new(Role)
	err = json.Unmarshal([]byte(data), r)
	return r, err
}

func (a *RedisAdapter) ChannelAdd(guildID string, channel *Channel) error {
	data, err := json.Marshal(channel)
	if err != nil {
		return err
	}

	_, err = a.client.HSet("cache:channel", channel.ID, []byte(data)).Result()
	return err
}

func (a *RedisAdapter) ChannelRemove(channel *Channel) error {
	_, err := a.client.HDel("cache:channel", channel.ID).Result()
	return err
}

func (a *RedisAdapter) Channel(channelID string) (*Channel, error) {
	data, err := a.client.HGet("cache:channel", channelID).Result()
	if err != nil {
		return nil, err
	}

	c := new(Channel)
	err = json.Unmarshal([]byte(data), c)
	return c, err
}

func (a *RedisAdapter) Emoji(guildID, emojiID string) (*Emoji, error) {
	data, err := a.client.HGet(fmt.Sprintf("cache:emoji:%s", guildID), emojiID).Result()
	if err != nil {
		return nil, err
	}

	e := new(Emoji)
	err = json.Unmarshal([]byte(data), e)
	return e, err
}

func (a *RedisAdapter) EmojiAdd(guildID string, emoji *Emoji) error {
	data, err := json.Marshal(emoji)
	if err != nil {
		return err
	}

	_, err = a.client.HSet(fmt.Sprintf("cache:emoji:%s", guildID), emoji.ID, []byte(data)).Result()
	return err
}
