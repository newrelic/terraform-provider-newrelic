package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/tomnomnom/linkheader"
)

type Pager interface {
	Parse(res *resty.Response) Paging
}

type Paging struct {
	Next string
}

type LinkHeaderPager struct{}

func (l *LinkHeaderPager) Parse(res *resty.Response) Paging {
	paging := Paging{}
	header := res.Header().Get("Link")
	if header != "" {
		links := linkheader.Parse(header)

		for _, link := range links.FilterByRel("next") {
			paging.Next = link.URL
			break
		}
	}

	return paging
}
