package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// DatabaseConfig 数据库配置
type databaseConfig struct {
	Host     string // 主机
	Port     int    // 端口
	Username string // 用户名
	Password string // 密码
	Name     string // 数据库名
	Type     string // 类型
}

// ServerConfig 服务器信息
type serverConfig struct {
	IP      string // 主机
	Port    int    // 端口
	AppKey  string // app key
	Mode    string // 模式
	Webhook string // webhook地址
	LRUSize int    `toml:"lruSize"` // lru缓存大小
}

// redis配置
type redisConfig struct {
	Network  string // 网络类型
	Address  string // 连接地址
	Password string // 密码
	DB       int    // 数据库
}

// Config 配置
type Config struct {
	Title string
	DB    databaseConfig `toml:"database"`
	Svr   serverConfig   `toml:"server"`
	Redis redisConfig    `toml:"redis"`
}

// Conf 全局配置
var Conf Config

func init() {
	if len(os.Args) < 2 {
		panic("no config file")
	}
	file := os.Args[1]
	if _, err := toml.DecodeFile(file, &Conf); err != nil {
		panic(err)
	}
}
