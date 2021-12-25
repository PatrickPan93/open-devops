package common

import (
	"log"

	"github.com/ning1875/inverted-index/labels"
)

// ResourceQueryReq 请求结构体
type ResourceQueryReq struct {
	ResourceType string          `json:"resource_type" binding:"required"`
	Labels       []*SingleTagReq `json:"labels" binding:"required"`
	TargetLabel  string          `json:"target_label"`
}

// QueryResponse 响应结构体
type QueryResponse struct {
	Code        int         `json:"code"`
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	PageCount   int         `json:"page_count"`
	TotalCount  int         `json:"total_count"`
	Result      interface{} `json:"result"`
}

// SingleTagReq 查询标签
type SingleTagReq struct {
	Key   string `json:"key" binding:"required"`   // 标签名称
	Value string `json:"value" binding:"required"` // 标签值, 可以是正则表达式
	Type  int    `json:"type" binding:"required"`
	/* 查询类型1-4
	1: MatchEqual,
	2: MatchNotEqual,
	3: MatchRegexp,
	4: MatchNotRegexp,
	*/

}

// FormatLabelMatcher 将前段请求转化为[]*labels.Matcher
func FormatLabelMatcher(ls []*SingleTagReq) []*labels.Matcher {
	matchers := make([]*labels.Matcher, 0)
	for _, i := range ls {
		// 通过MatchMap返回这次的match type值
		mType, ok := labels.MatchMap[i.Type]
		if !ok {
			log.Printf("common.FormatLabelMatcher: querty type %d is not supported", i.Type)
			continue
		}
		// gen NewMatcher via mType, key(label name), value(label values | Regex)
		matchers = append(matchers, labels.MustNewMatcher(mType, i.Key, i.Value))
	}
	return matchers
}
