package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	//Tags       json.RawMessage `json:"tags"`
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
	SN         string    `json:"sn" xorm:"-"`
	CPU        string    `json:"cpu" xorm:"cpu""`
	Mem        string    `json:"mem"`
	Disk       string    `json:"disk"`
	IPAddr     string    `json:"ip_addr" xorm:"-"`
	Hostname   string    `json:"hostname" xorm:"-"`
	CreateTime time.Time `json:"create_time" xorm:"create_time created"`
	UpdateTime time.Time `json:"update_time" xorm:"update_time updated"`
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

func (rh *ResourceHost) Update() (int64, error) {
	return DB["stree"].Update(rh)
}
