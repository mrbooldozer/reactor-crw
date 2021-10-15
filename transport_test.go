// +build unit

package reactor_crw_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reactor-crw"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpTransport_FetchData(t *testing.T) {
	t.Log("Given the need to fetch the data.")
	{
		t.Log("When request params are valid.")
		{
			expected := "response"

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "test-val", r.Header.Get("test-key"))
				_, _ = fmt.Fprintf(w, expected)
			}))
			defer srv.Close()

			httpTransport := reactor_crw.NewHttpTransport(
				http.DefaultClient,
				reactor_crw.Headers{"test-key": "test-val"},
			)

			data, err := httpTransport.FetchData(srv.URL)
			require.NoErrorf(t, err, "Wasn't expected an error during http call")

			response, _ := ioutil.ReadAll(data)
			require.Equal(t, expected, string(response))
		}

		t.Log("When request cannot be prepared.")
		{
			httpTransport := reactor_crw.NewHttpTransport(http.DefaultClient, nil)

			_, err := httpTransport.FetchData("\u007F")
			require.Error(t, err, "Expected an error during request")
		}

		t.Log("When endpoint is invalid.")
		{
			httpTransport := reactor_crw.NewHttpTransport(http.DefaultClient, nil)

			_, err := httpTransport.FetchData(" ")
			require.Error(t, err, "Expected an error during request")
		}
	}
}
