package iam

// LoginResponse 登录响应
//
//	@Description	登录成功后返回的 Token 信息和用户信息
type LoginResponse struct {
	AccessToken      string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 访问令牌（短期有效）
	RefreshToken     string    `json:"refresh_token" example:"dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4="`       // 刷新令牌（长期有效）
	ExpiresIn        int64     `json:"expires_in" example:"1800"`                                      // AccessToken 过期时间（秒）
	RefreshExpiresIn int64     `json:"refresh_expires_in" example:"604800"`                            // RefreshToken 过期时间（秒）
	ClientID         string    `json:"client_id,omitempty" example:"web-admin"`                        // 客户端ID
	Scope            string    `json:"scope,omitempty" example:""`                                     // 令牌权限
	OpenID           string    `json:"openId,omitempty" example:"openid"`                              // 微信OpenID
	UserInfo         *UserInfo `json:"user_info"`                                                      // 用户信息
}

// RefreshTokenRequest 刷新令牌请求
//
//	@Description	使用 RefreshToken 刷新 AccessToken 的请求参数
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required" example:"dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4="` // 刷新令牌
	ClientID     string `json:"clientId" binding:"required" example:"web-admin"`                            // 客户端ID
}

// RefreshTokenResponse 刷新令牌响应
//
//	@Description	刷新成功后返回的新 Token 对
type RefreshTokenResponse struct {
	AccessToken      string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 新的访问令牌
	RefreshToken     string `json:"refresh_token" example:"bmV3IHJlZnJlc2ggdG9rZW4="`               // 新的刷新令牌（轮换）
	ExpiresIn        int64  `json:"expires_in" example:"1800"`                                      // AccessToken 过期时间（秒）
	RefreshExpiresIn int64  `json:"refresh_expires_in" example:"604800"`                            // RefreshToken 过期时间（秒）
}

// UserInfo 用户信息
//
//	@Description	用户基本信息
type UserInfo struct {
	UserId      int64  `json:"userId" example:"1"`                              // 用户ID
	Username    string `json:"username" example:"admin"`                        // 用户名
	Nickname    string `json:"nickname" example:"系统管理员"`                        // 昵称
	Phonenumber string `json:"phonenumber" example:"13800138000"`               // 手机号
	Email       string `json:"email" example:"admin@example.com"`               // 邮箱
	Avatar      string `json:"avatar" example:"https://example.com/avatar.jpg"` // 头像URL
	OrgID       int64  `json:"orgId,omitempty" example:"1"`                     // 组织ID
	UserType    int32  `json:"userType" example:"0"`                            // 用户类型：0系统用户 1微信用户 2APP用户
	OpenID      string `json:"openId,omitempty" example:"openid"`               // 微信OpenID
	UnionID     string `json:"unionId,omitempty" example:"unionid"`             // 微信UnionID
}

// XcxUserInfo 小程序用户信息
type XcxUserInfo struct {
	UserID       int64  `json:"userId"`
	OrgID        int64  `json:"orgId"`
	Phonenumber  string `json:"phonenumber"`
	OpenID       string `json:"openId"`
	UnionID      string `json:"unionId"`
	UserName     string `json:"userName"`
	NickName     string `json:"nickName"`
	Sex          string `json:"sex"`
	HeadPortrait string `json:"headPortrait"`
	RoleName     string `json:"roleName"`
	RoleKey      string `json:"roleKey"`
}
