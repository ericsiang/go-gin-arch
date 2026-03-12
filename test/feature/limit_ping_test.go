package feature

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLimitPing(t *testing.T) {
	get, err := http.Get("http://127.0.0.1:5000/api/v1/limit_ping")
	if err != nil {
		t.Error(err)
	}
	if get.StatusCode != http.StatusOK {
		t.Error("http StatusCode not 200 ")
	}
}

// 通用http请求
func PerformRequest(path, method string, param []byte, router http.Handler) *http.Response {
	// 初始化响应
	w := httptest.NewRecorder()
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest(method, path, bytes.NewReader(param))

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	return result
}
