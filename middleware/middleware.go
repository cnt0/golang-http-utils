package middleware

import (
	"net/http"
	"time"
)

// IfModifiedSince uses "If-Modified-Since" and "Last-Modified" http headers
// to optimize http requests. Obviously not a thread-safe
type IfModifiedSince struct {
	t time.Time
}

// Handler adds "If-Modified-Since" header to the given request
func (i *IfModifiedSince) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

// WithTimeout middleware implements a timeout for http request handlers
type WithTimeout struct {
	d time.Duration
}

// NewWithTimeout is a constructor for WithTimeout struct
func NewWithTimeout(d time.Duration) WithTimeout {
	return WithTimeout{d: d}
}

// Handler implements a timeout for http request handlers
func (wt *WithTimeout) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := make(chan bool)

		go func(ok chan bool, wt *WithTimeout) {
			next.ServeHTTP(w, r)
			ok <- true
		}(ok, wt)

		select {
		case <-ok:
			return
		case <-time.After(wt.d):
			return
		}
	})
}
