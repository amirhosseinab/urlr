package urlr

import (
	"fmt"
	"net/http"
	"strings"
)

type (
	URLRewrite struct {
		Options *Options
	}

	Options struct {
		AcceptHTTP   bool
		AcceptWWW    bool
		SchemeHeader string
	}
)

func New(opt *Options) *URLRewrite {
	return &URLRewrite{opt}
}

func Default() *URLRewrite {
	return New(&Options{
		AcceptHTTP:   false,
		AcceptWWW:    false,
		SchemeHeader: "X-Forwarded-Proto",
	})
}

func (u *URLRewrite) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rewrite(w, r, u.Options)
		handler.ServeHTTP(w, r)
	})
}

func rewrite(w http.ResponseWriter, r *http.Request, opt *Options) {
	scheme := r.Header.Get(opt.SchemeHeader)
	host := r.Host
	rewrite := false

	if !opt.AcceptHTTP && scheme == "http" {
		scheme = "https"
		rewrite = true
	}
	if !opt.AcceptWWW && strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(r.Host, "www.")
		rewrite = true
	}

	if rewrite {
		url := fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
