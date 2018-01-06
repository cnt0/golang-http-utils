package utils

import (
	"net/http"
	"sync"
)

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
