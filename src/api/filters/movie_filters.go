package filters

import (
	"math"
	"strings"
)

type Filter struct {
	Page         int      `json:"page"`
	PageSize     int      `json:"pageSize"`
	Sort         string   `json:"sort"`
	SortSafeList []string `json:"sortSafeList"`
}

func (f Filter) SortCol() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter" + f.Sort)
}

func (f Filter) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

type MetaData struct {
	CurrentPage int `json:"currentPage,omitempty"`
	PageSize    int `json:"pageSize,omitempty"`
	FirstPage   int `json:"firstPage,omitempty"`
	LastPage    int `json:"lastPage,omitempty"`
	TotalRecord int `json:"totalRecord,omitempty"`
}

func CalculateMetaData(totalRecord, page, pageSize int) MetaData {
	if totalRecord == 0 {
		return MetaData{}
	}
	return MetaData{
		CurrentPage: page,
		PageSize:    pageSize,
		FirstPage:   1,
		LastPage:    int(math.Ceil(float64(totalRecord) / float64(pageSize))),
		TotalRecord: totalRecord,
	}
}
