package urlrewrite

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
	scheme := strings.ToLower(r.Header.Get(opt.SchemeHeader))
	host := r.Host

	targetScheme := scheme
	bareHost := strings.TrimPrefix(r.Host, "www.")

	if !opt.AcceptHTTP && scheme == "http" {
		scheme = "https"
	}

	if !opt.AcceptWWW && strings.HasPrefix(host, "www") {
		host = bareHost
	}

	if scheme != targetScheme || host != bareHost {
		url := fmt.Sprintf("%v://%v%v", scheme, host, r.RequestURI)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
