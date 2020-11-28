package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"seedHabits/conf"
	"seedHabits/handler"
	"seedHabits/handler/dao"
	"seedHabits/sdk/log"
)

var (
	flagSet = flag.NewFlagSet("seed-api", flag.ExitOnError)
	cfgPath = flagSet.String("config", `./conf.conf.toml`, "Path of config files")
	version = flagSet.Bool("version", false, "show relate version info")
)

var (
	Version  string
	CommitId string
	Built    string
)

type program struct {
	env    string
	worker *handler.Worked
}

func (p *program)Init(env string)error  {
log.Init(conf.Config.LogConfig)
dao.Init()
w.Server 
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	err := conf.Init("./conf/conf.toml")
	if err != nil {
		fmt.Println("werr", err.Error())
	}

	log.Init(conf.Config.LogConfig)
	dao.Init()

	server := &handler.Server{}
	err = server.Init()
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}

	server.Launch()
	addr := fmt.Sprintf("%s :%d", conf.Config.Server.Host, conf.Config.Server.Port)
	log.Logger.Info(addr)
	log.Logger.Info(conf.Config.Database)
	//r := InitRouter()
	//r.Run()
}
