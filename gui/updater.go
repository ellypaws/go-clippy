package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-clippy/database/functions"
	"go-clippy/database/functions/scraper"
	"log"
	"os"
	"strings"
	"time"
)

var sliceFuncs []functions.Function

// var excelUrl = UrlToScrape("https://support.microsoft.com/en-us/office/excel-functions-alphabetical-b3944572-255d-4efb-bb96-c6d90033e188")
var sheetsUrl = scraper.UrlToScrape("https://support.google.com/docs/table/25273?hl=en")

type MainModel struct {
	senders []tea.Model
	jobs    chan functions.Function
	results chan functions.Function
}

func newMain() MainModel {
	var senders []tea.Model
	if sliceFuncs == nil {
		sliceFuncs = sheetsUrl.Scrape()
	}

	var numJobs = len(sliceFuncs)
	const numWorkers = 3
	jobs := make(chan functions.Function, numJobs)
	results := make(chan functions.Function, numJobs)

	// start workers
	for id := 0; id <= numWorkers-1; id++ {
		sender := newSender()
		senders = append(senders, sender)
		go worker(id, jobs, results)
	}

	// send all jobs
	for _, function := range sliceFuncs {
		jobs <- function
	}
	close(jobs)

	return MainModel{
		senders: senders,
		jobs:    jobs,
		results: results,
	}
}

func (m MainModel) Init() tea.Cmd {
	return func() tea.Msg { return spinner.TickMsg{} }
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sendTo:
		m.senders[msg.id], _ = m.senders[msg.id].Update(msg.msg)
		return m, nil
	case spinner.TickMsg:
		for id, sender := range m.senders {
			sender, _ = sender.Update(msg)
			m.senders[id] = sender
		}
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return spinner.TickMsg{}
		})
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlF {
			finish(m.results)
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m MainModel) View() string {
	views := make([]string, 0)
	for _, sender := range m.senders {
		views = append(views, sender.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, views...)
}

func finish(results chan functions.Function) {
	// wait for all results
	for a := 1; a <= len(sliceFuncs); a++ {
		<-results
	}

	sliceFuncs = make([]functions.Function, 0)

	for function := range results {
		sliceFuncs = append(sliceFuncs, function)
	}

	// save to json
	indent, _ := json.MarshalIndent(sliceFuncs, "", "    ")
	toPrint := strings.ReplaceAll(string(indent), "\\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#e07a00")).Render("↵\n"))
	log.Printf("%+v\n", toPrint)
	os.WriteFile("sheets.json", indent, 0644)
}

var p *tea.Program

func main() {
	p = tea.NewProgram(newMain())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func worker(id int, jobs <-chan functions.Function, results chan<- functions.Function) {
	for j := range jobs {
		//fmt.Println("worker", id, "updating:", j.Name)
		now := time.Now()
		url := scraper.SheetsUrl(j.URL)
		url.UpdateSingleUrl(&j)
		results <- j
		//fmt.Println("worker", id, "updated:", j.Name, "in", time.Since(now))
		p.Send(sendTo{
			id: id,
			msg: resultMsg{
				food:     j.Name,
				duration: time.Since(now),
			},
		})
	}
}

type sendTo struct {
	id  int
	msg resultMsg
}
