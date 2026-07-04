package model

// APIResponse 统一 API 响应格式，见 DEVELOPMENT.md §5.2
type APIResponse struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(data interface{}) APIResponse {
	return APIResponse{OK: true, Data: data}
}

func Fail(message, code string) APIResponse {
	return APIResponse{OK: false, Message: message, Code: code}
}
