package clippy

import (
	"fmt"
	"github.com/nokusukun/bingo"
)

func (c Clippy) Record() {
	id, err := ClippyCollection.Insert(c, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted", id)
}
