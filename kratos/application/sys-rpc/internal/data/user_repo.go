package data

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	entpkg "github.com/gcc798/lightning/kratos/application/sys-rpc/ent"
	"golang.org/x/crypto/bcrypt"
)

func userEntityToItem(item *entpkg.User) *v1.UserItem {
	if item == nil {
		return nil
	}
	return &v1.UserItem{
		UserId:      item.ID,
		UserName:    item.UserName,
		NickName:    item.NickName,
		UserType:    item.UserType,
		Email:       item.Email,
		Phonenumber: item.Phonenumber,
		Sex:         item.Sex,
		Avatar:      item.Avatar,
		Status:      item.Status,
		Sort:        item.Sort,
		LoginIp:     item.LoginIP,
		LoginDate:   item.LoginDate,
		OpenId:      item.OpenID,
		UnionId:     item.UnionID,
		Remark:      item.Remark,
		CreateBy:    item.CreateBy,
		UpdateBy:    item.UpdateBy,
		CreatedTime: formatTime(item.CreatedTime),
		UpdatedTime: formatTime(item.UpdatedTime),
	}
}

func (r *Resources) activeUsers(ctx context.Context) ([]*entpkg.User, error) {
	items, err := r.Ent.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*entpkg.User, 0, len(items))
	for _, item := range items {
		if item.DeletedAt == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *Resources) GetUser(ctx context.Context, id int64) (*v1.UserItem, error) {
	item, err := r.Ent.User.Get(ctx, id)
	if err != nil {
		if entpkg.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if item.DeletedAt != nil {
		return nil, nil
	}
	return userEntityToItem(item), nil
}

func ensureUniqueUser(ctx context.Context, r *Resources, userID int64, username, phone, email string) error {
	items, err := r.activeUsers(ctx)
	if err != nil {
		return err
	}
	for _, item := range items {
		if userID > 0 && item.ID == userID {
			continue
		}
		if username != "" && item.UserName == username {
			return errors.New("username already exists")
		}
		if phone != "" && item.Phonenumber == phone {
			return errors.New("phonenumber already exists")
		}
		if email != "" && item.Email == email {
			return errors.New("email already exists")
		}
	}
	return nil
}

func (r *Resources) CreateUser(ctx context.Context, req *v1.CreateUserRequest) error {
	if err := ensureUniqueUser(ctx, r, 0, req.GetUserName(), req.GetPhonenumber(), req.GetEmail()); err != nil {
		return err
	}
	now := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	operator := currentOperatorID(ctx)
	_, err = r.Ent.User.Create().
		SetID(nextID()).
		SetUserName(req.GetUserName()).
		SetNickName(req.GetNickName()).
		SetUserType(req.GetUserType()).
		SetEmail(req.GetEmail()).
		SetPhonenumber(req.GetPhonenumber()).
		SetSex(req.GetSex()).
		SetAvatar(req.GetAvatar()).
		SetPassword(string(hash)).
		SetStatus(req.GetStatus()).
		SetRemark(req.GetRemark()).
		SetCreateBy(operator).
		SetUpdateBy(operator).
		SetCreatedTime(now).
		SetUpdatedTime(now).
		Save(ctx)
	return err
}

func (r *Resources) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) error {
	item, err := r.Ent.User.Get(ctx, req.GetUserId())
	if err != nil {
		return err
	}
	if item.DeletedAt != nil {
		return errors.New("user not found")
	}
	if err := ensureUniqueUser(ctx, r, req.GetUserId(), req.GetUserName(), req.GetPhonenumber(), req.GetEmail()); err != nil {
		return err
	}
	operator := currentOperatorID(ctx)
	_, err = r.Ent.User.UpdateOneID(req.GetUserId()).
		SetUserName(req.GetUserName()).
		SetNickName(req.GetNickName()).
		SetUserType(req.GetUserType()).
		SetEmail(req.GetEmail()).
		SetPhonenumber(req.GetPhonenumber()).
		SetSex(req.GetSex()).
		SetAvatar(req.GetAvatar()).
		SetStatus(req.GetStatus()).
		SetRemark(req.GetRemark()).
		SetUpdateBy(operator).
		SetUpdatedTime(time.Now()).
		Save(ctx)
	return err
}

