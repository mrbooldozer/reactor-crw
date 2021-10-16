package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
)

// Html represents a HTML parser implementation.
type Html struct{}

// FindContent searches for the first element that satisfies the query and returns
// its text node as content.
func (h *Html) FindContent(r io.Reader, query string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", fmt.Errorf("cannot parse document: %w", err)
	}

	return doc.Find(query).First().Text(), nil
}

// FindAttrMap parses HTML documents with multiple queries and retrieves the
// corresponding attributes of found elements. All queries and related attributes
// are stores within QueryAttrMap. All results will be stored in QueryResult.
//
// Example: p.FindAttrMap(body, QueryAttrMap{"div": "class"}, QueryResult{})
func (h *Html) FindAttrMap(r io.Reader, q QueryAttrMap, res QueryResult) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("cannot parse document: %w", err)
	}

	for query, attr := range q {
		doc.Find(query).Each(func(_ int, s *goquery.Selection) {
			val, ok := s.Attr(attr)
			if !ok || val == "javascript:" {
				return
			}

			val, _ = url.QueryUnescape(val)
			res[val] = struct{}{}
		})
	}

	return nil
}
