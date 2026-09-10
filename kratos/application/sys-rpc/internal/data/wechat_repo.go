package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	entpkg "github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent/user"
)

func (r *Resources) ResolveXcxUser(ctx context.Context, phonenumber, wxCode string) (*v1.UserItem, error) {
	if strings.TrimSpace(phonenumber) == "" {
		return nil, errors.New("手机号不能为空")
	}
	if strings.TrimSpace(wxCode) == "" {
		return nil, errors.New("微信code不能为空")
	}
	if r.WeChat == nil {
		return nil, errors.New("微信配置缺失")
	}
	wxResp, err := r.WeChat.Code2Session(ctx, wxCode)
	if err != nil {
		return nil, err
	}
	found, err := r.Ent.User.Query().
		Where(
			user.Phonenumber(phonenumber),
			user.DeletedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		if !entpkg.IsNotFound(err) {
			return nil, err
		}
		now := time.Now()
		nickName := buildXcxNickname(phonenumber)
		created, createErr := r.Ent.User.Create().
			SetID(nextID()).
			SetUserName(phonenumber).
			SetNickName(nickName).
			SetUserType(1).
			SetPhonenumber(phonenumber).
			SetStatus(0).
			SetOpenID(wxResp.OpenID).
			SetUnionID(wxResp.UnionID).
			SetCreatedTime(now).
			SetUpdatedTime(now).
			Save(ctx)
		if createErr != nil {
			return nil, fmt.Errorf("创建微信用户失败: %w", createErr)
		}
		return userEntityToItem(created), nil
	}
	if _, err := r.Ent.User.UpdateOneID(found.ID).
		SetOpenID(wxResp.OpenID).
		SetUnionID(wxResp.UnionID).
		SetUpdatedTime(time.Now()).
		Save(ctx); err != nil {
		return nil, err
	}
	found.OpenID = wxResp.OpenID
	found.UnionID = wxResp.UnionID
	found.UpdatedTime = time.Now()
	return userEntityToItem(found), nil
}

func buildXcxNickname(phonenumber string) string {
	phonenumber = strings.TrimSpace(phonenumber)
	if len(phonenumber) >= 4 {
		return "微信用户" + phonenumber[len(phonenumber)-4:]
	}
	return "微信用户"
}
