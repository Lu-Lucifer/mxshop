package v1

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	proto "mxshop/api/goods/v1"
	v1 "mxshop/app/goods/srv/internal/data/v1"
	v12 "mxshop/app/goods/srv/internal/data_search/v1"
	"mxshop/app/goods/srv/internal/domain/do"
	"mxshop/app/goods/srv/internal/domain/dto"
	metav1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/log"
	"sync"
)

type GoodsSrv interface {
	//根据id查找商品详情
	Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error)
	//商品列表
	List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error)
	Create(ctx context.Context, goods *dto.GoodsDTO) error
	Update(ctx context.Context, goods *dto.GoodsDTO) error
	Delete(ctx context.Context, ID uint64) error
	//批量查询商品
	BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error)
}

//type goodsService struct {
//	data         v1.GoodsStore
//	categoryData v1.CategoryStore
//	searchData   v12.GoodsStore
//	brandData    v1.BrandsStore
//}

// 工厂模式
type goodsService struct {
	data       v1.DataFactory
	searchData v12.SearchFactory
}

//func NewGoodsService(data v1.GoodsStore, categoryData v1.CategoryStore, searchData v12.GoodsStore, brandData v1.BrandsStore) *goodsService {
//	return &goodsService{
//		data:         data,
//		categoryData: categoryData,
//		searchData:   searchData,
//		brandData:    brandData}
//}

func NewGoodsService(dataFactory v1.DataFactory, searchData v12.SearchFactory) *goodsService {
	return &goodsService{
		data:       dataFactory,
		searchData: searchData,
	}
}

// 工厂模式
func newGoods(srv *service) *goodsService {
	return &goodsService{
		data:       srv.data,
		searchData: srv.dataSearch,
	}

}

// 使用工厂模式
func (gs *goodsService) Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error) {
	//goodsDO, err := gs.data.Get(ctx, ID)
	goodsDO, err := gs.data.Goods().Get(ctx, ID)
	if err != nil {
		log.Errorf("data,Get err: %v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *goodsDO,
	}, nil
}

// 递归遍历分类结构
func retrieveIDs(category *do.CategoryDO) []uint64 {
	ids := make([]uint64, 0)
	if category == nil || category.ID == 0 {
		return ids
	}
	ids = append(ids, uint64(category.ID))
	for _, child := range category.SubCategory {
		subIds := retrieveIDs(child)
		ids = append(ids, subIds...)
	}
	return ids
}

func (gs *goodsService) List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error) {
	searchReq := v12.GoodsFilterRequest{
		GoodsFilterRequest: req,
	}
	// category
	if req.TopCategory > 0 {
		//先通过category id查询该分类信息
		categoryDO, err := gs.data.Category().Get(ctx, uint64(req.TopCategory))
		if err != nil {
			log.Errorf("categoryData.Get err: %v", err)
			return nil, err
		}
		var ids []interface{}
		for _, id := range retrieveIDs(categoryDO) {
			ids = append(ids, id)
		}
		searchReq.CategoryIDs = ids
	}
	goodsList, err := gs.searchData.Goods().Search(ctx, &searchReq)
	if err != nil {
		log.Errorf("searchData.Search err: %v", err)
		return nil, err
	}
	// 用从es中搜索的商品，提取出ids，然后去数据库中查询出这些商品
	goodsIDs := []uint64{}
	for _, value := range goodsList.Items {
		goodsIDs = append(goodsIDs, uint64(value.ID))
	}
	goodsDOList, err := gs.data.Goods().ListByIDs(ctx, goodsIDs, orderby)
	if err != nil {
		log.Errorf("data.ListByIDs err: %v", err)
		return nil, err
	}
	var ret dto.GoodsDTOList
	ret.TotalCount = goodsDOList.TotalCount
	for _, goodsDO := range goodsDOList.Items {
		ret.Items = append(ret.Items, &dto.GoodsDTO{
			*goodsDO,
		})
	}
	return &ret, nil
}

func (gs *goodsService) Create(ctx context.Context, goods *dto.GoodsDTO) error {
	/*
		创建商品
		数据先写入mysql，再写入es
	*/
	//先判断CategoryID，BrandsID是否存在
	_, err := gs.data.Category().Get(ctx, uint64(goods.CategoryID))
	if err != nil {
		return err
	}
	_, err = gs.data.Brands().Get(ctx, uint64(goods.BrandsID))
	if err != nil {
		return err
	}

	/*
		之前写入es的方案是给gorm添加钩子函数AfterCreate（gorm会保证事务）
		这里用另一种方案，因为涉及到mysql，es，要用分布式事务：之前基于可靠消息最终一致性，但这个方案很重，每次需要发送事务消息
	*/

	/*
		开启事务后，要非常小心，一定要有回滚和提交操作；
		（实际开发中严重的坑）为避免代码运行到一半，程序崩掉了，后面的代码执行不到，无法执行rollback或commit，事务的数据会被锁住，这里需要用捕获异常处理，执行回滚
		事务开启后，defer语句一定会在return之前执行，可以保证捕获异常处理，执行回滚
		这种事务只针对对数据一致性不高的场景
	*/
	txn := gs.data.Begin() //开启事务
	defer func() {
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()
	err = gs.data.Goods().CreateInTxn(ctx, txn, &goods.GoodsDO)
	if err != nil {
		log.Errorf("data.CreateInTxn err: %v", err)
		txn.Rollback()
		return err
	}
	searchDO := do.GoodsSearchDO{
		ID:          goods.ID,
		CategoryID:  goods.CategoryID,
		BrandsID:    goods.BrandsID,
		OnSale:      goods.OnSale,
		ShipFree:    goods.ShipFree,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		Name:        goods.Name,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavNum:      goods.FavNum,
		MarketPrice: goods.MarketPrice,
		GoodsBrief:  goods.GoodsBrief,
		ShopPrice:   goods.ShopPrice,
	}
	err = gs.searchData.Goods().Create(ctx, &searchDO)
	if err != nil {
		log.Errorf("searchData.Create err: %v", err)
		txn.Rollback()
		return err
	}
	txn.Commit() //提交事务
	return nil
}

func (gs *goodsService) Update(ctx context.Context, goods *dto.GoodsDTO) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	//第一种常规做法
	//ds, err := gs.data.ListByIDs(ctx, ids, []string{})
	//if err != nil {
	//	return nil, err
	//}
	//var ret []*dto.GoodsDTO
	//for _, goodsDO := range ds.Items {
	//	ret = append(ret, &dto.GoodsDTO{
	//		*goodsDO,
	//	})
	//}
	//return ret, nil

	//第二中方案：使用go-zero的mapreduce来处理并发跳用
	var fns []func() error
	var ret []*dto.GoodsDTO
	var mu sync.Locker
	for _, id := range ids {
		fns = append(fns, func() error {
			//调用自己内部的get方法处理
			goodsDTO, err := gs.Get(ctx, id)
			//多个goroutine执行时，切片，map都是线程不安全的，要加锁；或者用sync map处理
			mu.Lock()
			ret = append(ret, goodsDTO)
			mu.Unlock()
			return err
		})
	}

	err := mr.Finish(fns...)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

var _ GoodsSrv = (*goodsService)(nil)
