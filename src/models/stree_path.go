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
	queryGPAByGPA = 4

	deleteGIfNoPUnderG = 1
	deletePIfNoAUnderP = 2
	deleteAifGPAExist  = 3

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

func (sp *StreePath) AddOne() (int64, error) {
	return DB[streeDB].InsertOne(sp)
}

// GetOne query one record base on query param
func (sp *StreePath) GetOne() (*StreePath, error) {
	has, err := DB[streeDB].Get(sp)
	if err != nil {
		return nil, errors.Wrap(err, "GetOne: failed to get one streePath record")
	}
	if !has {
		return nil, nil
	}
	return sp, nil
}

// DelOne  delete one record by obj
func (sp *StreePath) DelOne() (int64, error) {
	return DB["stree"].Delete(sp)
}

// CheckExist Check streePath if it's exists
func (sp *StreePath) CheckExist() (bool, error) {
	return DB[streeDB].Exist(sp)
}

// StreePathGetOne  函数的get
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

func StreePathQuery(req *common.NodeCommonReq) []string {
	// whatever you want to query. we must make sure that g exists.
	var dbg *StreePath
	// split by .  get g name from index 0
	gName := strings.Split(req.Node, ".")[0]
	nodeG := StreePath{
		Level:    group,
		Path:     originPath,
		NodeName: gName,
	}
	dbg, err := nodeG.GetOne()
	if err != nil {
		log.Printf("%+v\n", err)
		return nil
	}
	if dbg == nil {
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: g does not exist '%s'", req.Node)))
		return nil
	}
	var res []string
	switch req.QueryType {
	case queryAllPByG:
		if req.IsExpectedLenFormat(1) {
			pathP := fmt.Sprintf("/%d", dbg.Id)
			whereStr := "level=? and path=?"
			sps, err := StreePathGetMany(whereStr, department, pathP)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil
			}
			if len(sps) == 0 {
				log.Printf("%+v", errors.New(
					fmt.Sprintf(
						"StreePathQuery: no p under g '%s'", dbg.NodeName)))
				return nil
			}

			for _, i := range sps {
				res = append(res, i.NodeName)
			}
			//log.Println(res)
			return res
		}
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: Invalid group name '%s'", req.Node)))
		return nil
	case queryALlPAByG:
		if req.IsExpectedLenFormat(1) {
			// 根据g查询所有的g.p.a
			pathP := fmt.Sprintf("/%d", dbg.Id)
			whereStr := "level=? and path=?"
			sps, err := StreePathGetMany(whereStr, department, pathP)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil
			}
			if len(sps) == 0 {
				log.Println("StreePathQuery: no p under g")
				return nil
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
			//log.Println(res)
			return res
		}
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: Invalid group name '%s'", req.Node)))
		return nil
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
				return nil
			}
			if dbp == nil {
				log.Printf("%+v", errors.New(
					fmt.Sprintf(
						"StreePathQuery: p does not exist '%s'", p)))
				return nil
			}
			pathA := fmt.Sprintf("%s/%d", dbp.Path, dbp.Id)
			whereStr = "level = ? and path = ?"
			as, err := StreePathGetMany(whereStr, application, pathA)
			var res []string
			for _, a := range as {
				fullPath := fmt.Sprintf("%s.%s.%s", dbg.NodeName, dbp.NodeName, a.NodeName)
				res = append(res, fullPath)
			}
			sort.Strings(res)
			//fmt.Println(res)
			return res
		}
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: Invalid group.department name '%s'", req.Node)))
		return nil
	case queryGPAByGPA:
		if req.IsExpectedLenFormat(3) {
			p := strings.Split(req.Node, ".")[1]
			pathP := fmt.Sprintf("/%d", dbg.Id)
			whereStr := "level=? and path=? and node_name=?"
			dbp, err := StreePathGetOne(whereStr, department, pathP, p)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil
			}
			if dbp == nil {
				log.Printf("%+v", errors.New(
					fmt.Sprintf(
						"StreePathQuery: p does not exist '%s'", p)))
				return nil
			}
			pathA := fmt.Sprintf("%s/%d", dbp.Path, dbp.Id)
			whereStr = "level = ? and path = ? and node_name = ?"
			a := strings.Split(req.Node, ".")[2]

			fmt.Println(pathA, a)
			dba, err := StreePathGetOne(whereStr, application, pathA, a)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil
			}
			if dba == nil {
				log.Printf("%+v", errors.New(
					fmt.Sprintf(
						"StreePathQuery: a does not exist '%s'", a)))
				return nil
			}
			fullPath := fmt.Sprintf("%s.%s.%s", dbg.NodeName, dbp.NodeName, dba.NodeName)
			res = append(res, fullPath)
			return res
		}
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: Invalid group.department.appolication name '%s'", req.Node)))
		return nil
	default:
		log.Printf(
			"%+v", errors.New(
				fmt.Sprintf(
					"StreePathQuery: target query type not supported '%d'", req.QueryType)))
	}
	return nil
}

