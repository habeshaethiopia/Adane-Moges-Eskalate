package handlers

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Object  interface{} `json:"object,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Object     interface{} `json:"object"`
	PageNumber int         `json:"pageNumber"`
	PageSize   int         `json:"pageSize"`
	TotalSize  int64       `json:"totalSize"`
	Errors     []string    `json:"errors,omitempty"`
}
