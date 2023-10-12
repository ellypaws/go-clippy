package clippy

import "github.com/nokusukun/bingo"

var ClippyCollection *bingo.Collection[Clippy]

type Clippy struct {
	Username  string   `json:"username,omitempty"`
	Snowflake string   `json:"snowflake,omitempty"`
	Points    []Points `json:"points,omitempty"`
}

type Points struct {
	Awards    []Award `json:"score,omitempty"`
	GuildName string  `json:"guild_name,omitempty"`
	GuildID   string  `json:"guild_id,omitempty"`
}

type Award struct {
	Channel         string `json:"channel,omitempty"`
	ChannelID       string `json:"channel_id,omitempty"`
	MessageID       string `json:"message_id,omitempty"`
	CallerSnowflake string `json:"caller_snowflake,omitempty"`
	CallerUsername  string `json:"caller_username,omitempty"`
}

func (c Clippy) Key() []byte {
	return []byte(c.Snowflake)
}
