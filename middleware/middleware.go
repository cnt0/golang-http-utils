package middleware

import (
	"net/http"
	"sync"
	"time"
)

// IfModifiedSince uses "If-Modified-Since" and "Last-Modified" http headers
// to optimize http requests.
type IfModifiedSince struct {
	m sync.Mutex
	t time.Time
}

// Handler adds "If-Modified-Since" header to the given request
func (i *IfModifiedSince) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i.m.Lock()
		defer i.m.Unlock()
		r.Header.Set("If-Modified-Since", i.t.Format(http.TimeFormat))
		next.ServeHTTP(w, r)
		if t, err := time.Parse(http.TimeFormat, w.Header().Get("Last-Modified")); err == nil {
			i.t = t
		}
	})
}

// NewIfModifiedSince is just an IfModifiedSince constructor
func NewIfModifiedSince() IfModifiedSince {
	return IfModifiedSince{t: time.Now()}
}
