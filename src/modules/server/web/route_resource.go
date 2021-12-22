package web

import (
	"fmt"
	"net/http"
	"open-devops/src/common"
	"open-devops/src/models"

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
