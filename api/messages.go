package api

import (
	"net/http"
)

func PutMessage(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}

func ConsumeHead(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}

func ConsumeTail(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}

func ConsumeIndex(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}

func PeekIndex(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}
