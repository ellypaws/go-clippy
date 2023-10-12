package clippy

import (
	"fmt"
	"github.com/nokusukun/bingo"
)

func (c Clippy) Record() {
	id, err := Collection.Insert(c, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted", id)
}

func RecordMany(c []Clippy) {
	id, err := Collection.InsertMany(c, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	for _, i := range id {
		fmt.Println("Inserted", i)
	}
}

func GetClippy(s string) (*Clippy, error) {
	q := QueryClippy(s)
	if q.Error != nil {
		return nil, q.Error
	}
	return q.First(), nil
}

func QueryClippy(s string) *bingo.QueryResult[Clippy] {
	return Collection.Query(bingo.Query[Clippy]{
		Filter: func(doc Clippy) bool {
			return doc.Snowflake == s
		},
	})
}
