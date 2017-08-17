package urlrewrite

import (
	"net/http"
	"strings"
	"fmt"
	"log"
)

func UrlRewrite(r *http.Request, w http.ResponseWriter) bool {
	scheme := r.Header.Get("X-Forwarded-Proto")
	host := strings.TrimPrefix(r.Host, "www.")
	if scheme == "http" || len(host) < len(r.Host) {
		scheme = "https"
		url := fmt.Sprintf("%v://%v%v", scheme, host, r.RequestURI)

		log.Printf("\nRedirect to: %q - HTTP Code: %q\n"+url, http.StatusPermanentRedirect)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		return true
	}
	return false
}
