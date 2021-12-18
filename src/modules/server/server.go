package main

import (
	"log"
	"open-devops/src/models"
	"open-devops/src/modules/server/config"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/common/promlog"
	promlogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// 命令行解析
	app = kingpin.New(filepath.Base(os.Args[0]), "The open-devops server")
	// 指定配置文件参数
	configFile = app.Flag("config.file", "open-devops-server configuration file").Short('c').Default("open-devops-server.yaml").String()
)

func main() {
	// Help Info
	app.HelpFlag.Short('h')
	// 基于prometheus库通过build ldflags版本信息注入
	promlogConfig := promlog.Config{}
	app.Version(version.Print("open-devops-server"))
	promlogflag.AddFlags(app, &promlogConfig)

	// Get Param From Command line from $1 to the end
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.Println("Start loading config...")
	serverConfig, err := config.LoadFile(*configFile)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	log.Println("Loading config successfully")

	// Init MySQL
	models.InitMySQL(serverConfig.MysqlS)

}