// StreePathAddOne 新增Path
func StreePathAddOne(req *common.NodeCommonReq) error {
	// 要求新增对象必须是 g.p.a 三段式
	res := strings.Split(req.Node, ".")
	if len(res) != 3 {
		log.Printf("StreePathAddOne: Invalid path format: %s", req.Node)
		return errors.Errorf("StreePathAddOne: Invalid path format: %s", req.Node)
	}
	g, p, a := res[0], res[1], res[2]

	// 查询g
	nodeG := &StreePath{
		Level:    group,
		Path:     "0",
		NodeName: g,
	}
	dbG, err := nodeG.GetOne()
	if err != nil {
		log.Printf("StreePathAddOne: Failed to get node g '%+v'", err)
		return errors.Wrap(err, "")
	}
	// 根据g查询结果再判断
	switch dbG {
	case nil:
		log.Println("StreePathAddOne: g is not exists, creating..")
		// 说明g不存在,依次插入g.p.a
		_, err := nodeG.AddOne()
		if err != nil {
			log.Printf("StreePathAddOne: Adding g failed: '%+v'", err)
			return errors.Wrap(err, "")
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
		_, err = nodeP.AddOne()
		if err != nil {
			log.Printf("StreePathAddOne: Adding p failed: '%+v'", err)
			return errors.Wrap(err, "")
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
		_, err = nodeA.AddOne()
		if err != nil {
			log.Printf("StreePathAddOne: Adding p failed: '%+v'", err)
			return errors.Wrap(err, "")
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
		dbP, err := nodeP.GetOne()
		if err != nil {
			log.Printf("StreePathAddOne: g exists but check p failed, path %s, err %+v", req.Node, err)
			return errors.Wrap(err, "")
		}
		if dbP != nil {
			log.Println("StreePathAddOne: g.p exists and going to find a")
			pathA := fmt.Sprintf("%s/%d", nodeP.Path, dbP.Id)
			nodeA := &StreePath{
				Level:    application,
				Path:     pathA,
				NodeName: a,
			}
			dbA, err := nodeA.GetOne()
			if err != nil {
				log.Printf("StreePathAddOne: g.p exists but check a failed, path %s, err %+v", req.Node, err)
				return errors.Wrap(err, "")
			}
			if dbA != nil {
				log.Printf("StreePathAddOne: g.p.a exists, dont need to write again")
				return errors.Wrap(err, "")
			}
			_, err = nodeA.AddOne()
			if err != nil {
				log.Printf("StreePathAddOne: g.p exists, writing a faied: %+v", err)
				return errors.Wrap(err, "")
			}
			log.Println("StreePathAddOne: adding a successfully")
			return errors.Wrap(err, "")
		}
		// 说明g存在,但p不存在,那么需要创建p的同时也创建a
		_, err = nodeP.AddOne()
		if err != nil {
			log.Printf("StreePathAddOne: g exists, writing p faied: %+v", err)
			return errors.Wrap(err, "")
		}
		// 创建p成功,构造A并创建
		pathA := fmt.Sprintf("%s/%d", nodeP.Path, nodeP.Id)
		nodeA := &StreePath{
			Level:    application,
			Path:     pathA,
			NodeName: a,
		}
		_, err = nodeA.AddOne()
		if err != nil {
			log.Printf("StreePathAddOne: g.p exists, writing a faied: %+v", err)
			return errors.Wrap(err, "")
		}
	}
	log.Println("StreePathAddOne: Adding g.p.a successfully")
	return nil
}

func StreePathDelete(req *common.NodeCommonReq) int64 {
	// whatever we want to delete, we must make sure g existed.
	var dbg *StreePath

	// split by .  get g name from index 0
	path := strings.Split(req.Node, ".")
	pLevel := len(path)
	nodeG := &StreePath{
		Level:    group,
		Path:     originPath,
		NodeName: path[0], // means to get g_name
	}
	dbg, err := nodeG.GetOne()
	if err != nil {
		log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while getting g %s", nodeG.NodeName)))
		return 0
	}
	if dbg == nil {
		log.Printf("%+v", errors.New(
			fmt.Sprintf(
				"StreePathQuery: g does not exist '%s'", req.Node)))
		return 0
	}
	if req.ForceDelete {
		whereStr := "path like ? or path= ?"
		sp, err := StreePathGetMany(whereStr, fmt.Sprintf(`/%d/%%%%`, dbg.Id), fmt.Sprintf("/%d", dbg.Id))
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while getting all child nodes of g %s", dbg.NodeName)))
			return 0
		}
		for _, v := range sp {
			delNum, err := v.DelOne()
			if err != nil {
				log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: error while deleting data via force delete mode %v", v)))
				return 0
			}
			if delNum < 1 {
				log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: error while deleting data via force delete mode %v", v)))
				return delNum
			}
			log.Println(fmt.Sprintf("StreePathDelete: deleting data via force mode successfully %v", v))
		}
		delNum, err := dbg.DelOne()
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: error while deleting data via force delete mode %v", dbg.NodeName)))
			return 0
		}
		if delNum < 1 {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: error while deleting data via force delete mode %v", dbg.NodeName)))
			return delNum
		}
		log.Printf(fmt.Sprintf("StreePathDelete: deleting g %s with force mode successfully", dbg.NodeName))
		return delNum
	}
	switch pLevel {
	case deleteGIfNoPUnderG:
		// g existed. trying to find p
		pathP := fmt.Sprintf("/%d", dbg.Id)
		whereStr := "level=? and path=?"
		ps, err := StreePathGetMany(whereStr, department, pathP)
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: failed to get ps %s", pathP)))
			return 0
		}
		if len(ps) > 0 {
			log.Printf(
				fmt.Sprintf(
					"StreePathDelete: failed to delete g %s , there are p %v under g %s", dbg.NodeName, ps, dbg.NodeName))
			return 0
		}
		delNum, err := dbg.DelOne()
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while deleting g %s", dbg.NodeName)))
			return 0
		}
		log.Println(fmt.Sprintf("StreePathDelete: Deleting g %s successfully.. %d g get deleted", dbg.NodeName, delNum))
		return delNum
	case deletePIfNoAUnderP:
		// g existed. trying to find p
		pathP := fmt.Sprintf("/%d", dbg.Id)
		nodeP := &StreePath{
			Level:    department,
			Path:     pathP,
			NodeName: path[1],
		}
		dbp, err := nodeP.GetOne()
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while deleting p %s", nodeP.NodeName)))
			return 0
		}
		if dbp == nil {
			log.Printf("%+v", errors.New(
				fmt.Sprintf(
					"StreePathQuery: p does not exist '%s'", req.Node)))
			return 0
		}
		pathA := fmt.Sprintf("%s/%d", dbp.Path, dbp.Id)
		whereStr := "level=? and path=?"
		as, err := StreePathGetMany(whereStr, application, pathA)
		if err != nil {
			log.Printf("%+v\n", fmt.Sprintf("StreePathQuery: error while getting a %s", pathA))
			return 0
		}
		if len(as) > 0 {
			log.Printf(
				fmt.Sprintf(
					"StreePathDelete: failed to delete p %s.%s , there are a %v under p %s.%s",
					dbg.NodeName, dbp.NodeName, as, dbg.NodeName, dbp.NodeName))
			return 0
		}
		delNum, err := dbp.DelOne()
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while deleting p %s", dbp.NodeName)))
			return 0
		}
		log.Println(fmt.Sprintf("StreePathDelete: Deleting p %s.%s successfully.. %s.%s p get deleted",
			dbg.NodeName, dbp.NodeName, dbg.NodeName, dbp.NodeName))
		return delNum
	case deleteAifGPAExist:
		pathP := fmt.Sprintf("/%d", dbg.Id)
		nodeP := &StreePath{
			Level:    department,
			Path:     pathP,
			NodeName: path[1],
		}
		dbp, err := nodeP.GetOne()
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while getting p %s.%s.%s", dbg.NodeName, nodeP.NodeName, path[2])))
			return 0
		}
		if dbp == nil {
			log.Printf("%+v", errors.New(
				fmt.Sprintf(
					"StreePathQuery: p does not exist '%s'", req.Node)))
			return 0
		}
		pathA := fmt.Sprintf("%s/%d", dbp.Path, dbp.Id)
		whereStr := "level=? and path=? and node_name=?"
		dba, err := StreePathGetOne(whereStr, application, pathA, path[2])
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while getting a %s.%s.%s", dbg.NodeName, dbp.NodeName, path[2])))
			return 0
		}
		if dba == nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete:  a does not exist %s.%s.%s", dbg.NodeName, dbp.NodeName, path[2])))
			return 0
		}
		delNum, err := dba.DelOne()
		if err != nil {
			log.Printf("%+v", errors.New(fmt.Sprintf("StreePathDelete: Error while deleting a %s.%s.%s", dbg.NodeName, dbp.NodeName, path[2])))
			return 0
		}
		log.Println(fmt.Sprintf("StreePathDelete: Deleting p %s.%s.%s successfully.. %s.%s.%s p get deleted",
			dbg.NodeName, dbp.NodeName, dba.NodeName, dbg.NodeName, dbp.NodeName, dba.NodeName))
		return delNum

	default:
		log.Println("StreePathDelete: default logic")
		return 0
	}
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
		"inf.building",
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
		"inf.cicd"}
	for _, n := range ns {
		req := &common.NodeCommonReq{
			Node:      n,
			QueryType: 3,
		}
		StreePathQuery(req)
	}
}

func StreePathDeleteTest() {
	ns := []string{
		"waimai",
	}
	for _, n := range ns {
		req := &common.NodeCommonReq{Node: n, ForceDelete: true}
		res := StreePathDelete(req)
		fmt.Println(res)
	}

}

// TODO: A function to support RAW sql execution.
