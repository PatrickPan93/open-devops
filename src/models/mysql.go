package models

import (
	"log"
	"open-devops/src/modules/server/config"
	"time"

	"github.com/pkg/errors"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

var DB = map[string]*xorm.Engine{}

func InitMySQL(mysqlS []*config.MysqlConf) {
	for _, conf := range mysqlS {
		db, err := xorm.NewEngine("mysql", conf.Addr)
		if err != nil {
			log.Fatalf("[init.mysql.erropr][cannot connect to mysql][addr:%v][error:%+v]\n", conf.Addr, errors.WithStack(err))
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("[init.mysql.erropr][cannot connect to mysql][addr:%v][error:%+v]\n", conf.Addr, errors.WithStack(err))

		}
		db.SetMaxIdleConns(conf.MaxIdle)
		db.SetMaxOpenConns(conf.MaxCon)
		db.SetConnMaxLifetime(time.Hour)
		db.ShowSQL(conf.Debug)
		db.Logger().SetLevel(xlog.LOG_INFO)
		DB[conf.Name] = db
	}
}
