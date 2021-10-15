package parser

import "io"

// QueryAttrMap stores a mapping of parser query and an attribute which data
// should be retrieved. Example: QueryAttrMap{"div": "class"}. Here div elements
// should be found and the class value of each will be collected.
type QueryAttrMap map[string]string

// QueryResult stores all results parsed using the QueryAttrMap value. For preventing
// data duplicates the values are stored in a map structure.
type QueryResult map[string]struct{}

// Parser describes a generic set of parser functions.
type Parser interface {
	// FindContent searches for only one element and returns its text content.
	FindContent(r io.Reader, query string) (string, error)

	// FindAttrMap searches by multiple queries and retrieves corresponding attributes
	// of found elements. All queries and related attributes stores within QueryAttrMap.
	// All results will be stored in QueryResult.
	FindAttrMap(io.Reader, QueryAttrMap, QueryResult) error
}
