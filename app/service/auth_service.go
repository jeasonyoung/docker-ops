package service

import (
	"context"
	"docker-ops-server/app/model"
	"docker-ops-server/library/common"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/mojocn/base64Captcha"
)

var Auth = new(authService)

const (
	typeDir    = "M"
	typeMenu   = "C"
	typeButton = "F"
	noFrame    = 0
)

type authService struct {
	routerCache *gcache.Cache
}

// GetVerifyImg 获取字母数字混合验证码
func (a *authService) GetVerifyImg() (idKey, base64Img string) {
	driver := &base64Captcha.DriverString{
		Height:          80,
		Width:           240,
		NoiseCount:      50,
		ShowLineOptions: 20,
		Length:          4,
		Source:          "abcdefghjkmnpqrstuvwxyz23456789",
		Fonts:           []string{"chromohv.ttf"},
	}
	driver = driver.ConvertFonts()
	store := base64Captcha.DefaultMemStore
	c := base64Captcha.NewCaptcha(driver, store)
	idKey, base64Img, err := c.Generate()
	if err != nil {
		g.Log().Error(err)
	}
	return
}

func (a *authService) VerifyString(id, answer string) bool {
	driver := new(base64Captcha.DriverString)
	store := base64Captcha.DefaultMemStore
	c := base64Captcha.NewCaptcha(driver, store)
	answer = gstr.ToLower(answer)
	return c.Verify(id, answer, true)
}

// Login 用户登录
func (a *authService) Login(ctx context.Context, account, passwd string) (*common.ContextUser, error) {
	if account == "" || passwd == "" {
		return nil, gerror.New("账号或密码错误为空")
	}
	pwd := g.Cfg().GetString(fmt.Sprintf("users.%s", account), "")
	if pwd == "" {
		g.Log().Warningf("用户名不存在: %s", account)
		return nil, gerror.New("账号或密码错误为空.")
	}
	if pwd != passwd {
		g.Log().Warningf("密码错误: %s => %s", passwd, pwd)
		return nil, gerror.New("账号或密码错误为空..")
	}
	return &common.ContextUser{
		Id:       0,
		Account:  account,
		NickName: account,
	}, nil
}

// GetRouters 获取路由集合
func (a *authService) GetRouters(ctx context.Context, userId uint64) ([]*model.RouterOutput, error) {
	cacheKey := "cache-routers"
	//检查缓存
	if a.routerCache != nil {
		if cache, _ := a.routerCache.Get(cacheKey); cache != nil {
			if out, ok := cache.([]*model.RouterOutput); ok {
				return out, nil
			}
		}
	}
	//没有缓存
	data := g.Cfg("routers")
	if data == nil {
		return nil, gerror.New("没有获取到菜单路由数据!")
	}
	var output []*model.RouterOutput
	if err := gconv.Structs(data.Map(), &output); err != nil {
		g.Log().Error(err)
		return nil, err
	}
	//检查数据
	if output != nil && len(output) > 0 {
		if a.routerCache == nil {
			a.routerCache = gcache.New()
		}
	}
	return output, nil
}
