package router

import (
	"docker-ops-server/app/api"
	"docker-ops-server/library/common"
	"docker-ops-server/middleware"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	//token
	token := common.BuildGTokenInstance(api.Auth)
	//Web服务
	s := g.Server()
	//中间件处理器
	s.Use(middleware.Cors)
	//路由
	s.Group("/", func(group *ghttp.RouterGroup) {
		//token 拦截器
		_ = token.Middleware(group)
		// 图片验证码
		group.GET("/captcha", api.Auth.GetVerifyImage)
		//用户信息
		group.GET("/getInfo", api.Auth.GetInfo)
		//路由数据
		group.GET("/getRouters", api.Auth.GetRouters)
	})
}
