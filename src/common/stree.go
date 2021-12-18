package common

import (
	"strings"
)

// NodeCommonReq 操作树结构的通用对象

type NodeCommonReq struct {
	Node        string `json:"node"`         // 服务节点名称: 可以是一段式 也可以是二段 inf | inf.mon
	QueryType   int    `json:"query_type"`   // 查询模式
	ForceDelete bool   `json:"force_delete"` // (父节点被删除)子节点是否需要强制删除
}

func (nc *NodeCommonReq) IsExpectedLenFormat(l int) bool {
	return len(strings.Split(nc.Node, ".")) == l
}
