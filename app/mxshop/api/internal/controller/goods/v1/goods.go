package goods

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	proto "mxshop/api/goods/v1"
	"mxshop/app/mxshop/api/internal/domain/request"
	"mxshop/app/mxshop/api/internal/service"
	gin2 "mxshop/app/pkg/translator/gin"
	"mxshop/pkg/common/core"
)

type goodsController struct {
	srv   service.ServiceFactory
	trans ut.Translator
}

func NewGoodsController(srv service.ServiceFactory, trans ut.Translator) *goodsController {
	return &goodsController{
		srv:   srv,
		trans: trans,
	}

}

func (gc *goodsController) List(ctx *gin.Context) {
	var req request.GoodsFilter
	if err := ctx.ShouldBind(&req); err != nil {
		gin2.HandleValidatorError(ctx, err, gc.trans)
		return
	}
	goodsFilterRequest := proto.GoodsFilterRequest{
		PriceMin:    req.PriceMin,
		PriceMax:    req.PriceMax,
		IsHot:       req.IsHot,
		IsNew:       req.IsNew,
		IsTab:       req.IsTab,
		TopCategory: req.TopCategory,
		Pages:       req.Pages,
		PagePerNums: req.PagePerNums,
		KeyWords:    req.KeyWords,
		Brand:       req.Brand,
	}
	goodsListResponse, err := gc.srv.Goods().List(ctx, &goodsFilterRequest)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	//ctx.JSON(http.StatusOK, goodsListResponse)
	core.WriteResponse(ctx, nil, goodsListResponse)
}
