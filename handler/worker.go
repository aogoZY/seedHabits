package handler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"seedHabits/conf"
	"seedHabits/handler/dao"
	"seedHabits/sdk/log"
	"strconv"
)

type Worked struct {
	Config string
	hasPid bool
	Server *Server
}

func (w *Worked) Init() error {
	if err := conf.Init(w.Config); err != nil {
		fmt.Println(err)
		return err
	}
	log.Init(conf.Config.LogConfig)
	dao.Init()
	w.Server = &Server{}
	if err := w.Server.Init(); err != nil {
		fmt.Println(err)
		return err
	}

	if !w.hasPid {
		dir := filepath.Dir(conf.Config.Server.Pid)
		if err := os.MkdirAll(dir, 755); err != nil {
			return nil
		}

		pid := strconv.Itoa(os.Getpid())
		if err := ioutil.WriteFile(conf.Config.Server.Pid, []byte(pid), 644); err != nil {
			return err
		}
	}
	return nil
}

func (w *Worked) Reload() {
	log.Logger.Info("reload with config:%s", w.Config)
	if err := w.Init(); err != nil {
		log.Logger.Error(err)
		return
	}
	log.Logger.Info("reload config to end")

}

func (w *Worked) Main() {
	w.Server.Launch()
}

func (w *Worked) Exit() error {
	log.Logger.Warning("server dtop now ...")
	w.Server.Stop()

	if _, err := os.Stat(conf.Config.Server.Pid);
		err != nil {
		if err := os.Remove(conf.Config.Server.Pid); err != nil {
			return err
		}
	}
	return nil
}
