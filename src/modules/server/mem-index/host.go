package mem_index

import (
	"encoding/json"
	"fmt"
	"log"
	"open-devops/src/models"
	"strconv"

	"github.com/ning1875/inverted-index/labels"

	"github.com/pkg/errors"

	ii "github.com/ning1875/inverted-index"
)

type HostIndex struct {
	Ir      *ii.HeadIndexReader
	Logger  log.Logger
	Modulus int // 分片总数
	Num     int // 分片第几个
}

func (hi *HostIndex) FlushIndex() {
	var limit int64
	var offset int64
	r := new(models.ResourceHost)
	// 计算resource_host表中数据量
	count, err := r.Count()
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, "mem-index.FlushIndex: Error while counting num of resource_host"))
		return
	}
	for i := 0; i < int(count); i++ {
		// 如果不存在分片（来源于配置)
		if hi.Modulus == 0 {
			limit = count
			offset = 0
			break
		}
		/*
			//TODO 解决分片取数据问题
				// 如果对i取模,与当前节点分片编号相同,则keep住该id
				if i%hi.Modulus == hi.Num {
					ids = append(ids, i)
					continue
				}

		*/
	}

	whereInSql := "id > 0"
	//objs, err := models.ResourceHostGetManyWithLimit(ids[len(ids)-1], , whereInSql)
	objs, err := models.ResourceHostGetManyWithLimit(int(limit), int(offset), whereInSql)

	if err != nil {
		log.Printf("%+v", errors.Wrapf(err, "mem-index.FlushIndex: Error while getting data from resource_host by ids %d\n", limit))
		return
	}
	// 自动刷node path(g.p.a)
	thisGPAS := map[string]struct{}{}

	thisH := ii.NewHeadReader()
	for _, item := range objs {
		// 取出tags字段并json化
		m := make(map[string]string)
		m["hash"] = item.Hash

		tags := make(map[string]string)
		// 数组型数据 内网ips 公网ips 安全组
		prIps := []string{}
		pbIps := []string{}
		//secGs := []string{}

		m["uid"] = item.Uid
		m["name"] = item.Name
		m["cloud_provider"] = item.CloudProvider
		m["charging_mode"] = item.ChargingMode
		m["region"] = item.Region
		m["instance_type"] = item.InstanceType
		m["availability_zone"] = item.AvailabilityZone
		m["vpc_id"] = item.VpcId
		m["subnet_id"] = item.SubnetId
		m["status"] = item.Status

		m["account_id"] = strconv.FormatInt(item.AccountId, 10)

		// json列表型
		json.Unmarshal(item.PrivateIps, &prIps)
		json.Unmarshal(item.PublicIps, &pbIps)

		// json map型
		json.Unmarshal(item.Tags, &tags)
		// cpu mem
		m["cpu"] = item.CPU
		m["mem"] = item.Mem
		m["disk"] = item.Disk

		// g.p.a
		m["stree_group"] = item.StreeGroup
		m["stree_product"] = item.StreeProduct
		m["stree_app"] = item.StreeApp
		thisGPAS[fmt.Sprintf("%s.%s.%s", item.StreeGroup, item.StreeProduct, item.StreeApp)] = struct{}{}

		// 调用倒排索引库刷新索引
		thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapToLabelsSet(m))
		thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapToLabelsSet(tags))

		// 数组型
		for _, i := range prIps {
			mp := map[string]string{
				"private_ips": i,
			}
			thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapToLabelsSet(mp))
		}
		for _, i := range pbIps {
			mp := map[string]string{
				"public_ips": i,
			}
			thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapToLabelsSet(mp))
		}
	}
	hi.Ir.Reset(thisH)
}

//map转换为labels
func mapToLabelsSet(m map[string]string) labels.Labels {
	var lSet labels.Labels
	for k, v := range m {
		l := labels.Label{
			Name:  k,
			Value: v,
		}
		lSet = append(lSet, l)
	}
	return lSet
}

func (hi *HostIndex) GetIndexReader() *ii.HeadIndexReader {
	return hi.Ir
}
