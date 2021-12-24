package main

import (
	"context"
	"log"
	"open-devops/src/models"
	"open-devops/src/modules/server/cloud_sync"
	"open-devops/src/modules/server/config"
	mem_index "open-devops/src/modules/server/mem-index"
	"open-devops/src/modules/server/rpc"
	"open-devops/src/modules/server/web"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"

	"github.com/oklog/run"

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

	// Init Inverted-index
	mem_index.Init(serverConfig.IndexModulesConf)

	// testing function
	models.StreePathAddTest()
	//models.StreePathQueryTest()
	//models.StreePathDeleteTest()
	//models.InvertIndexTest()

	// design the running group(multi go routine)
	var g run.Group
	ctx, cancel := context.WithCancel(context.Background())

	{
		// 接收signal的chan
		signalChan := make(chan os.Signal, 1)
		// 接收cancel信息的chan
		cancelChan := make(chan struct{})
		// 监听来自系统的terminal相关信号
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		// the first: execution func
		// the second: the error handling func
		// 通过第一个g.add 来控制整体go goroutine的生命周期.
		// 假设第一组g.add的os.notify嗅探到了term信号,便会执行cancel()
		// 由于其它g.add都有listen <-ctx.Done()，所以当主routine嗅探到term信号后执行cancel()，其它go routine都会开始退出
		//
		g.Add(func() error {
			select {
			case <-signalChan:
				log.Println("notify a SIGTERM syscall.. process will exit soon")
				// cancel() if signalChan got an term signal
				cancel()
				return nil
			case <-cancelChan:
				log.Println("Received a cancel event")
				return nil
			}
		},
			func(error) {
				close(cancelChan)
			})
	}
	{
		// starting rpc server here
		g.Add(func() error {
			errChan := make(chan error, 1)
			// 由于rpc server特殊性,无法传入ctx,所以需要用go routine去进行启动
			go func() {
				errChan <- rpc.Start(serverConfig.RPCAddr)
			}()

			select {
			// 如果rpc server阻塞解决并且抛出异常，则退出当前go routine并抛出错误
			case err := <-errChan:
				log.Printf("%+v", errors.Wrap(err, "rpc server running error"))
				return err
			// 如果ctx被显式cancel,退出当前go routine

			case <-ctx.Done():
				log.Println("rpc server receive quit signal.. would be stopped soon")
				return nil

			}
		},
			func(err error) {
				cancel()
			})
	}
	{
		// starting rpc server here
		g.Add(func() error {
			errChan := make(chan error, 1)
			// 由于rpc server特殊性,无法传入ctx,所以需要用go routine去进行启动
			go func() {
				errChan <- web.StartGin(serverConfig.HTTPAddr)
			}()

			select {
			// 如果rpc server阻塞解决并且抛出异常，则退出当前go routine并抛出错误
			case err := <-errChan:
				log.Printf("%+v", errors.Wrap(err, "web server running error"))
				return err
			// 如果ctx被显式cancel,退出当前go routine
			case <-ctx.Done():
				log.Println("web server receive quit signal.. would be stopped soon")
				return nil

			}
		},
			func(err error) {
				cancel()
			})
	}
	{
		//公有云同步
		if serverConfig.PublicCloudSyncConf.Enable {
			g.Add(func() error {
				return cloud_sync.CloudSyncManager(ctx)
			}, func(err error) {
				log.Printf("%+v", errors.Wrap(err, "cloud_sync.CloudSyncManager running error"))
				cancel()
			})
		}

	}
	{
		// 倒排索引更新

		g.Add(func() error {
			return mem_index.RevertedIndexSyncManager(ctx)
		}, func(err error) {
			log.Printf("%+v", errors.Wrap(err, "mem_index.RevertedIndexSyncManager running error"))
			cancel()
		})
	}
	g.Run()
}
