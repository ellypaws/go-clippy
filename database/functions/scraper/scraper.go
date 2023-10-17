package main

import (
	"go-clippy/database/functions"
	"strings"
)

type FunctionScraper interface {
	Scrape() []functions.Function
	UpdateUrls([]functions.Function)
	UpdateSingleUrl(*functions.Function)
}

func UrlToScrape(url string) FunctionScraper {
	// check if url contains microsoft or excel
	if strings.Contains(url, "microsoft") {
		excelUrl := ExcelUrl(url)
		return &excelUrl
	}
	if strings.Contains(url, "google") {
		sheetsUrl := SheetsUrl(url)
		return &sheetsUrl
	}
	return nil
}
