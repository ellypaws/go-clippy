package functions

import (
	"fmt"
	"github.com/nokusukun/bingo"
	"log"
	"strings"
)

func (f Function) Record(collection *bingo.Collection[Function]) {
	id, err := collection.Insert(f, bingo.Upsert)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted", id)
}

func RecordMany(f []Function, collection *bingo.Collection[Function]) {
	id, err := collection.InsertMany(f, bingo.Upsert)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range id {
		fmt.Println("Inserted", i)
	}
}

func GetFunction(s string, collection *bingo.Collection[Function]) (*Function, error) {
	q := QueryFunction(s, collection)
	if q.Error != nil {
		return nil, q.Error
	}
	return q.First(), nil
}

func QueryFunction(s string, collection *bingo.Collection[Function]) *bingo.QueryResult[Function] {
	return collection.Query(bingo.Query[Function]{
		Filter: func(doc Function) bool {
			return strings.Contains(strings.ToLower(doc.Name), strings.ToLower(s))
		},
	})
}

func GetByKey(s string, collection *bingo.Collection[Function]) (*Function, error) {
	q := QueryKey(strings.ToUpper(s), collection)
	if q.Error != nil {
		return nil, q.Error
	}
	return q.First(), nil
}

func QueryKey(s string, collection *bingo.Collection[Function]) *bingo.QueryResult[Function] {
	return collection.Query(bingo.Query[Function]{
		Keys: [][]byte{[]byte(s)},
	})
}

var cache Cache

type Cache []Function

func Cached(collection *bingo.Collection[Function]) Cache {
	if cache != nil {
		return cache
	}
	result := collection.Query(bingo.Query[Function]{
		Filter: func(doc Function) bool {
			return true
		},
	})
	if !result.Any() {
		return []Function{}
	}
	var items []Function
	for _, i := range result.Items {
		if i != nil {
			items = append(items, *i)
		}
	}
	return items
}

func (c Cache) String(i int) string {
	switch {
	//case c[i].Syntax.Layout != "":
	//	return c[i].Syntax.Layout
	//case c[i].Syntax.Raw != "":
	//	return c[i].Syntax.Raw
	default:
		return c[i].Name
	}
}

func (c Cache) Len() int {
	return len(c)
}
