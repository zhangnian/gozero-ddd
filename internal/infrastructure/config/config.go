package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config REST API 应用配置
type Config struct {
	rest.RestConf             // go-zero REST 配置
	MySQL         MySQLConfig `json:",optional"` // MySQL 配置
	Redis         RedisConfig `json:",optional"` // Redis 配置
	Kafka         KafkaConfig `json:",optional"` // Kafka 配置
	UseKafka      bool        `json:",default=false"` // 是否使用 Kafka 事件总线
}

// RpcConfig gRPC 服务配置
// 组合了 go-zero 的 RpcServerConf 和自定义配置
type RpcConfig struct {
	zrpc.RpcServerConf             // go-zero gRPC 服务配置
	MySQL              MySQLConfig `json:",optional"` // MySQL 配置
	Kafka              KafkaConfig `json:",optional"` // Kafka 配置
	UseKafka           bool        `json:",default=false"` // 是否使用 Kafka 事件总线
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

// KafkaConfig Kafka 消息队列配置
type KafkaConfig struct {
	Brokers         []string      `json:",optional"`               // Kafka broker 地址列表
	Topic           string        `json:",default=domain-events"`  // 事件主题
	GroupID         string        `json:",default=knowledge-service"` // 消费者组ID
	WriteTimeout    time.Duration `json:",default=10s"`            // 写入超时
	ReadTimeout     time.Duration `json:",default=10s"`            // 读取超时
	BatchSize       int           `json:",default=100"`            // 批量发送大小
	BatchTimeout    time.Duration `json:",default=1s"`             // 批量发送超时
	RequiredAcks    int           `json:",default=-1"`             // 确认模式: -1=all, 0=none, 1=leader
	Async           bool          `json:",default=false"`          // 是否异步发送
	AutoCreateTopic bool          `json:",default=true"`           // 是否自动创建主题
}
