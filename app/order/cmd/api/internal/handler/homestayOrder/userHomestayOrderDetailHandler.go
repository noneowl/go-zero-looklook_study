package homestayOrder

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"looklook_study/app/order/cmd/api/internal/logic/homestayOrder"
	"looklook_study/app/order/cmd/api/internal/svc"
	"looklook_study/app/order/cmd/api/internal/types"
)

func UserHomestayOrderDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserHomestayOrderDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestayOrder.NewUserHomestayOrderDetailLogic(r.Context(), svcCtx)
		resp, err := l.UserHomestayOrderDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
