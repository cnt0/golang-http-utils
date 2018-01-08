package utils

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// defined in src/activate/activate.c from systemd source tree
const sdListenFDsStart = 3

// HTTPResponseWithErr is just what it is
type HTTPResponseWithErr struct {
	http.Response
	error
}

// ParallelRequests runs http requests in parallel
func ParallelRequests(reqs []*http.Request, client *http.Client) []*HTTPResponseWithErr {
	resps := make([]*HTTPResponseWithErr, len(reqs))
	var wg sync.WaitGroup
	wg.Add(len(reqs))
	for idx, req := range reqs {
		go func(r *http.Request, c *http.Client, idx int) {
			resp, err := c.Do(r)
			resps[idx] = &HTTPResponseWithErr{*resp, err}
			wg.Done()
		}(req, client, idx)
	}
	wg.Wait()
	return resps
}

// ListenAndServeSA listens on unix socket provided by systemd socket activation.
// If there is no socket activation, it will act just like http.ListenAndServe
func ListenAndServeSA(addr string, handler http.Handler) error {
	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
		l, err := net.FileListener(os.NewFile(sdListenFDsStart, "socket"))
		if err != nil {
			return err
		}
		if err := http.Serve(l, handler); err != nil {
			return err
		}
	}
	return http.ListenAndServe(addr, handler)
}
