package urlr_test

import (
	"io"
	"net/http"
	"testing"

	"net/http/httptest"

	"fmt"

	"github.com/amirhosseinab/urlrewrite"
)

func TestURLRewrite_Handler(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "dummy content")
	})
	table := []struct {
		target, headerName, headerValue string

		acceptHTTP, acceptWWW bool
		optionHeaderName      string

		statusCode int
		location   string
	}{
		{
			target:           "http://www.domain.tld",
			headerName:       "dummy",
			headerValue:      "dummy",
			acceptHTTP:       true,
			acceptWWW:        true,
			optionHeaderName: "dummy",
			statusCode:       http.StatusOK,
			location:         "",
		},
		{
			target:           "http://www.domain.tld",
			headerName:       "scheme",
			headerValue:      "http",
			acceptHTTP:       false,
			acceptWWW:        true,
			optionHeaderName: "scheme",
			statusCode:       http.StatusPermanentRedirect,
			location:         "https://www.domain.tld",
		},
		{
			target:           "http://www.domain.tld",
			headerName:       "scheme",
			headerValue:      "http",
			acceptHTTP:       true,
			acceptWWW:        false,
			optionHeaderName: "scheme",
			statusCode:       http.StatusPermanentRedirect,
			location:         "http://domain.tld",
		},
		{
			target:           "http://www.domain.tld",
			headerName:       "scheme",
			headerValue:      "http",
			acceptHTTP:       false,
			acceptWWW:        false,
			optionHeaderName: "scheme",
			statusCode:       http.StatusPermanentRedirect,
			location:         "https://domain.tld",
		},
		{
			target:           "http://www.domain.tld/some/request/uri",
			headerName:       "scheme",
			headerValue:      "http",
			acceptHTTP:       false,
			acceptWWW:        false,
			optionHeaderName: "scheme",
			statusCode:       http.StatusPermanentRedirect,
			location:         "https://domain.tld/some/request/uri",
		},
	}

	for _, tt := range table {
		name := fmt.Sprintf("target:[%s] http:%v www:%v/code:%d %s", tt.target, tt.acceptHTTP, tt.acceptWWW, tt.statusCode, tt.location)
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.target, nil)
			req.Header.Set(tt.headerName, tt.headerValue)
			w := httptest.NewRecorder()
			opt := &urlr.Options{AcceptHTTP: tt.acceptHTTP, AcceptWWW: tt.acceptWWW, SchemeHeader: tt.headerName}

			sut := urlr.New(opt).Handler(h)
			sut.ServeHTTP(w, req)

			resp := w.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("%s\n\tStatusCode should be %d, got: %d", name, tt.statusCode, resp.StatusCode)
			}
			if tt.statusCode == http.StatusPermanentRedirect {
				loc := resp.Header.Get("Location")
				if tt.location != loc {
					t.Errorf("%s\n\tLocation Header should be %s, got: %s", name, tt.location, loc)
				}
			}
		})
	}
}
