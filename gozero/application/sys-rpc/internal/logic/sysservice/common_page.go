package sysservicelogic

import "github.com/gcc798/lightning/application/sys-rpc/pb"

func normalizePage(pageNum, pageSize int64) (int64, int64) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 1000 {
		pageSize = 1000
	}
	return pageNum, pageSize
}

func toPageInfo(total, pageNum, pageSize int64) *pb.PageInfo {
	pages := int64(0)
	if pageSize > 0 {
		pages = (total + pageSize - 1) / pageSize
	}
	return &pb.PageInfo{Total: total, Size: pageSize, Current: pageNum, Pages: pages}
}
