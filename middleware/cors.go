package middleware

import "github.com/gogf/gf/net/ghttp"

//Cors 跨域处理
func Cors(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
