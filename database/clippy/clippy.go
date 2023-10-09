package clippy

import "github.com/nokusukun/bingo"

var ClippyCollection *bingo.Collection[Clippy]

type Clippy struct {
	Username string `json:"username,omitempty" validate:"required,email"`
	Points   int    `json:"points,omitempty" validate:"required,min=3,max=64"`
}

func (c Clippy) Key() []byte {
	return []byte(c.Username)
}

func (c Clippy) Record() {
	id, err := ClippyCollection.Insert(c, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	println("Inserted", id)
}
