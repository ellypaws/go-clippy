package bingo

import "strings"

type FunctionScraper interface {
	Scrape() []Function
	UpdateUrls([]Function)
	UpdateSingleUrl(*Function)
}

func UrlToScrape(url string) FunctionScraper {
	// check if url contains microsoft or excel
	if strings.Contains(url, "microsoft") {
		excelUrl := ExcelUrl(url)
		return &excelUrl
	}
	if strings.Contains(url, "google") {
		return &SheetsScraper{}
	}
	return nil
}
