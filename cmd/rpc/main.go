package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gozero-ddd/internal/interfaces/rpc/pb"
	"gozero-ddd/internal/interfaces/rpc/server"
	"gozero-ddd/internal/interfaces/rpc/svc"
)

var configFile = flag.String("f", "etc/knowledge-rpc.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 1. 加载配置
	var c svc.RpcConfig
	conf.MustLoad(*configFile, &c)

	// 2. 创建服务上下文（依赖注入容器）
	// ServiceContext 初始化所有依赖：数据库连接、仓储、领域服务、处理器等
	ctx := svc.NewServiceContext(c)

	// 3. 创建 gRPC 服务器
	// go-zero 的 zrpc.MustNewServer 封装了 gRPC 服务器的创建
	// 自动添加了拦截器、中间件等
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 4. 注册 gRPC 服务
		// KnowledgeServer 实现了 pb.KnowledgeServiceServer 接口
		pb.RegisterKnowledgeServiceServer(grpcServer, server.NewKnowledgeServer(ctx))

		// 5. 注册反射服务（开发调试用）
		// 允许使用 grpcurl 等工具查询服务方法
		if c.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 打印启动信息
	fmt.Printf("🚀 知识库管理系统 gRPC 服务启动成功\n")
	fmt.Printf("📍 服务地址: %s\n", c.ListenOn)
	fmt.Printf("📚 gRPC 接口:\n")
	fmt.Printf("   GetKnowledgeBase    - 获取知识库详情（Query 演示）\n")
	fmt.Printf("   CreateKnowledgeBase - 创建知识库（Command 演示）\n")
	fmt.Printf("\n")
	fmt.Printf("💡 测试命令:\n")
	fmt.Printf("   # 使用 grpcurl 测试（需要先安装 grpcurl）\n")
	fmt.Printf("   grpcurl -plaintext -d '{\"name\":\"测试知识库\",\"description\":\"这是一个测试\"}' localhost:9999 knowledge.KnowledgeService/CreateKnowledgeBase\n")
	fmt.Printf("   grpcurl -plaintext -d '{\"id\":\"<knowledge_base_id>\"}' localhost:9999 knowledge.KnowledgeService/GetKnowledgeBase\n")
	fmt.Printf("\n")

	// 6. 启动服务器
	s.Start()
}
