// +build unit

package reactor_crw

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"reactor-crw/parser"
)

type transportMock struct {
	mock.Mock
}

func (m *transportMock) FetchData(url string) (io.ReadCloser, error) {
	args := m.Called(url)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

type parserMock struct {
	mock.Mock
}

func (m *parserMock) FindContent(r io.Reader, q string) (string, error) {
	args := m.Called(r, q)
	return args.String(0), args.Error(1)
}

func (m *parserMock) FindAttrMap(r io.Reader, qa parser.QueryAttrMap, qr parser.QueryResult) error {
	args := m.Called(r, qa, qr)
	return args.Error(0)
}

func TestHtmlCrawler_Fetch(t *testing.T) {
	path := "https://test.com/test/path"

	trp := &transportMock{}
	prs := &parserMock{}

	t.Log("Given the need to crawl html page.")
	{
		t.Log("When multiple pages requested")
		{
			c := &HtmlCrawler{trp, prs, true}

			rc := ioutil.NopCloser(strings.NewReader(""))
			trp.On("FetchData", path).Return(rc, nil).Once()

			prs.On("FindContent", rc, ".pagination_expanded .current").Return("2", nil)
			trp.On("FetchData", path+"/1").Return(rc, nil).Once()
			trp.On("FetchData", path+"/2").Return(rc, nil).Once()

			prs.On("FindAttrMap", rc, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					qr := args.Get(2).(parser.QueryResult)
					qr["link_1"] = struct{}{}
				}).
				Return(nil).
				Once()

			prs.On("FindAttrMap", rc, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					qr := args.Get(2).(parser.QueryResult)
					qr["link_2"] = struct{}{}
				}).
				Return(nil).
				Once()

			res, err := c.Fetch(path, []string{"image"})
			require.NoErrorf(t, err, "Wasn't expected an error during crawl")
			require.Equal(t, parser.QueryResult{"link_1": struct{}{}, "link_2": struct{}{}}, res)
		}

		t.Log("When parser returned an error")
		{
			expectedErr := errors.New("error")
			c := &HtmlCrawler{trp, prs, false}

			rc := ioutil.NopCloser(strings.NewReader(""))
			trp.On("FetchData", path).Return(rc, nil).Once()

			prs.On("FindAttrMap", rc, mock.Anything, mock.Anything).
				Return(expectedErr).
				Once()

			_, err := c.Fetch(path, []string{"image"})
			assert.ErrorIs(t, err, expectedErr)
		}

		t.Log("When transport returned an error")
		{
			expectedErr := errors.New("error")
			c := &HtmlCrawler{trp, prs, false}

			rc := ioutil.NopCloser(strings.NewReader(""))
			trp.On("FetchData", path).Return(rc, expectedErr).Once()

			_, err := c.Fetch(path, []string{"image"})
			assert.ErrorIs(t, err, expectedErr)
		}

		t.Log("When single page requested")
		{
			c := &HtmlCrawler{trp, prs, false}

			rc := ioutil.NopCloser(strings.NewReader(""))
			trp.On("FetchData", path).Return(rc, nil).Once()

			prs.On("FindAttrMap", rc, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					qr := args.Get(2).(parser.QueryResult)
					qr["link_1"] = struct{}{}
				}).
				Return(nil).
				Once()

			res, err := c.Fetch(path, nil)
			require.NoErrorf(t, err, "Wasn't expected an error during crawl")
			require.Equal(t, parser.QueryResult{"link_1": struct{}{}}, res)
		}
	}
}