func (r *Resources) DeleteUsers(ctx context.Context, ids ...int64) error {
	now := time.Now()
	operator := currentOperatorID(ctx)
	return r.withTx(ctx, func(tx *entpkg.Tx) error {
		for _, id := range ids {
			item, err := tx.User.Get(ctx, id)
			if err != nil {
				return err
			}
			if item.DeletedAt != nil {
				return errors.New("user not found")
			}
			if _, err := tx.User.UpdateOneID(id).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
				return err
			}
			userRoles, err := tx.UserRole.Query().All(ctx)
			if err != nil {
				return err
			}
			for _, item := range userRoles {
				if item.DeletedAt == nil && item.UserID == id {
					if _, err := tx.UserRole.UpdateOneID(item.ID).SetDeletedAt(now).SetUpdateBy(operator).SetUpdatedTime(now).Save(ctx); err != nil && !entpkg.IsNotFound(err) {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (r *Resources) PageUsers(ctx context.Context, req *v1.PageUsersRequest) (*v1.PageUsersReply, error) {
	items, err := r.activeUsers(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]*v1.UserItem, 0, len(items))
	for _, item := range items {
		if req.GetUsername() != "" && !strings.Contains(item.UserName, req.GetUsername()) {
			continue
		}
		if req.GetPhonenumber() != "" && !strings.Contains(item.Phonenumber, req.GetPhonenumber()) {
			continue
		}
		if req.Status != nil && item.Status != req.GetStatus() {
			continue
		}
		filtered = append(filtered, userEntityToItem(item))
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].GetCreatedTime() == filtered[j].GetCreatedTime() {
			return filtered[i].GetUserId() < filtered[j].GetUserId()
		}
		return filtered[i].GetCreatedTime() > filtered[j].GetCreatedTime()
	})
	pageNum, pageSize := buildPage(req.GetPageNum(), req.GetPageSize())
	start, end := pageBounds(len(filtered), pageNum, pageSize)
	return &v1.PageUsersReply{List: filtered[start:end], Total: int64(len(filtered)), PageNum: pageNum, PageSize: pageSize}, nil
}

func (r *Resources) ResetPassword(ctx context.Context, userID int64, newPassword string) error {
	item, err := r.Ent.User.Get(ctx, userID)
	if err != nil {
		return err
	}
	if item.DeletedAt != nil {
		return errors.New("user not found")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.Ent.User.UpdateOneID(userID).
		SetPassword(string(hash)).
		SetUpdateBy(currentOperatorID(ctx)).
		SetUpdatedTime(time.Now()).
		Save(ctx)
	return err
}

func (r *Resources) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	if userID <= 0 {
		return errors.New("current user not found")
	}
	item, err := r.Ent.User.Get(ctx, userID)
	if err != nil {
		return err
	}
	if item.DeletedAt != nil {
		return errors.New("user not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(oldPassword)) != nil {
		return errors.New("old password is invalid")
	}
	return r.ResetPassword(ctx, userID, newPassword)
}

func (r *Resources) ImportUsers(ctx context.Context, users []*v1.CreateUserRequest) (*v1.ImportUsersReply, error) {
	reply := &v1.ImportUsersReply{}
	for _, item := range users {
		found, _, err := r.FindUserByAccount(ctx, item.GetUserName())
		if err != nil {
			return nil, err
		}
		if found != nil {
			reply.FailCount++
			reply.Errors = append(reply.Errors, "duplicate username: "+item.GetUserName())
			continue
		}
		if err := r.CreateUser(ctx, item); err != nil {
			reply.FailCount++
			reply.Errors = append(reply.Errors, err.Error())
			continue
		}
		reply.SuccessCount++
	}
	return reply, nil
}
