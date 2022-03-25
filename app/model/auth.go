package model

//LoginBodyReq 登录请求报文体
type LoginBodyReq struct {
	DeviceCode string `json:"deviceCode"`                                             //设备代码(非必传)
	Account    string `json:"account" v:"required|length:4,20#请输入账号|账号长度为:min到:max位"` //账号
	Password   string `json:"password" v:"required|min-length:6#请输入密码|密码长度不够"`        //密码
	VerifyCode string `json:"verifyCode"`                                             //验证码
	VerifyKey  string `json:"verifyKey"`                                              //验证码键
}

//LoginBodyRes 登录响应报文体
type LoginBodyRes struct {
	Token string `json:"token"` //登录令牌
}

//UserInfo 用户信息
type UserInfo struct {
	Id      uint64 `json:"id"`      //用户ID
	Account string `json:"account"` //用户账号
	Nick    string `json:"nick"`    //用户昵称
}

// UserInfoOutput 认证用户信息
type UserInfoOutput struct {
	UserInfo //用户信息

	Roles       []string `json:"roles"`       //角色集合
	Permissions []string `json:"permissions"` //权限集合
}

// RouterOutput 路由配置
type RouterOutput struct {
	Name       string          `json:"name"`                    //路由名字
	Path       string          `json:"path"`                    //路由地址
	Hidden     bool            `json:"hidden" d:"false"`        //是否隐藏路由
	Redirect   string          `json:"redirect" d:"noRedirect"` //重定向地址(当设置为noRedirect时该路由在面包屑导航中不可被点击)
	Component  string          `json:"component"`               //组件地址
	Query      string          `json:"query"`                   //路由参数(如:{"id": 1, "name": "ry"})
	AlwaysShow bool            `json:"alwaysShow" d:"false"`    //当路由下children大于1时，自动会变成嵌套的模式
	Meta       *MetaOutput     `json:"meta"`                    //其它元素
	Children   []*RouterOutput `json:"children"`                //子路由集合
}

// MetaOutput 路由显示信息
type MetaOutput struct {
	Title    string `json:"title"`            //路由在侧边栏和面包屑中展示的名字
	Icon     string `json:"icon"`             //路由的图标(对应路径src/assets/icons/svg)
	NoCache  bool   `json:"noCache" d:"true"` //设置为true，则不会被 <keep-alive>缓存
	MenuType string `json:"menuType" d:"C"`   //类型(M目录 C菜单 F按钮)
	Link     string `json:"link"`             //内链地址(http(s)://开头)
}
