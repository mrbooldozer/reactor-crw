//go:build unit
// +build unit

package parser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"reactor-crw/parser"
)

func TestHtml_FindOne(t *testing.T) {
	p := parser.Html{}

	r := strings.NewReader(`
		<body>
			<h1>heading</h1>
			<a class="test-class" href="href-test" target="_blank">First Link</a>
			<a class="class" href="href-test" target="_blank">Second Link</a>
			<a href="href-test" target="_blank">Third Link</a>
		</body>
	`)

	res, err := p.FindContent(r, "body .test-class")
	require.NoErrorf(t, err, "Wasn't expected an error on html parse")

	require.Equal(t, "First Link", res)
}

func TestHtml_FindAttrMap(t *testing.T) {
	p := parser.Html{}

	r := strings.NewReader(`
		<body>
			<img src="image-src">
			<a class="test-class" href="href-test" target="_blank">First Link</a>
			<div class="div-class"></div>
		</body>
	`)

	queryAttrMap := parser.QueryAttrMap{
		".test-class": "href",
		"div":         "class",
		"img":         "src",
		"body":        "src",
	}
	res := make(parser.QueryResult)

	err := p.FindAttrMap(r, queryAttrMap, res)
	require.NoErrorf(t, err, "Wasn't expected an error on html parse")
	require.Equal(
		t,
		parser.QueryResult{
			"image-src": struct{}{},
			"href-test": struct{}{},
			"div-class": struct{}{},
		},
		res,
	)
}
