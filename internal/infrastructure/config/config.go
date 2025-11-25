package config

import "github.com/zeromicro/go-zero/rest"

// Config 应用配置
type Config struct {
	rest.RestConf            // go-zero REST 配置
	Database      Database   `json:",optional"` // 数据库配置
	Redis         Redis      `json:",optional"` // Redis配置
}

// Database 数据库配置
type Database struct {
	Driver string `json:",default=memory"` // 驱动类型: memory, mysql, postgres
	DSN    string `json:",optional"`       // 数据源名称
}

// Redis 缓存配置
type Redis struct {
	Host     string `json:",optional"`
	Password string `json:",optional"`
	DB       int    `json:",default=0"`
}

