package handler

import (
	"net/http"

	"github.com/delyr1c/dechoric/src/domain/internal/logic"
	"github.com/delyr1c/dechoric/src/domain/internal/svc"
	"github.com/delyr1c/dechoric/src/domain/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DomainHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDomainLogic(r.Context(), svcCtx)
		resp, err := l.Domain(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
