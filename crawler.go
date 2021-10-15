package reactor_crw

import (
	"fmt"
	"io"
	"strconv"

	"reactor-crw/parser"
)

// Transport is an interface that wraps all network related operations performed
// by crawler.
type Transport interface {
	FetchData(url string) (io.ReadCloser, error)
}

// HtmlCrawler allows crawling HTML pages using provided transport.Transport and
// parser.Parser. The crawler doesn't do anything with the content itself it only
// gathers the collection of sources links.
type HtmlCrawler struct {
	// Transport performs all network request.
	Transport Transport

	// Parser contains parser.Parser implementation for HTML.
	Parser parser.Parser

	// MultiPage allows applying the crawler for multiple pages depending on its
	// value. If its value set to true the crawler will try to calculate the
	// number of available pages using the source page.
	MultiPage bool
}

// Fetch retrieves content sources from the page using the path value. Depending
// on HtmlCrawler.MultiPage it may fetch content links from multiple pages.
func (c *HtmlCrawler) Fetch(path string, search []string) (parser.QueryResult, error) {
	if c.MultiPage {
		return c.multiPage(path, search)
	}

	collectedData := make(parser.QueryResult)

	err := c.fetch(path, search, collectedData)
	if err != nil {
		return nil, err
	}

	return collectedData, nil
}

// fetchMultiPage will fetch content sources from multiple pages. It'll try to
// retrieve the number of pages and consequently crawl all pages.
func (c *HtmlCrawler) multiPage(path string, search []string) (parser.QueryResult, error) {
	maxPage, err := c.resolveMaxPage(path)
	if err != nil {
		return nil, err
	}

	collectedData := make(parser.QueryResult)

	for p := 1; p <= maxPage; p++ {
		err = c.fetch(fmt.Sprintf("%s/%d", path, p), search, collectedData)
		if err != nil {
			return nil, err
		}
	}

	return collectedData, nil
}

func (c *HtmlCrawler) fetch(path string, search []string, qr parser.QueryResult) error {
	body, err := c.Transport.FetchData(path)
	if err != nil {
		return err
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(body)

	err = c.Parser.FindAttrMap(body, buildQuery(search), qr)
	if err != nil {
		return fmt.Errorf("cannot apply crawler: %w", err)
	}

	return nil
}

func (c *HtmlCrawler) resolveMaxPage(path string) (int, error) {
	const htmlPagination = ".pagination_expanded .current"

	body, err := c.Transport.FetchData(path)
	if err != nil {
		return 0, err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(body)

	pa, err := c.Parser.FindContent(body, htmlPagination)
	if err != nil {
		return 0, err
	}

	intPa, err := strconv.Atoi(pa)
	if err != nil {
		return 0, err
	}

	return intPa, nil
}

// query represents a search query that will be applied against the parser.
// Each query has a query itself to find corresponding elements, its attr,
// content of which should be returned, contentType marks a type of content
// and used for fetching specific types only.
type query struct {
	contentType string
	query       string
	attr        string
}

var queries = []query{
	{"image", ".post_content .image > img", "src"},
	{"image", ".post_content .image > a", "href"},
	{"gif", ".post_content .video_gif_source", "href"},
	{"mp4", ".post_content .video_gif source[type='video/mp4']", "src"},
	{"webm", ".post_content .video_gif source[type='video/webm']", "src"},
}

// buildQuery builds a parser.QueryAttrMap according to provided search list.
// The resulting parser.QueryAttrMap will contain only those queries that meet
// required content types.
//
// Example: buildQuery([]string{"image", "mp4"}).
func buildQuery(search []string) parser.QueryAttrMap {
	qa := parser.QueryAttrMap{}

	searchMap := make(map[string]struct{}, len(search))
	for _, s := range search {
		searchMap[s] = struct{}{}
	}

	for _, sq := range queries {
		if _, ok := searchMap[sq.contentType]; ok {
			qa[sq.query] = sq.attr
		}
	}

	return qa
}
