package controller

import (
	"context"

	"focus-single/api/v1"
	"focus-single/internal/service/category"
)

// 栏目管理
var Category = cCategory{}

type cCategory struct{}

func (a *cCategory) Tree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error) {
	res = &v1.CategoryTreeRes{}
	res.List, err = category.GetTree(ctx, req.ContentType)
	return
}
