package generic

import (
	"fmt"
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

// gets query param 'pageSize' for a request as string and as an added to the provided values parameter
func PageSizeParameter(pageSize int, params *url.Values) (string, error) {
	if pageSize < 1 || pageSize > 2000 {
		return "", fmt.Errorf("The page size must be between 1 and 2000. Was %d", pageSize)
	}

	if params == nil {
		params = &url.Values{}
	}
	params.Add("pageSize", strconv.Itoa(pageSize))

	return params.Encode(), nil
}

