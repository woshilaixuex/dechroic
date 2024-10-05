package http_model

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: http请求与返回模型
 * @Date: 2024-10-03 23:54
 */

type Response[T any] struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
