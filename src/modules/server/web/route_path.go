package web

import (
	"fmt"
	"net/http"
	"open-devops/src/common"
	"open-devops/src/models"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func NodePathAdd(c *gin.Context) {
	var inputs common.NodeCommonReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, http.StatusBadRequest)
		return
	}
	res := strings.Split(inputs.Node, ".")
	if len(res) != 3 {
		common.JSONR(c, http.StatusBadRequest, errors.Errorf("path_invalidate: %s", inputs.Node))
		return
	}
	err := models.StreePathAddOne(&inputs)
	if err != nil {
		common.JSONR(c, http.StatusInternalServerError, err)
		return
	}
	common.JSONR(c, http.StatusOK, "path_add_success")
}

func NodePathQuery(c *gin.Context) {
	var inputs common.NodeCommonReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, http.StatusBadRequest)
		return
	}
	if inputs.QueryType == 3 && len(strings.Split(inputs.Node, ".")) != 2 {
		common.JSONR(c, http.StatusBadRequest, errors.Errorf("query_type=3 path shoud be a.b: %v", inputs.Node))
		return
	}

	res := models.StreePathQuery(&inputs)
	common.JSONR(c, http.StatusOK, res)
}

func NodePathDelete(c *gin.Context) {
	var inputs common.NodeCommonReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, http.StatusBadRequest)
		return
	}
	delNum := models.StreePathDelete(&inputs)
	if delNum == 0 {
		common.JSONR(c, http.StatusOK, "nothing_deleted")
		return
	}
	common.JSONR(c, http.StatusOK, fmt.Sprintf("node_path_delete: %d has been deleted", delNum))
}
