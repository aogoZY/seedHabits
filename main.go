package main

import (
	"flag"
	"fmt"
	"github.com/kwanhur/go-svc/svc"
	"os"
	"path/filepath"
	"runtime"
	"seedHabits/handler"
	"seedHabits/sdk/log"
	"syscall"
)

var (
	flagSet = flag.NewFlagSet("seed-api", flag.ExitOnError)
	cfgPath = flagSet.String("config", `./conf/conf.toml`, "Path of config files")
	version = flagSet.Bool("version", false, "show relate version info")
)

var (
	Version  string
	CommitId string
	Built    string
)

type program struct {
	env    svc.Environment
	worker *handler.Worked
}

func (p *program) Init(env svc.Environment) error {
	// 检查是否是windows 服务。。。目测一般时候也用不到
	p.env = env
	if env.IsWindowsService(){
		dir:=filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	flagSet.Parse(os.Args[1:])
	if *version {
		fmt.Printf("commit id :%s\n", CommitId)
		fmt.Printf("build by %s %s / %s at %s \n", runtime.Version(), runtime.GOOS, runtime.GOARCH, Built)
		os.Exit(2)
	}
	daemon := &handler.Worked{
		Config: *cfgPath,
		Server: &handler.Server{},
	}
	err := daemon.Init()
	if err != nil {
		return err
	}
	p.worker = daemon
	return nil
}

func (p *program) Start() error {
	log.Logger.Info("starting...")
	if p.worker != nil {
		p.worker.Main()
	} else {
		log.Logger.Warning("worker is nil")
	}
	return nil
}

func (p *program) Stop() error {
	log.Logger.Warning("stopping")
	if p.worker != nil {
		p.worker.Exit()
	} else {
		log.Logger.Warning("worker is nil")
	}
	return nil
}

func (p *program) Reload(signal os.Signal) () {
	log.Logger.Info("got signal %s", signal.String())
	switch signal {
	case syscall.SIGHUP:
		p.worker.Reload()
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	pg := &program{}
	svc.Notify(syscall.SIGHUP, pg.Reload)
	if err := svc.Run(pg, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT); err != nil {
		log.Logger.Error(err)
		os.Exit(2)
	} else {
		log.Logger.Info("bye ~")
	}

	//r := InitRouter()
	//r.Run()
}
