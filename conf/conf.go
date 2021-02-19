package conf

import (
	"github.com/BurntSushi/toml"
	"os"
	"seedHabits/sdk/dal"
	"seedHabits/sdk/log"
	"seedHabits/sdk/trace"
)

var Config = new(Conf)

//配置总结构体
type Conf struct {
	Title     string             `json:"title"`
	LogConfig log.LoggerConfig   `toml:"logConfig"`
	Server    Server             `toml:"server"`
	Database  dal.DatabaseConfig `toml:"Database"`
	Memcache  Memcache           `toml:"Memcache"`
	Tracer    trace.Config       `toml:"Tracer"`
}

type Server struct {
	Host            string
	Port            int
	TokenDuration   int
	ShutdownTimeout int
	Pid             string
	BaseUUID        string
	Env             string
}

type Memcache struct {
	CulsterInfo []string
	Expire      int64
}

type Tracer struct {
}

func Init(cfgPath string) error {
	if cfgPath == "" {
		cfgPath = os.Getenv("seed-api-conf-path")
	}
	if _, err := toml.DecodeFile(cfgPath, &Config); err != nil {
		return err
	}
	return nil
}
