package model

import (
	"github.com/olivere/elastic/v7"
)

type UserEs struct {
	ID         uint64 `json:"id,omitempty" mapstructure:"id"`
	Username   string `json:"username,omitempty" mapstructure:"username"`
	Nickname   string `json:"nickname,omitempty" mapstructure:"nickname"`
	Phone      string `json:"phone,omitempty" mapstructure:"phone"`
	Age        uint64 `json:"age,omitempty" mapstructure:"age"`
	Ancestral  string `json:"ancestral,omitempty" mapstructure:"Ancestral"`
	Identity   string `json:"identity,omitempty" mapstructure:"identity"`
	UpdateTime uint64 `json:"update_time,omitempty" mapstructure:"update_time"`
	CreateTime uint64 `json:"create_time,omitempty" mapstructure:"create_time"`
}

type SearchRequest struct {
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Identity  string `json:"identity"`
	Ancestral string `json:"ancestral"`
	Num       int    `json:"num"`
	Size      int    `json:"size"`
}

//bool query 条件
type EsSearch struct {
	MustQuery    []elastic.Query
	MustNotQuery []elastic.Query
	ShouldQuery  []elastic.Query
	Filters      []elastic.Query
	Sorters      []elastic.Sorter
	From         int //分页
	Size         int
}

func (r *SearchRequest) ToFilter() *EsSearch {
	var search EsSearch
	if len(r.Nickname) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("nickname", r.Nickname))
	}
	if len(r.Phone) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewTermsQuery("phone", r.Phone))
	}
	if len(r.Ancestral) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("ancestral", r.Ancestral))
	}
	if len(r.Identity) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("identity", r.Identity))
	}

	if search.Sorters == nil {
		search.Sorters = append(search.Sorters, elastic.NewFieldSort("create_time").Desc())
	}

	search.From = (r.Num - 1) * r.Size
	search.Size = r.Size
	return &search
}
