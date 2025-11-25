package main

import (
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"gozero-ddd/internal/infrastructure/config"
	"gozero-ddd/internal/interfaces/api/routes"
	"gozero-ddd/internal/interfaces/api/svc"
)

var configFile = flag.String("f", "etc/knowledge.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建 REST 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 创建服务上下文（依赖注入容器）
	ctx := svc.NewServiceContext(c)

	// 注册路由
	routes.RegisterRoutes(server, ctx)

	// 启动服务器
	fmt.Printf("🚀 知识库管理系统启动成功\n")
	fmt.Printf("📍 服务地址: http://%s:%d\n", c.Host, c.Port)
	fmt.Printf("📚 API 文档:\n")
	fmt.Printf("   POST   /api/v1/knowledge           - 创建知识库\n")
	fmt.Printf("   GET    /api/v1/knowledge           - 获取知识库列表\n")
	fmt.Printf("   GET    /api/v1/knowledge/:id       - 获取知识库详情\n")
	fmt.Printf("   PUT    /api/v1/knowledge/:id       - 更新知识库\n")
	fmt.Printf("   DELETE /api/v1/knowledge/:id       - 删除知识库\n")
	fmt.Printf("   POST   /api/v1/knowledge/:id/documents      - 添加文档\n")
	fmt.Printf("   GET    /api/v1/knowledge/:id/documents      - 获取文档列表\n")
	fmt.Printf("   DELETE /api/v1/knowledge/:id/documents/:doc_id - 删除文档\n")
	fmt.Printf("\n")

	server.Start()
}
