package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go"
)

type HMap = map[string]string

type RequestMaker struct {
	client           *http.Client
	headers          HMap
	status_forcelist []uint
	max_retries      uint
	backoff          uint
}

type RequestParams struct {
	Method          string
	Url             string
	Headers         HMap
	BodyReader      io.Reader
	StatusForcelist []uint
}

type InitParams struct {
	Timeout         time.Duration
	Headers         HMap
	StatusForcelist []uint
	MaxRetries      uint
	BackoffFactor   uint
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}


func UrlCombine(url, urlp string) string {
	if urlp == `` || urlp == `/` {
		return url
	}
	up := strings.SplitN(url, "?", 2)
	left := strings.TrimRight(up[0], `/`)
	is_multi, has_prefix := len(up) > 1, urlp[0] == '/'
		
	if is_multi && has_prefix {
		return left + urlp + `?` + up[1]
	}
	if ! is_multi && has_prefix {
		 return left + urlp
	}
	handled_urlp := strings.TrimFunc(urlp, func(r rune) bool { return r == '?' || r == '&' })
	if is_multi && ! has_prefix {
		return left + `?` + handled_urlp + `&` + up[1] 
	}
	// if ! is_multi && ! has_prefix
	return left + `?` + handled_urlp
}

func (r *RequestMaker) Init(params *InitParams) *RequestMaker {
	r.client = &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: params.Timeout,
		},
	}
	r.headers = params.Headers
	r.status_forcelist = params.StatusForcelist
	r.max_retries = params.MaxRetries
	r.backoff = params.BackoffFactor
	return r
}

func (r *RequestMaker) InitDefault() *RequestMaker {
	return r.Init(&InitParams{
		Timeout: 30 * time.Second,
		Headers: HMap{
			`user-agent`: `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0`,
		},
		StatusForcelist: []uint{500, 502, 503, 504},
		MaxRetries: 3,
		BackoffFactor: 2,
	})
}

func UrlCombineAll(urls ...string) string {
	url := urls[0]
	for _, u := range urls[1:] {
		url = UrlCombine(url, u)
	}
	return url
}

func set_headers(r *http.Request, map_array ...*HMap) {
	for _, m := range map_array {
		for k, v := range *m {
			r.Header.Set(k, v)
		}
	}
}

func (r *RequestMaker) Request(params *RequestParams) (resp *http.Response, err error) {
	req, err := http.NewRequest(params.Method, params.Url, params.BodyReader)
	if err != nil {
		return
	}
	set_headers(req, &r.headers, &params.Headers)

	retry.Do(
		func() error {
			resp, err = r.client.Do(req)

			if params.StatusForcelist != nil && Contains(params.StatusForcelist, uint(resp.StatusCode)) {
				resp.Body.Close()
				err = errors.New(fmt.Sprintf("Status code from force list %d", resp.StatusCode))
				return err
			}

			return err
		},
		retry.Attempts(uint(r.max_retries)),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying request after error: %v", err)
		}),
		retry.DelayType(retry.BackOffDelay),
	)
	if err != nil {
		return
	}
	//defer resp.Body.Close()
	return
}
