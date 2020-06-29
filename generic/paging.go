package generic

import (
	"fmt"
	"net/http"
	"net/url"
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

// gets query param 'pageSize' for a request
func PageSizeParameter(pageSize int) (string, error) {
	if pageSize < 1 || pageSize > 2000 {
		return "", fmt.Errorf("The page size must be between 1 and 2000. Was %d", pageSize)
	}

	params := url.Values{}
	params.Add("pageSize", strconv.Itoa(pageSize))

	return params.Encode(), nil
}
