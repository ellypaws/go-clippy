package clippy

import "github.com/nokusukun/bingo"

var Collection *bingo.Collection[Award]
var UserSettings *bingo.Collection[User]

type Award struct {
	Username        string `json:"username,omitempty"`
	Snowflake       string `json:"snowflake,omitempty"`
	GuildName       string `json:"guild_name,omitempty"`
	GuildID         string `json:"guild_id,omitempty"`
	Channel         string `json:"channel,omitempty"`
	ChannelID       string `json:"channel_id,omitempty"`
	MessageID       string `json:"message_id,omitempty"`
	OriginUsername  string `json:"origin_username,omitempty"`
	OriginSnowflake string `json:"origin_snowflake,omitempty"`
	InteractionID   string `json:"interaction_id,omitempty"`
}

type User struct {
	Username  string `json:"username,omitempty"`
	Snowflake string `json:"snowflake,omitempty"`
	OptOut    bool   `json:"opt_out,omitempty"` // Assume that Private is true if OptOut is true
	Private   bool   `json:"hide_points,omitempty"`
	Points    int    `json:"points,omitempty"`
}

func (point Award) Key() []byte {
	// nil because this not a keyed slice/collection
	//return nil
	//return []byte(point.Snowflake + point.MessageID)
	return []byte(point.InteractionID)
}

func (config User) Key() []byte {
	return []byte(config.Snowflake)
}
