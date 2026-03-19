// Package ginresp 定義了 Gin 應用程式的通用回應結構
package ginresp

// FailResponse 失敗回應
type FailResponse struct {
	Msg string `json:"msg"`
}
