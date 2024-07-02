package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	v1 "mxshop/app/goods/srv/internal/data_search/v1"
	"mxshop/app/goods/srv/internal/domain/do"
	"mxshop/app/pkg/code"
	"mxshop/pkg/errors"
	"strconv"
)

type goods struct {
	esClient *elastic.Client
}

func NewGoods(esClient *elastic.Client) *goods {
	return &goods{esClient: esClient}
}

// 工厂模式
func newGoods(ds *dataSearch) *goods {
	return &goods{
		esClient: ds.esClient,
	}
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsSearchDO) error {
	_, err := g.esClient.Index().
		Index(goods.GetIndexName()).
		Id(strconv.Itoa(int(goods.ID))).
		BodyJson(&goods).
		Do(ctx)
	return err
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	_, err := g.esClient.Delete().
		Index(do.GoodsSearchDO{}.GetIndexName()).
		Id(strconv.Itoa(int(ID))).
		Refresh("true").
		Do(ctx)
	return err
}

func (g *goods) Update(ctx context.Context, goods *do.GoodsSearchDO) error {
	// 先删除再新增
	err := g.Delete(ctx, uint64(goods.ID))
	if err != nil {
		return err
	}
	err = g.Create(ctx, goods)
	if err != nil {
		return err
	}
	return nil
}

// 将所有涉及到es查询的条件，语句全部抽离出来
func (g *goods) Search(ctx context.Context, req *v1.GoodsFilterRequest) (*do.GoodsSearchDOList, error) {
	// 定义es复合查询 bool查询
	q := elastic.NewBoolQuery()
	if req.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	// IsHot IsNew 默认为false
	if req.IsHot {
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsHot))
	}
	if req.PriceMin > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}
	if req.TopCategory > 0 {
		q = q.Filter(elastic.NewTermsQuery("category_id", req.CategoryIDs...))
	}

	// 分页
	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums < 0:
		req.PagePerNums = 10
	}

	res, err := g.esClient.Search().Index(do.GoodsSearchDO{}.GetIndexName()).Query(q).From(int(req.Pages) * int(req.PagePerNums)).Size(int(req.PagePerNums)).Do(ctx)
	if err != nil {
		return nil, err
	}
	var ret do.GoodsSearchDOList
	ret.TotalCount = res.Hits.TotalHits.Value
	for _, value := range res.Hits.Hits {
		var goodsDo do.GoodsSearchDO
		err = json.Unmarshal(value.Source, &goodsDo)
		if err != nil {
			return nil, errors.WithCode(code.ErrEsUnmarshal, err.Error())
		}
		ret.Items = append(ret.Items, &goodsDo)
	}
	return &ret, nil
}

var _ v1.GoodsStore = (*goods)(nil)
