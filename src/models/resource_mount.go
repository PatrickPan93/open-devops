package models

import (
	"fmt"
	"log"
	"open-devops/src/common"
	"strings"

	"github.com/pkg/errors"
)

var availableResources = map[string]struct{}{
	"resource_host": {},
}

func CheckResources(resource string) bool {
	_, ok := availableResources[resource]
	return ok
}

func ResourceMount(req *common.ResourceMountReq) (int64, error) {
	gpas := strings.Split(req.TargetPath, ".")
	g, p, a := gpas[0], gpas[1], gpas[2]
	ids := ""
	for _, id := range req.ResourceIds {
		ids += fmt.Sprintf("%d,", id)
	}
	ids = strings.TrimRight(ids, ",")
	rawSql := fmt.Sprintf(`update %s set stree_group="%s", stree_product="%s", stree_app="%s" where id in (%s)`, req.ResourceType, g, p, a, ids)
	log.Printf("models.ResourceMount: raw sql is %s\n", rawSql)
	res, err := DB["stree"].Exec(rawSql)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("models.ResourceMount: error while executing raw sql: %s", rawSql))
	}
	return res.RowsAffected()
}
