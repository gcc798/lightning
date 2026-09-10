package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAttachmentPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentPageLogic {
	return &AttachmentPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AttachmentPageLogic) AttachmentPage(in *pb.AttachmentPageReq) (*pb.AttachmentPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"status = 0"}
	args := make([]interface{}, 0)
	if in.FileName != "" {
		args = append(args, "%"+in.FileName+"%")
		where = append(where, fmt.Sprintf("file_name like $%d", len(args)))
	}
	if in.FileType != "" {
		args = append(args, in.FileType)
		where = append(where, fmt.Sprintf("file_type = $%d", len(args)))
	}
	if in.BusinessType != "" {
		args = append(args, in.BusinessType)
		where = append(where, fmt.Sprintf("business_type = $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.biz_attachment where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []attachmentRow
	query := `select id, env_id, file_name, file_key, file_size, file_type, file_ext, business_type, business_id, business_field, is_public, access_url, metadata, status, expire_time, create_by, create_time, update_time from public.biz_attachment where ` + whereSQL + ` order by create_time desc, id desc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.AttachmentPageResp{Records: toAttachmentList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
