package config

import "github.com/zeromicro/go-zero/rest"

// Config 应用配置
type Config struct {
	rest.RestConf               // go-zero REST 配置
	MySQL         MySQLConfig   `json:",optional"` // MySQL 配置
	Redis         RedisConfig   `json:",optional"` // Redis 配置
	UseMemory     bool          `json:",default=false"` // 是否使用内存存储（开发测试用）
}

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	DataSource string `json:",optional"` // 数据源 DSN
}

// RedisConfig 缓存配置
type RedisConfig struct {
	Host     string `json:",optional"`
	Password string `json:",optional"`
	DB       int    `json:",default=0"`
}
