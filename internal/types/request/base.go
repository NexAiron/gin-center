package request

type PaginationRequest struct {
	Page     int `json:"page" binding:"required,min=1"`
	PageSize int `json:"page_size" binding:"required,min=1,max=100"`
}
type SortRequest struct {
	SortBy    string `json:"sort_by" binding:"omitempty"`
	SortOrder string `json:"sort_order" binding:"omitempty,oneof=asc desc"`
}
type SearchRequest struct {
	Keyword string `json:"keyword" binding:"omitempty"`
}
type BaseRequest struct {
	RequestID   string `json:"request_id" binding:"required,uuid4"`
	ClientIP    string `json:"client_ip" binding:"required,ipv4"`
	UserAgent   string `json:"user_agent" binding:"omitempty,max=256"`
	RequestTime int64  `json:"request_time" binding:"required,number"`
}
