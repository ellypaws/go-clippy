package bingo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nokusukun/bingo"
	. "go-clippy/database/clippy"
	. "go-clippy/database/functions"
	"net/http"
)

var db *bingo.Driver

func Init() {
	db = getDriver()
	ClippyCollection = bingo.CollectionFrom[Clippy](db, "clippy")
	SheetsCollection = bingo.CollectionFrom[Function](db, "sheets")
	ExcelCollection = bingo.CollectionFrom[Function](db, "excel")
}

func getDriver() *bingo.Driver {
	if db != nil {
		return db
	}
	config := bingo.DriverConfiguration{
		DeleteNoVerify: false,
		Filename:       "database/clippy.bingo.db",
	}
	db, _ = bingo.NewDriver(config)
	return db
}

func UrlToDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to request the webpage")
		return nil, err
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Failed to parse the webpage")
		return nil, err
	}
	return doc, nil
}
