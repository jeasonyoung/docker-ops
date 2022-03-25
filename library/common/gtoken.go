package common

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

//GTokenConfig gToken配置接口
type GTokenConfig interface {
	// GetCacheMode 缓存模式(1:Cache,2:Redis,默认:1)
	GetCacheMode() int8
	// GetCacheKey 缓存key
	GetCacheKey() string
	// GetTimeout 超时时间 默认10天（毫秒）
	GetTimeout() int
	// GetEncryptKey Token加密key
	GetEncryptKey() []byte
	// IsMultiLogin 是否支持多端登录，默认false
	IsMultiLogin() bool
	// GetLoginPath 登录路径
	GetLoginPath() string
	// LoginBeforeFunc 登录验证方法 return userKey 用户标识 如果userKey为空，结束执行
	LoginBeforeFunc(r *ghttp.Request) (string, interface{})
	// LoginAfterFunc 登录返回方法
	LoginAfterFunc(r *ghttp.Request, respData gtoken.Resp)
	// AuthAfterFunc 认证返回方法
	AuthAfterFunc(r *ghttp.Request, respData gtoken.Resp)
	// GetLogoutPath 登出地址
	GetLogoutPath() string
}

// BuildGTokenInstance 构建GToken对象实例
func BuildGTokenInstance(cfg GTokenConfig) *gtoken.GfToken {
	return &gtoken.GfToken{
		CacheMode:        cfg.GetCacheMode(),
		CacheKey:         cfg.GetCacheKey(),
		Timeout:          cfg.GetTimeout(),
		EncryptKey:       cfg.GetEncryptKey(),
		MultiLogin:       cfg.IsMultiLogin(),
		LoginPath:        cfg.GetLoginPath(),
		LoginBeforeFunc:  cfg.LoginBeforeFunc,
		LoginAfterFunc:   cfg.LoginAfterFunc,
		AuthAfterFunc:    cfg.AuthAfterFunc,
		LogoutPath:       cfg.GetLogoutPath(),
		AuthExcludePaths: g.SliceStr{"/captcha"},
	}
}
