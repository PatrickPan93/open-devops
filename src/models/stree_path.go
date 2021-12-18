package models

import (
	"fmt"
	"log"
	"open-devops/src/common"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

const (
	streeDB       = "stree"
	queryAllPByG  = 1
	queryALlPAByG = 2
	queryAllAByPG = 3

	group       = 1
	department  = 2
	application = 3

	originPath = "0"
)

// StreePath table stree_path
type StreePath struct {
	Id       int64  `json:"id"`
	Level    int64  `json:"level"`
	Path     string `json:"path"`
	NodeName string `json:"node_name"`
}

func (sp *StreePath) addOne() (int64, error) {
	return DB[streeDB].InsertOne(sp)
}

// GetOne query one record base on query param
func (sp *StreePath) getOne() (*StreePath, error) {
	has, err := DB[streeDB].Get(sp)
	if err != nil {
		return nil, errors.Wrap(err, "GetOne: failed to get one streePath record")
	}
	if !has {
		return nil, nil
	}
	return sp, nil
}

// CheckExist Check streePath if it's exists
func (sp *StreePath) CheckExist() (bool, error) {
	return DB[streeDB].Exist(sp)
}

// StreePathGet 函数的get
func StreePathGetOne(where string, args ...interface{}) (*StreePath, error) {
	var obj StreePath
	has, err := DB[streeDB].Where(where, args...).Get(&obj)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("StreePathGet: failed to get path using where '%s,%v'\n", where, args))
	}
	if !has {
		return nil, nil
	}
	return &obj, nil

}

// StreePathGetMany 根据条件获取多条记录
func StreePathGetMany(where string, args ...interface{}) ([]StreePath, error) {
	var objs []StreePath
	err := DB[streeDB].Where(where, args...).Find(&objs)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("StreePathGet: failed to get path using where '%s,%v'\n", where, args))
	}

	return objs, nil

}

func StreePathQuery(req *common.NodeCommonReq) {
	// whatever you want to query. we must make sure that g exists.
	var dbg *StreePath
	// split by . , get g name from index 0
	g := strings.Split(req.Node, ".")[0]
	nodeG := StreePath{
		Level:    group,
		Path:     originPath,
		NodeName: g,
	}
	dbg, err := nodeG.getOne()
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	if dbg == nil {
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: g does not exist '%s'", req.Node)))
		return
	}

	switch req.QueryType {
	case queryAllPByG:
		if req.IsExpectedLenFormat(1) {
			pathP := fmt.Sprintf("/%d", dbg.Id)
			whereStr := "level=? and path=?"
			sps, err := StreePathGetMany(whereStr, department, pathP)
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}
			if len(sps) == 0 {
				log.Printf("%+v", errors.New(
					fmt.Sprintf(
						"StreePathQuery: no p under g '%s'", dbg.NodeName)))
				return
			}
			var res []string
			for _, i := range sps {
				res = append(res, i.NodeName)
			}
			log.Println(res)
			return
		}
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: Invalid group name '%s'", req.Node)))
		return
	case queryALlPAByG:
		if req.IsExpectedLenFormat(1) {
			// 根据g查询所有的g.p.a
			pathP := fmt.Sprintf("/%d", dbg.Id)
			whereStr := "level=? and path=?"
			sps, err := StreePathGetMany(whereStr, department, pathP)
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}
			if len(sps) == 0 {
				log.Println("StreePathQuery: no p under g")
				return
			}

			// 根据p的结果拼接pathA,找出所有匹配的a
			var res []string
			for _, p := range sps {
				pathA := fmt.Sprintf("%s/%d", p.Path, p.Id)
				// 根据逐条pathA查处对应的p
				asp, err := StreePathGetMany(whereStr, application, pathA)
				if err != nil {
					log.Printf("%+v\n", err)
					continue
				}
				if len(asp) == 0 {
					continue
				}
				// 根据当前p结果拼接fullPath
				for _, a := range asp {
					fullPath := fmt.Sprintf("%s.%s.%s", dbg.NodeName, p.NodeName, a.NodeName)
					res = append(res, fullPath)
				}
			}
			sort.Strings(res)
			log.Println(res)
			return
		}
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: Invalid group name '%s'", req.Node)))
		return
	case queryAllAByPG:
		// 由于switch前的逻辑,到这里可以确认g存在
		if req.IsExpectedLenFormat(2) {
			// 通过p的level, path, node_name来找到具体的p
			p := strings.Split(req.Node, ".")[1]
			pathP := fmt.Sprintf("/%d", dbg.Id)
			whereStr := "level=? and path=? and node_name=?"
			dbp, err := StreePathGetOne(whereStr, department, pathP, p)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			if dbp == nil {
				log.Printf("%+v", errors.New(
					fmt.Sprintf(
						"StreePathQuery: p does not exist '%s'", p)))
				return
			}
			pathA := fmt.Sprintf("%s/%d", dbp.Path, dbp.Id)
			whereStr = "level = ? and path = ?"
			as, err := StreePathGetMany(whereStr, application, pathA)
			fmt.Println(as)
			return
		}
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: Invalid group.department name '%s'", req.Node)))
		return

	default:
		log.Printf(
			"%+v", errors.New(
				fmt.Sprintf(
					"StreePathQuery: target query type not supported '%d'", req.QueryType)))
	}
}

