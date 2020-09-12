package dao

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"

	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/model"
)

const (
	author     = "asong"
	project    = "golang_dream"
	mappingTpl = `{
	"mappings":{
		"properties":{
			"id": 				{ "type": "long" },
			"username": 		{ "type": "keyword" },
			"nickname":			{ "type": "text" },
			"phone":			{ "type": "keyword" },
			"age":				{ "type": "long" },
			"ancestral":		{ "type": "text" },
			"identity":         { "type": "text" },
			"update_time":		{ "type": "long" },
			"create_time":		{ "type": "long" }
			}
		}
	}`
	esRetryLimit = 3 //bulk 错误重试机制
)

type UserES struct {
	index   string
	mapping string
	client  *elastic.Client
}

func NewUserES(client *elastic.Client) *UserES {
	index := fmt.Sprintf("%s_%s", author, project)
	userEs := &UserES{
		client:  client,
		index:   index,
		mapping: mappingTpl,
	}

	userEs.init()

	return userEs
}

func (es *UserES) init() {
	ctx := context.Background()

	exists, err := es.client.IndexExists(es.index).Do(ctx)
	if err != nil {
		fmt.Printf("userEs init exist failed err is %s\n", err)
		return
	}

	if !exists {
		_, err := es.client.CreateIndex(es.index).Body(es.mapping).Do(ctx)
		if err != nil {
			fmt.Printf("userEs init failed err is %s\n", err)
			return
		}
	}
}

func (es *UserES) BatchAdd(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchAdd(ctx, user); err != nil {
			fmt.Println("batch add failed ", err)
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchAdd(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
		u.CreateTime = uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
		doc := elastic.NewBulkIndexRequest().Id(strconv.FormatUint(u.ID, 10)).Doc(u)
		req.Add(doc)
	}
	if req.NumberOfActions() < 0 {
		return nil
	}
	if _, err := req.Do(ctx); err != nil {
		return err
	}
	return nil
}

func (es *UserES) BatchUpdate(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchUpdate(ctx, user); err != nil {
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchUpdate(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
		doc := elastic.NewBulkUpdateRequest().Id(strconv.FormatUint(u.ID, 10)).Doc(u)
		req.Add(doc)
	}

	if req.NumberOfActions() < 0 {
		return nil
	}
	if _, err := req.Do(ctx); err != nil {
		return err
	}
	return nil
}

func (es *UserES) BatchDel(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchDel(ctx, user); err != nil {
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchDel(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		doc := elastic.NewBulkDeleteRequest().Id(strconv.FormatUint(u.ID, 10))
		req.Add(doc)
	}

	if req.NumberOfActions() < 0 {
		return nil
	}

	if _, err := req.Do(ctx); err != nil {
		return err
	}
	return nil
}

// 根据id 批量获取
func (es *UserES) MGet(ctx context.Context, IDS []uint64) ([]*model.UserEs, error) {
	userES := make([]*model.UserEs, 0, len(IDS))
	idStr := make([]string, 0, len(IDS))
	for _, id := range IDS {
		idStr = append(idStr, strconv.FormatUint(id, 10))
	}
	resp, err := es.client.Search(es.index).Query(
		elastic.NewIdsQuery().Ids(idStr...)).Size(len(IDS)).Do(ctx)

	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}
	for _, e := range resp.Each(reflect.TypeOf(&model.UserEs{})) {
		us := e.(*model.UserEs)
		userES = append(userES, us)
	}
	return userES, nil
}

func (es *UserES) Search(ctx context.Context, filter *model.EsSearch) ([]*model.UserEs, error) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(filter.MustQuery...)
	boolQuery.MustNot(filter.MustNotQuery...)
	boolQuery.Should(filter.ShouldQuery...)
	boolQuery.Filter(filter.Filters...)

	// 当should不为空时，保证至少匹配should中的一项
	if len(filter.MustQuery) == 0 && len(filter.MustNotQuery) == 0 && len(filter.ShouldQuery) > 0 {
		boolQuery.MinimumShouldMatch("1")
	}

	service := es.client.Search().Index(es.index).Query(boolQuery).SortBy(filter.Sorters...).From(filter.From).Size(filter.Size)
	resp, err := service.Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}
	userES := make([]*model.UserEs, 0)
	for _, e := range resp.Each(reflect.TypeOf(&model.UserEs{})) {
		us := e.(*model.UserEs)
		userES = append(userES, us)
	}
	return userES, nil
}
