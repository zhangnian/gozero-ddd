package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config REST API 应用配置
type Config struct {
	rest.RestConf             // go-zero REST 配置
	MySQL         MySQLConfig `json:",optional"`      // MySQL 配置
	Redis         RedisConfig `json:",optional"`      // Redis 配置
	UseMemory     bool        `json:",default=false"` // 是否使用内存存储（开发测试用）
}

// RpcConfig gRPC 服务配置
// 组合了 go-zero 的 RpcServerConf 和自定义配置
type RpcConfig struct {
	zrpc.RpcServerConf             // go-zero gRPC 服务配置
	MySQL              MySQLConfig `json:",optional"`      // MySQL 配置
	UseMemory          bool        `json:",default=false"` // 是否使用内存存储
}

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	DataSource  string `json:",optional"`      // 数据源 DSN
	AutoMigrate bool   `json:",default=false"` // 是否自动迁移表结构
}

// RedisConfig 缓存配置
type RedisConfig struct {
	Host     string `json:",optional"`
	Password string `json:",optional"`
	DB       int    `json:",default=0"`
}
