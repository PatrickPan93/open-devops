package models

import (
	"fmt"
	"log"
	"open-devops/src/common"
	"strings"

	"github.com/pkg/errors"
)

func ResourceQuery(resourceType string, matchIds []uint64, limit, offset int) (interface{}, error) {
	ids := ""
	for _, id := range matchIds {
		ids += fmt.Sprintf("%d,", id)
	}
	ids = strings.TrimRight(ids, ",")

	whereInSql := fmt.Sprintf("id in (%s)", ids)
	log.Printf("models.ResourceQuery: generated ids %s and whereInSql %s", ids, whereInSql)

	var (
		res interface{}
		err error
	)
	switch resourceType {
	case common.ResourceHost:
		res, err = ResourceHostGetManyWithLimit(limit, offset, whereInSql)
	default:
		log.Printf("models.ResourceQuery: resourceType %s is not supported\n", resourceType)
	}
	return res, errors.Wrap(err, "")
}
