package logic

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/internal/svc"
	"github.com/delyr1c/dechoric/src/domain/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DomainLogic {
	return &DomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DomainLogic) Domain(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
