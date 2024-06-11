package logic

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/domainApi/internal/svc"
	"github.com/delyr1c/dechoric/src/domain/domainApi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DomainApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDomainApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DomainApiLogic {
	return &DomainApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DomainApiLogic) DomainApi(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
