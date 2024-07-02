package v1

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "mxshop/api/goods/v1"
	"mxshop/app/goods/srv/internal/domain/dto"
	v1 "mxshop/app/goods/srv/internal/service/v1"
	v12 "mxshop/pkg/common/meta/v1"
)

type goodsServer struct {
	//srv v1.GoodsSrv
	//工厂模式
	srv v1.ServiceFactory
	proto.UnimplementedGoodsServer
}

func NewGoodsServer(srv v1.ServiceFactory) *goodsServer {
	return &goodsServer{
		srv: srv,
	}
}

func ModelToResponse(goods *dto.GoodsDTO) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		//使用外键注意加载进来
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (gs *goodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	opts := v12.ListMeta{
		Page:     int(request.Pages),
		PageSize: int(request.PagePerNums),
	}
	list, err := gs.srv.Goods().List(ctx, opts, request, []string{})
	if err != nil {
		return nil, err
	}
	var ret proto.GoodsListResponse
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, ModelToResponse(item))
	}
	return &ret, nil

}

func (gs *goodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) BannerList(ctx context.Context, req *proto.BannerListReq) (*proto.BannerListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) mustEmbedUnimplementedGoodsServer() {
	//TODO implement me
	panic("implement me")
}
