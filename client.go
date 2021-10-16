package reactor_crw

import (
	"reactor-crw/handler"
	"strings"
	"sync"

	"reactor-crw/parser"
)

// Crawler is an interface for crawler used by the client. It fetches the content
// sources by the provided path and search parameters which are just a list of
// required content types. The crawler itself doesn't do anything with fetched
// data but only collects content sources.
type Crawler interface {
	Fetch(path string, search []string) (parser.QueryResult, error)
}

// Client defines a facade for a specific crawler implementation and a list of
// content handlers. The client runs the whole process of crawling the data and
// apply it against a list of provided handlers.
type Client struct {
	TotalSources chan int
	Progress     chan int
	Errors       chan error

	crawler    Crawler
	handler    handler.ContentHandler
	maxWorkers int
}

// NewClient creates a new crawler client and initialize client channels.
func NewClient(c Crawler, maxWorkers int, ch handler.ContentHandler) *Client {
	return &Client{
		crawler:      c,
		handler:      ch,
		maxWorkers:   maxWorkers,
		TotalSources: make(chan int, 1),
		Progress:     make(chan int),
		Errors:       make(chan error),
	}
}

// Run starts all the work by running Crawler and processes its search results
// with the corresponding handler.ContentHandler.
//
// All progress and errors signals will be sent to Progress and Errors channels
// respectively. TotalSources channel is used to notify the client about the
// amount of found content links.
//
// Content sources will be processed by handler.ContentHandler through a simple
// worker pool.
func (c *Client) Run(path string, search string) error {
	defer func() {
		close(c.Progress)
		close(c.Errors)
	}()

	collectedData, err := c.crawler.Fetch(path, strings.Split(search, ","))
	c.TotalSources <- len(collectedData)
	if err != nil || len(collectedData) == 0 {
		return err
	}

	contentHandlerTasks := make(chan string, len(collectedData))
	for task := range collectedData {
		contentHandlerTasks <- task
	}
	close(contentHandlerTasks)

	var wg sync.WaitGroup

	for i := 0; i < c.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range contentHandlerTasks {
				c.handler.Process(t, c.Progress, c.Errors)
			}
		}()
	}

	wg.Wait()

	return nil
}
