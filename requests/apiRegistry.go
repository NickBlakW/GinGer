package requests

var APIRequestRegistry []ApiRequest

type ApiRequest struct {
	Path string
	DTO any
}

func RegisterPostRequest(req ApiRequest) {
	APIRequestRegistry = append(APIRequestRegistry, req)
}