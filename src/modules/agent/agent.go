package main

import (
	"context"
	"log"
	"open-devops/src/modules/agent/config"
	"open-devops/src/modules/agent/info"
	"open-devops/src/modules/agent/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/oklog/run"

	"github.com/prometheus/common/promlog"
	promlogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// 命令行解析
	app = kingpin.New(filepath.Base(os.Args[0]), "The open-devops-agent")
	// 指定配置文件参数
	configFile = app.Flag("config.file", "open-devops-agent configuration file").Short('c').Default("open-devops-agent.yaml").String()
)

func main() {
	// Help Info
	app.HelpFlag.Short('h')
	// 基于prometheus库通过build ldflags版本信息注入
	promlogConfig := promlog.Config{}
	app.Version(version.Print("open-devops-agent"))
	promlogflag.AddFlags(app, &promlogConfig)

	// Get Param From Command line from $1 to the end
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.Println("Start loading config...")
	agentConfig, err := config.LoadFile(*configFile)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	log.Println("Loading config successfully")

	rpcCli := rpc.InitRpcCli(agentConfig.RPCServerAddr)

	if err := rpcCli.Ping(); err != nil {
		log.Printf("%+v\n", err)
		return
	}

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
	// 采集基础信息
	{
		g.Add(func() error {
			err := info.TickerInfoCollectAndReport(rpcCli, ctx)
			return err
		}, func(err error) {
			cancel()
		})
	}
	g.Run()
}
