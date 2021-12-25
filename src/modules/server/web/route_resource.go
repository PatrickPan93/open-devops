package web

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"open-devops/src/common"
	"open-devops/src/models"
	mem_index "open-devops/src/modules/server/mem-index"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ResourceMount(c *gin.Context) {
	var inputs common.ResourceMountReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, http.StatusBadRequest, err)
		return
	}
	// Validate the resource name
	ok := models.CheckResources(inputs.ResourceType)
	if !ok {
		common.JSONR(c, http.StatusBadRequest, fmt.Sprintf("resource_node_not_exists: %v", inputs.ResourceType))
		return
	}
	// Check g.p.a if it's exist
	qReq := &common.NodeCommonReq{
		Node:      inputs.TargetPath,
		QueryType: 4,
	}
	gpa := models.StreePathQuery(qReq)
	if len(gpa) == 0 {
		common.JSONR(c, http.StatusBadRequest, fmt.Sprintf("target_type_no_exists: %v", inputs.TargetPath))
		return
	}
	rowsAff, err := models.ResourceMount(&inputs)
	if err != nil {
		common.JSONR(c, http.StatusInternalServerError, err)
		return
	}

	common.JSONR(c, http.StatusOK, fmt.Sprintf("resource_mount_add_success: rowsAff: %d", rowsAff))
}

func ResourceUnmount(c *gin.Context) {
	var inputs common.ResourceMountReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, http.StatusBadRequest, err)
		return
	}
	// Validate the resource name
	ok := models.CheckResources(inputs.ResourceType)
	if !ok {
		common.JSONR(c, http.StatusBadRequest, fmt.Sprintf("resource_node_not_exists: %v", inputs.ResourceType))
		return
	}
	// Check g.p.a if it's exist
	qReq := &common.NodeCommonReq{
		Node:      inputs.TargetPath,
		QueryType: 4,
	}
	gpa := models.StreePathQuery(qReq)
	if len(gpa) == 0 {
		common.JSONR(c, http.StatusBadRequest, fmt.Sprintf("target_type_no_exists: %v", inputs.TargetPath))
		return
	}
	rowsAff, err := models.ResourceUnmount(&inputs)
	if err != nil {
		common.JSONR(c, http.StatusInternalServerError, err)
		return
	}

	common.JSONR(c, http.StatusOK, fmt.Sprintf("resource_unmount_success: rowsAff: %d", rowsAff))
}

func ResourceQuery(c *gin.Context) {
	var inputs common.ResourceQueryReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, http.StatusBadRequest, err)
		return
	}
	// Check the resourceType index has registered to index container
	ok := mem_index.JudgeResourceIndexExists(inputs.ResourceType)
	if !ok {
		common.JSONR(c, http.StatusBadRequest, fmt.Sprintf("ResourceType_index_not_exists: %v", inputs.ResourceType))
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	if err != nil {
		common.JSONR(c, http.StatusBadRequest, fmt.Sprintf("invalid_page_size: %v", err))
		return
	}
	currentPage, err := strconv.Atoi(c.DefaultQuery("current_size", "1"))
	if err != nil {
		common.JSONR(c, http.StatusBadRequest, fmt.Sprintf("invalid_current_page_size: %v", err))
		return
	}
	offset := 0
	limit := 0
	limit = pageSize

	if currentPage > 1 {
		offset = (currentPage - 1) * limit
	}
	matchIds := mem_index.GetMatchIdsByIndex(inputs)
	// hardcode matchIds for testing
	//matchIds := []uint64{384, 385, 386, 387, 388, 389}
	totalCount := len(matchIds)

	// 计算应该分页的数量（取整）
	pageCount := int(math.Ceil(float64(totalCount) / float64(limit)))

	// 构造响应对象
	resp := common.QueryResponse{
		Code:        http.StatusOK,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageCount:   pageCount,
		TotalCount:  totalCount,
		Result:      nil,
	}

	if len(matchIds) == 0 {
		common.JSONR(c, http.StatusBadRequest, "failed to get matchIds by index")
		log.Println("route_resource.ResourceQuery: failed to get matchIds by index")
		return
	}
	res, err := models.ResourceQuery(inputs.ResourceType, matchIds, limit, offset)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Result = err
		common.JSONR(c, http.StatusInternalServerError, err)
		return
	}
	resp.Result = res
	common.JSONR(c, resp)
}
