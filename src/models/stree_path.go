package models

import (
	"fmt"

	"github.com/pkg/errors"
)

// StreePath table stree_path
type StreePath struct {
	Id       int64  `json:"id"`
	Level    int64  `json:"level"`
	Path     string `json:"path"`
	NodeName string `json:"node_name"`
}

func (sp *StreePath) addOne() (int64, error) {
	rowAffect, err := DB["stree"].InsertOne(sp)
	return rowAffect, err
}

func StreePathGet(where string, args ...interface{}) (*StreePath, error) {
	var obj StreePath
	isExists, err := DB["stree"].Where(where, args...).Get(&obj)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("StreePathGet: failed to get path using where '%s,%v'\n", where, args))
	}
	if !isExists {
		return nil, nil
	}
	return &obj, nil

}
