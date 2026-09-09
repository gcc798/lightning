package biz

import (
	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	kerrors "github.com/go-kratos/kratos/v2/errors"
)

func badRequest(msg string) error {
	return kerrors.BadRequest("BAD_REQUEST", msg)
}

func unauthorized(msg string) error {
	return kerrors.Unauthorized("UNAUTHORIZED", msg)
}

func forbidden(msg string) error {
	return kerrors.Forbidden("FORBIDDEN", msg)
}

func notFound(msg string) error {
	return kerrors.NotFound("NOT_FOUND", msg)
}

func okReply() *v1.MessageReply {
	return &v1.MessageReply{Message: "ok"}
}

func requireFound[T any](item *T, msg string) (*T, error) {
	if item == nil {
		return nil, notFound(msg)
	}
	return item, nil
}
