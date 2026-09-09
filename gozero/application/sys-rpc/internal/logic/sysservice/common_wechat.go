package sysservicelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	gzsqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	miniProgramUserType     int32 = 1
	defaultMiniProgramOrgID int64 = 1880159541355577346
)

type wechatCode2SessionResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func wechatCode2Session(appId, secret, wxCode string) (*wechatCode2SessionResp, error) {
	apiURL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		url.QueryEscape(appId), url.QueryEscape(secret), url.QueryEscape(wxCode))
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("调用微信接口失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取微信响应失败: %w", err)
	}
	var result wechatCode2SessionResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析微信响应失败: %w", err)
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信接口错误: %s", result.ErrMsg)
	}
	if result.OpenId == "" {
		return nil, fmt.Errorf("获取微信OpenID失败")
	}
	return &result, nil
}

func authenticateXcx(ctx context.Context, svcCtx *svc.ServiceContext, phonenumber, code, wxCode string) (*userAuthRow, error) {
	if phonenumber == "" || code == "" || wxCode == "" {
		return nil, fmt.Errorf("手机号、验证码和微信code不能为空")
	}
	if !svcCtx.Config.Wechat.Enabled {
		return nil, fmt.Errorf("微信小程序登录未启用")
	}
	wxResp, err := wechatCode2Session(svcCtx.Config.Wechat.AppId, svcCtx.Config.Wechat.Secret, wxCode)
	if err != nil {
		return nil, err
	}

	var row userAuthRow
	err = svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, org_id, user_name, nick_name, user_type, email, phonenumber, avatar, password, status
		from public.s_user
		where phonenumber = $1
		limit 1
	`, phonenumber)
	if err != nil && !errors.Is(err, gzsqlx.ErrNotFound) {
		return nil, fmt.Errorf("查询用户失败")
	}
	if errors.Is(err, gzsqlx.ErrNotFound) {
		now := time.Now().Unix()
		if _, err := svcCtx.DB.ExecCtx(ctx, `
			insert into public.s_user (user_name, nick_name, user_type, org_id, phonenumber, status, sex, open_id, union_id, login_date, created_time, updated_time)
			values ($1, $2, $3, $4, $5, 0, 2, $6, $7, $8, now(), now())
		`, phonenumber, phonenumber,
			int64(miniProgramUserType), defaultMiniProgramOrgID, phonenumber,
			sql.NullString{String: wxResp.OpenId, Valid: true},
			sql.NullString{String: wxResp.UnionId, Valid: wxResp.UnionId != ""},
			now); err != nil {
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}
		if err := svcCtx.DB.QueryRowCtx(ctx, &row, `
			select id, org_id, user_name, nick_name, user_type, email, phonenumber, avatar, password, status
			from public.s_user
			where phonenumber = $1
			limit 1
		`, phonenumber); err != nil {
			return nil, fmt.Errorf("查询新用户失败: %w", err)
		}
	} else if _, err := svcCtx.DB.ExecCtx(ctx, `
		update public.s_user
		set open_id = $1, union_id = $2, login_date = $3, updated_time = now()
		where id = $4
	`, sql.NullString{String: wxResp.OpenId, Valid: true},
		sql.NullString{String: wxResp.UnionId, Valid: wxResp.UnionId != ""},
		time.Now().Unix(), row.Id); err != nil {
		return nil, fmt.Errorf("更新微信用户信息失败: %w", err)
	}
	if row.Status != 0 {
		return nil, fmt.Errorf("用户已被停用")
	}
	return &row, nil
}

func authenticateWechat(ctx context.Context, svcCtx *svc.ServiceContext, wxCode string) (*userAuthRow, error) {
	if wxCode == "" {
		return nil, fmt.Errorf("微信code不能为空")
	}
	if !svcCtx.Config.Wechat.Enabled {
		return nil, fmt.Errorf("微信小程序登录未启用")
	}
	wxResp, err := wechatCode2Session(svcCtx.Config.Wechat.AppId, svcCtx.Config.Wechat.Secret, wxCode)
	if err != nil {
		return nil, err
	}

	var row userAuthRow
	err = svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, org_id, user_name, nick_name, user_type, email, phonenumber, avatar, password, status
		from public.s_user
		where open_id = $1
		limit 1
	`, wxResp.OpenId)
	if err != nil && !errors.Is(err, gzsqlx.ErrNotFound) {
		return nil, fmt.Errorf("查询用户失败")
	}
	if errors.Is(err, gzsqlx.ErrNotFound) {
		userName := "wx_" + wxResp.OpenId[:minInt(16, len(wxResp.OpenId))]
		now := time.Now().Unix()
		if _, err := svcCtx.DB.ExecCtx(ctx, `
			insert into public.s_user (user_name, nick_name, user_type, org_id, status, sex, open_id, union_id, login_date, created_time, updated_time)
			values ($1, $2, $3, $4, 0, 2, $5, $6, $7, now(), now())
		`, userName, "微信用户",
			int64(miniProgramUserType), defaultMiniProgramOrgID,
			sql.NullString{String: wxResp.OpenId, Valid: true},
			sql.NullString{String: wxResp.UnionId, Valid: wxResp.UnionId != ""},
			now); err != nil {
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}
		if err := svcCtx.DB.QueryRowCtx(ctx, &row, `
			select id, org_id, user_name, nick_name, user_type, email, phonenumber, avatar, password, status
			from public.s_user
			where open_id = $1
			limit 1
		`, wxResp.OpenId); err != nil {
			return nil, fmt.Errorf("查询新用户失败: %w", err)
		}
	} else if wxResp.UnionId != "" {
		_, _ = svcCtx.DB.ExecCtx(ctx, `update public.s_user set union_id = $1, login_date = $2, updated_time = now() where id = $3`,
			sql.NullString{String: wxResp.UnionId, Valid: true}, time.Now().Unix(), row.Id)
	}
	if row.Status != 0 {
		return nil, fmt.Errorf("用户已被停用")
	}
	return &row, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
