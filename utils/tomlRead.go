package utils

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/howeyc/fsnotify"
)

type tomlConfig struct {
	Environment string
	DB          database `toml:"database"`
	Cache       cache    `toml:"cache"`
	Auth        auth     `toml:"auth"`
	Logs        logs     `toml:"logs"`
	Cors        cors     `toml:"cors"`
	Produrce    produrce `toml:"produrce"`
	Port        port     `toml:"port"`
}
type port struct {
	Port string
}

type cache struct {
	Cluster cluster
	Single  single
}
type cluster struct {
	Addrs []string
}

type single struct {
	Server   string
	PassWord string
	DB       int
}

type produrce struct {
	Server string
	Port   string
}
type cors struct {
	Origin string
}
type database struct {
	Type       string
	DriverName string
	Server     string
	Port       string
	UserName   string
	SID        string
	PassWord   string
	NlsLang    string
}

type auth struct {
	Server string
}

type logs struct {
	FilePath  string
	Level     int
	Formatter string
}

var config tomlConfig

func init() {
	InitAll()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsAttrib() {
					InitAll()
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch("app.toml")
	if err != nil {
		log.Fatal(err)
	}
}

func InitAll() {
	LoadToml()

	InitDB()

}

func LoadToml() {
	if _, err := toml.DecodeFile("app.toml", &config); err != nil {
		panic("初始化 app.toml 配置文件失败\n" + err.Error())
	}
	InitLogs()
	SysLog.Debug("初始化 app.toml 配置文件成功", config)
}

// GetConfig func 获取数据库配置信息
func GetConfig() tomlConfig {
	return config
}

// GetConfig func 获取数据库配置信息
func GetProdurceConfig() produrce {
	return config.Produrce
}

// GetDataBaseCfg func 获取数据库配置信息
func GetDataBaseCfg() database {
	return config.DB
}

//GetCacheCfg 获取缓存配置信息
func GetCacheCfg() cache {
	return config.Cache
}

//GetCacheCfg 获取端口配置信息
func GetPort() port {
	return config.Port
}
