package api

type ApiResponse struct {
	Code    int	`json:"code"`
	Message string	`json:"message"`
	Data    interface{}	`json:"data"`
}


func ApiResponseFunc(code int, message string, data interface{}) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
