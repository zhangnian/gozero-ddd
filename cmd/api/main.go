package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	// 创建服务上下文（依赖注入容器）
	ctx := svc.NewServiceContext(c)

	// 注册路由
	routes.RegisterRoutes(server, ctx)

	// 打印启动信息
	fmt.Printf("🚀 知识库管理系统启动成功\n")
	fmt.Printf("📍 服务地址: http://%s:%d\n", c.Host, c.Port)
	fmt.Printf("📚 API 文档:\n")
	fmt.Printf("   POST   /api/v1/knowledge           - 创建知识库\n")
	fmt.Printf("   GET    /api/v1/knowledge           - 获取知识库列表\n")
	fmt.Printf("   GET    /api/v1/knowledge/:id       - 获取知识库详情\n")
	fmt.Printf("   PUT    /api/v1/knowledge/:id       - 更新知识库\n")
	fmt.Printf("   DELETE /api/v1/knowledge/:id       - 删除知识库\n")
	fmt.Printf("   POST   /api/v1/knowledge/merge     - 合并知识库（事务演示）\n")
	fmt.Printf("   POST   /api/v1/knowledge/:id/documents      - 添加文档\n")
	fmt.Printf("   GET    /api/v1/knowledge/:id/documents      - 获取文档列表\n")
	fmt.Printf("   DELETE /api/v1/knowledge/:id/documents/:doc_id - 删除文档\n")
	fmt.Printf("\n")

	// 优雅关闭
	// 监听系统信号，实现优雅停机
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		fmt.Println("\n🛑 收到关闭信号，正在优雅关闭...")
		server.Stop()
	}()

	// 启动服务器
	server.Start()
}