// StreePathAddOne 新增Path
func StreePathAddOne(req *common.NodeCommonReq) {
	// 要求新增对象必须是 g.p.a 三段式
	res := strings.Split(req.Node, ".")
	if len(res) != 3 {
		log.Printf("StreePathAddOne: Invalid path format: %s", req.Node)
		return
	}
	g, p, a := res[0], res[1], res[2]

	// 查询g
	nodeG := &StreePath{
		Level:    group,
		Path:     "0",
		NodeName: g,
	}
	dbG, err := nodeG.getOne()
	if err != nil {
		log.Printf("StreePathAddOne: Failed to get node g '%+v'", err)
		return
	}
	// 根据g查询结果再判断
	switch dbG {
	case nil:
		log.Println("StreePathAddOne: g is not exists, creating..")
		// 说明g不存在,依次插入g.p.a
		_, err := nodeG.addOne()
		if err != nil {
			log.Printf("StreePathAddOne: Adding g failed: '%+v'", err)
			return
		}
		log.Println("StreePathAddOne: g is created...")
		// 插入p并以g的id构造path
		log.Println("StreePathAddOne: p creating...")
		pathP := fmt.Sprintf("/%d", nodeG.Id)
		nodeP := &StreePath{
			Level:    department,
			Path:     pathP,
			NodeName: p,
		}
		_, err = nodeP.addOne()
		if err != nil {
			log.Printf("StreePathAddOne: Adding p failed: '%+v'", err)
			return
		}
		log.Println("StreePathAddOne: p is created...")
		// 插入a并以p,g id构造path
		log.Println("StreePathAddOne: a creating...")
		pathA := fmt.Sprintf("%s/%d", nodeP.Path, nodeP.Id)
		nodeA := &StreePath{
			Level:    application,
			Path:     pathA,
			NodeName: a,
		}
		_, err = nodeA.addOne()
		if err != nil {
			log.Printf("StreePathAddOne: Adding p failed: '%+v'", err)
			return
		}
		log.Println("StreePathAddOne: a is created...")
	default:
		// g存在查询p
		log.Println("StreePathAddOne: g exists")
		pathP := fmt.Sprintf("/%d", dbG.Id)
		nodeP := &StreePath{

			Level:    department,
			Path:     pathP,
			NodeName: p,
		}
		dbP, err := nodeP.getOne()
		if err != nil {
			log.Printf("StreePathAddOne: g exists but check p failed, path %s, err %+v", req.Node, err)
			return
		}
		if dbP != nil {
			log.Println("StreePathAddOne: g.p exists and going to find a")
			pathA := fmt.Sprintf("%s/%d", nodeP.Path, dbP.Id)
			nodeA := &StreePath{
				Level:    application,
				Path:     pathA,
				NodeName: a,
			}
			dbA, err := nodeA.getOne()
			if err != nil {
				log.Printf("StreePathAddOne: g.p exists but check a failed, path %s, err %+v", req.Node, err)
				return
			}
			if dbA != nil {
				log.Printf("StreePathAddOne: g.p.a exists, dont need to write again")
				return
			}
			_, err = nodeA.addOne()
			if err != nil {
				log.Printf("StreePathAddOne: g.p exists, writing a faied: %+v", err)
				return
			}
			log.Println("StreePathAddOne: adding a successfully")
			return
		}
		// 说明g存在,但p不存在,那么需要创建p的同时也创建a
		_, err = nodeP.addOne()
		if err != nil {
			log.Printf("StreePathAddOne: g exists, writing p faied: %+v", err)
			return
		}
		// 创建p成功,构造A并创建
		pathA := fmt.Sprintf("%s/%d", nodeP.Path, nodeP.Id)
		nodeA := &StreePath{
			Level:    application,
			Path:     pathA,
			NodeName: a,
		}
		_, err = nodeA.addOne()
		if err != nil {
			log.Printf("StreePathAddOne: g.p exists, writing a faied: %+v", err)
			return
		}
	}
	log.Println("StreePathAddOne: Adding g.p.a successfully")
}

// StreePathAddTest TESTING
func StreePathAddTest() {
	ns := []string{
		"inf.monitor.thanos",
		"inf.monitor.kafka",
		"inf.cicd.deploy",
		"inf.cicd.jenkins",
		"waimai.qiangdan.queue",
		"waimai.monitor.m3db",
		"waimai.ditu.kafka",
		"waimai.ditu.elasticsearch",
	}
	for _, n := range ns {
		req := &common.NodeCommonReq{
			Node: n,
		}
		StreePathAddOne(req)
	}
}

// StreePathQueryTest TESTING
func StreePathQueryTest() {
	ns := []string{
		"waimai.monitor",
		"waimai.ditu",
	}
	for _, n := range ns {
		req := &common.NodeCommonReq{
			Node:      n,
			QueryType: 3,
		}
		StreePathQuery(req)
	}
}
