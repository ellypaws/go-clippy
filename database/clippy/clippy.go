package clippy

import (
	"github.com/nokusukun/bingo"
	"strconv"
)

var ClippyCollection *bingo.Collection[Clippy]

type Clippy struct {
	Username  string   `json:"username,omitempty"`
	Snowflake int      `json:"snowflake,omitempty"`
	Points    []Points `json:"points,omitempty"`
}

type Points struct {
	Score   int    `json:"score,omitempty"`
	Guild   string `json:"guild,omitempty"`
	GuildID int    `json:"guild_id,omitempty"`
}

func (c Clippy) Key() []byte {
	return []byte(strconv.Itoa(c.Snowflake))
}

func (c Clippy) Record() {
	id, err := ClippyCollection.Insert(c, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	println("Inserted", id)
}
