package bingo

type FunctionScraper interface {
	Scrape() []Function
	UpdateUrls([]Function)
	UpdateSingleUrl(*Function)
}
