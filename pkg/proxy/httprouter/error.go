package httprouter

import (
	"github.com/valyala/fasthttp"
)

type errorProcessing func(res *fasthttp.Response, err error, statusCode int)

// ErrorProcessing ...
func ErrorProcessing(res *fasthttp.Response, err error, statusCode int) {
	res.SetBody([]byte(err.Error()))
	res.SetStatusCode(statusCode)
}
