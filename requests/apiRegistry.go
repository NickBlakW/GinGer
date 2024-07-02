package requests

var APIRequestRegistry []ApiRequest

type ApiRequest struct {
	Path string
	DTO any
}

func RegisterRequest(req ApiRequest) {
	APIRequestRegistry = append(APIRequestRegistry, req)
}

func RegisterRequests(reqs []ApiRequest) {
	APIRequestRegistry = append(APIRequestRegistry, reqs...)
}