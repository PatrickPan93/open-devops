package common

type ResourceMountReq struct {
	// required 必填参数 type: host rds dcs
	ResourceType string  `json:"resource_type" binding:"required"`
	ResourceIds  []int64 `json:"resource_ids" binding:"required"`
	TargetPath   string  `json:"target_path" binding:"required"`
}
