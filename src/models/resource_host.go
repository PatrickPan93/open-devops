package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// AgentCollectInfo 机器上采集info的字段
type AgentCollectInfo struct {
	SN       string `json:"sn"`
	CPU      string `json:"cpu"`
	Mem      string `json:"mem"`
	Disk     string `json:"disk"`
	IPAddr   string `json:"ip_addr"`
	Hostname string `json:"hostname"`
}
type ResourceHost struct {
	// 公共字段
	Id         int64           `json:"id"`
	Uid        string          `json:"uid"`
	Hash       string          `json:"hash"`
	Name       string          `json:"name"`
	PrivateIps json.RawMessage `json:"private_ips"`
	Tags       json.RawMessage `json:"tags"`
	// 公有云字段
	CloudProvider    string          `json:"cloud_provider"`
	ChargingMode     string          `json:"charging_mode"`
	Region           string          `json:"region"`
	AccountId        int64           `json:"account_id"`
	VpcId            string          `json:"vpc_id"`
	SubnetId         string          `json:"subnet_id"`
	SecurityGroups   string          `json:"security_group"`
	Status           string          `json:"status"`
	InstanceType     string          `json:"instance_type"`
	PublicIps        json.RawMessage `json:"public_ips"`
	AvailabilityZone string          `json:"availability_zone"`

	// 机器采集字段
	SN           string    `json:"sn" xorm:"-"`
	CPU          string    `json:"cpu" xorm:"cpu""`
	Mem          string    `json:"mem"`
	Disk         string    `json:"disk"`
	IPAddr       string    `json:"ip_addr" xorm:"-"`
	Hostname     string    `json:"hostname" xorm:"-"`
	StreeGroup   string    `json:"stree_group"`
	StreeProduct string    `json:"stree_product"`
	StreeApp     string    `json:"stree_app"`
	CreateTime   time.Time `json:"create_time" xorm:"create_time created"`
	UpdateTime   time.Time `json:"update_time" xorm:"update_time updated"`
}

// GetHash 判断资源是否发生变化函数
func (rh *ResourceHost) GetHash() string {
	h := md5.New()
	h.Write([]byte(rh.SN))
	h.Write([]byte(rh.Name))
	h.Write([]byte(rh.IPAddr))
	h.Write([]byte(rh.CPU))
	h.Write([]byte(rh.Mem))
	h.Write([]byte(rh.Disk))
	return hex.EncodeToString(h.Sum(nil))
}

func (rh *ResourceHost) GetOne() (*ResourceHost, error) {
	has, err := DB["stree"].Get(rh)
	if err != nil {
		return nil, errors.Wrap(err, "models.GetOne: Error while getting ResrouceHost")
	}
	if !has {
		return nil, nil
	}
	return rh, nil
}

func (rh *ResourceHost) AddOne() (int64, error) {
	return DB["stree"].InsertOne(rh)
}

func (rh *ResourceHost) UpdateByUid(uid string) (int64, error) {
	return DB["stree"].Where("uid=?", uid).Update(rh)
}

func (rh *ResourceHost) Count() (int64, error) {
	return DB["stree"].Where("id>0").Count(rh)
}

func GetHostUidAndHash() (map[string]string, error) {
	var objs []ResourceHost
	err := DB["stree"].Cols("uid", "hash").Find(&objs)
	if err != nil {
		return nil, errors.Wrap(err, "models.GetHostUidAndHash: error while getting host uid and hash")
	}
	m := make(map[string]string)
	for _, h := range objs {
		m[h.Uid] = h.Hash
	}
	return m, nil
}

func BatchDeleteResource(tableName string, idKey string, ids []string) (int64, error) {

	var whereInStr string
	for _, v := range ids {
		whereInStr += fmt.Sprintf("\"%s\",", v)
	}
	whereInStr = strings.TrimRight(whereInStr, ",")
	rawSql := fmt.Sprintf(`delete from %s where %s in (%s)`,
		tableName,
		idKey,
		whereInStr)
	res, err := DB["stree"].Exec(rawSql)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("models.BatchDeleteResource: error while deleting resource_hosts: %s", ids))
	}
	return res.RowsAffected()
}

func ResourceHostGetManyWithLimit(limit, offset int, whereStr string, args ...interface{}) ([]ResourceHost, error) {
	var objs []ResourceHost
	err := DB["stree"].Where(whereStr, args...).Limit(limit, offset).Find(&objs)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("models.ResourceHostGetManyWithLimit: error while getting resource_hosts by Ids: %s", whereStr))
	}
	return objs, nil

}

func ResourceHostGetMany(whereStr string, args ...interface{}) ([]ResourceHost, error) {
	var objs []ResourceHost
	err := DB["stree"].Where(whereStr, args...).Find(&objs)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("models.ResourceHostGetMany: error while getting resource_hosts by Ids: %s", whereStr))
	}
	return objs, nil

}

/*
type nodeMap map[string][]int64

type invertedIndex struct {
	Lock     sync.Mutex
	IndexMap map[string]nodeMap
}

func InvertIndexTest() {
	var jsonMap map[string]string
	var rId int64

	var ii = &invertedIndex{
		Lock:     sync.Mutex{},
		IndexMap: make(map[string]nodeMap, 0),
	}

	// Getting all result from db
	rhs := make([]ResourceHost, 0)
	err := DB["stree"].Cols("id", "tags").Find(&rhs)
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, ""))
	}

	// Get id and Tags value from result
	ii.Lock.Lock()
	for _, r := range rhs {
		rId = r.Id
		jsonBytes, _ := r.Tags.MarshalJSON()
		json.Unmarshal(jsonBytes, &jsonMap)
		fmt.Println(rId, jsonMap)

		for k, v := range jsonMap {
			//fmt.Println(k, v)
			// 如果key不存在,则直接创建key以及node map
			if result := ii.IndexMap[k]; result == nil {
				ii.IndexMap[k] = nodeMap{v: {rId}}
				continue
			}
			if nodeResult := ii.IndexMap[k][v]; nodeResult == nil {
				ii.IndexMap[k][v] = []int64{rId}
				//ii.IndexMap[k][v] = nodeMap{v: {rId}}
				continue
			}
			iSlice := ii.IndexMap[k][v]
			iSlice = append(iSlice, rId)
			ii.IndexMap[k][v] = iSlice
		}
	}
	ii.Lock.Unlock()
	fmt.Println(ii.IndexMap)
	bytes, _ := json.Marshal(ii.IndexMap)
	fmt.Println(string(bytes))

	if res := ii.IndexMap["test"]["test"]; res == nil {
		fmt.Println("结果为空 触发其它logic")
	}
	if res := ii.IndexMap["arch"]["amd64"]; res != nil {
		var rhs []ResourceHost
		var strs []string
		for _, v := range res {
			strRes := strconv.FormatInt(v, 10)
			strs = append(strs, strRes)
		}
		whereInStr := strings.Join(strs, ",")

		whereStr := fmt.Sprintf("id in (%s)", whereInStr)
		DB["stree"].Where(whereStr).Find(&rhs)
		fmt.Println(rhs)
	}
}

*/
