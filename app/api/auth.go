package api

import (
	"docker-ops-server/app/model"
	"docker-ops-server/app/service"
	"docker-ops-server/library/common"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

var Auth = new(authApi)

type authApi struct{}

//GetCacheMode 缓存模式(1:Cache, 2:Redis)
func (a *authApi) GetCacheMode() int8 {
	return g.Cfg().GetInt8("gToken.CacheMode", 1)
}

// GetCacheKey 缓存key
func (a *authApi) GetCacheKey() string {
	return g.Cfg().GetString("gToken.CacheKey")
}

// GetTimeout 超时时间 默认10天（毫秒）
func (a *authApi) GetTimeout() int {
	return g.Cfg().GetInt("gToken.Timeout", 86400*10*1000)
}

// GetEncryptKey Token加密key
func (a *authApi) GetEncryptKey() []byte {
	return g.Cfg().GetBytes("gToken.EncryptKey")
}

// IsMultiLogin 是否支持多端登录，默认false
func (a *authApi) IsMultiLogin() bool {
	return g.Cfg().GetBool("gToken.MultiLogin", false)
}

// GetLoginPath 登录路径
func (a *authApi) GetLoginPath() string {
	return "/auth/login"
}

//GetVerifyImage
//@id get-verify-image
//@tags 通用功能
//@summary 获取验证码图片
//@description 图片为base64格式
//@produce json
//@success 200 {object} response.RespResult{data=api.captchaOutput}
//@router /captcha [get]
func (a *authApi) GetVerifyImage(r *ghttp.Request) {
	captchaOnOff := g.Cfg().GetBool("server.captchaOnOff", true)
	idKey, base64 := service.Auth.GetVerifyImg()
	common.BuildRespSuccess(r, &captchaOutput{
		idKey, base64, captchaOnOff,
	})
}

type captchaOutput struct {
	Key          string `json:"key"`          //验证码键
	Content      string `json:"content"`      //验证码图片Base64编码
	CaptchaOnOff bool   `json:"captchaOnOff"` //是否关闭验证
}

// LoginBeforeFunc
//@id auth-login
//@tags 1.00.认证模块
//@summary 1.00.01.用户登录
//@accept json
//@param req body model.LoginBodyReq{} true "请求报文"
//@produce json
//@success 200 {object} response.RespResult{data=model.LoginBodyRes}
//@router /auth/login [post]
func (a *authApi) LoginBeforeFunc(r *ghttp.Request) (string, interface{}) {
	//获取登录参数
	var req *model.LoginBodyReq
	if err := r.Parse(&req); err != nil {
		common.BuildRespFailWithError(r, err.(gvalid.Error).Current())
	}
	//判断验证码是否正确
	captchaOnOff := g.Cfg().GetBool("server.captchaOnOff", true)
	if captchaOnOff && !service.Auth.VerifyString(req.VerifyKey, req.VerifyCode) {
		common.BuildRespFail(r, "验证码错误")
	}
	//验证用户
	account := req.Account
	user, err := service.Auth.Login(r.GetCtx(), account, req.Password)
	if err != nil {
		g.Log().Error(err)
		common.BuildRespFailWithError(r, err)
		return "", nil
	}
	//客户端IP
	ip := common.GetClientIp(r)
	//
	var key string
	if a.IsMultiLogin() {
		key = gconv.String(user.Id) + "-" + gmd5.MustEncryptString(account) + gmd5.MustEncryptString(req.Password+ip)
	} else {
		key = gconv.String(user.Id) + "-" + gmd5.MustEncryptString(account) + gmd5.MustEncryptString(req.Password)
	}
	return key, user
}

//LoginAfterFunc 登录成功后执行函数
func (a *authApi) LoginAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	if !respData.Success() {
		_ = r.Response.WriteJson(respData)
		return
	}
	token := respData.GetString("token")
	common.BuildRespSuccess(r, &model.LoginBodyRes{
		Token: token,
	})
}

// AuthAfterFunc 认证后执行
func (a *authApi) AuthAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	if r.Method == "OPTIONS" || respData.Success() {
		if respData.Success() {
			if val := respData.Data.(g.Map); val != nil {
				if user := val["data"].(*common.ContextUser); user != nil {
					//保存到上下文
					common.ContextService.Middleware(r, user)
					return
				}
			}
		}
	}
	common.BuildRespFailWithCode(r, 401, "未登录")
}

// GetLogoutPath
//@id auth-logout
//@tags 1.00.认证模块
//@summary 1.00.02.注销地址
//@accept json
//@produce json
//@success 200 {object} response.RespResult{}
//@router /auth/logout [post]
func (a *authApi) GetLogoutPath() string {
	return "/auth/logout"
}

//GetInfo
//@id auth-get-info
//@tags 1.00.认证模块
//@summary 1.00.03.用户信息
//@accept json
//@produce json
//@success 200 {object} response.RespResult{data=model.AuthUserInfoOutput}
//@router /getInfo [get]
func (a *authApi) GetInfo(r *ghttp.Request) {
	//获取当前用户
	ctxUser := common.ContextService.GetUser(r.GetCtx())
	if ctxUser == nil {
		common.BuildRespFail(r, "令牌解析失败")
		return
	}
	output := &model.UserInfoOutput{
		UserInfo: model.UserInfo{
			Id:      ctxUser.Id,
			Account: ctxUser.Account,
			Nick:    ctxUser.NickName,
		},
		Roles:       []string{"admin"},
		Permissions: []string{"*:*:*"},
	}
	common.BuildRespSuccess(r, output)
}

//GetRouters
//@id auth-get-routers
//@tags 1.00.认证模块
//@summary 1.00.04.路由数据
//@accept json
//@produce json
//@success 200 {object} response.RespResult{data=model.AuthUserInfoOutput}
//@router /routers [get]
func (a *authApi) GetRouters(r *ghttp.Request) {
	//获取当前用户
	ctxUser := common.ContextService.GetUser(r.GetCtx())
	if ctxUser == nil {
		common.BuildRespFail(r, "令牌解析失败")
		return
	}
	data, err := service.Auth.GetRouters(r.GetCtx(), ctxUser.Id)
	common.BuildRespWithError(r, data, err)
}
