package database

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nokusukun/bingo"
	"go-clippy/database/clippy"
	"go-clippy/database/functions"
	"log"
	"net/http"
)

var db *bingo.Driver

func init() {
	db = getDriver()
	clippy.Awards = bingo.CollectionFrom[clippy.Award](db, "clippy")
	clippy.Users = bingo.CollectionFrom[clippy.User](db, "users")
	clippy.Moderators = bingo.CollectionFrom[clippy.Moderator](db, "moderators")
	functions.SheetsCollection = bingo.CollectionFrom[functions.Function](db, "sheets")
	functions.ExcelCollection = bingo.CollectionFrom[functions.Function](db, "excel")
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

func Close() {
	err := getDriver().Close()
	if err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}
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
