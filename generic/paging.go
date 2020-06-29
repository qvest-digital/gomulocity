package generic

import (
	"net/http"
	"strconv"
)

/*
PagingStatistics represent cumulocity's 'application/vnd.com.nsn.cumulocity.pagingStatistics+json'.
See: https://cumulocity.com/guides/reference/rest-implementation/#pagingstatistics-application-vnd-com-nsn-cumulocity-pagingstatistics-json
*/
type PagingStatistics struct {
	TotalRecords int `json:"totalRecords,omitempty"`
	TotalPages   int `json:"totalPages,omitempty"`
	PageSize     int `json:"pageSize"`
	CurrentPage  int `json:"currentPage"`
}

// Page add query param 'currentPage' to request
func Page(Page int) func(*http.Request) {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set("currentPage", strconv.Itoa(Page))
		r.URL.RawQuery = q.Encode()
	}
}

// add query param 'pageSize' to request
func PageSize(pageSize int) func(*http.Request) {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set("pageSize", strconv.Itoa(pageSize))
		r.URL.RawQuery = q.Encode()
	}
}
