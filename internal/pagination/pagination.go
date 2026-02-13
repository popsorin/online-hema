package pagination

import (
	"net/http"
	"strconv"
)

type Params struct {
	Page     int
	PageSize int
	Offset   int
}

type Response struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

func ParseParams(r *http.Request) Params {
	page := parseIntParam(r, "page", 1)
	pageSize := parseIntParam(r, "page_size", 20)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	return Params{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

func NewResponse(data interface{}, params Params, totalCount int) Response {
	totalPages := (totalCount + params.PageSize - 1) / params.PageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return Response{
		Data:       data,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}

func parseIntParam(r *http.Request, key string, defaultValue int) int {
	str := r.URL.Query().Get(key)
	if str == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}
