// gRPC 客户端示例
// 演示如何调用 KnowledgeService 的 gRPC 接口
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gozero-ddd/internal/interfaces/rpc/pb"
)

func main() {
	// 1. 创建 gRPC 连接
	conn, err := grpc.Dial(
		"localhost:9999",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("❌ 连接 gRPC 服务失败: %v", err)
	}
	defer conn.Close()

	// 2. 创建客户端
	client := pb.NewKnowledgeServiceClient(conn)

	// 3. 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ==================== 演示：创建知识库 ====================
	fmt.Println("📝 演示 CreateKnowledgeBase 接口（Command 操作）")
	fmt.Println("================================================")

	createResp, err := client.CreateKnowledgeBase(ctx, &pb.CreateKnowledgeBaseRequest{
		Name:        "gRPC 测试知识库",
		Description: "这是通过 gRPC 接口创建的知识库，演示 Command 操作",
	})
	if err != nil {
		log.Fatalf("❌ 创建知识库失败: %v", err)
	}

	kb := createResp.KnowledgeBase
	fmt.Printf("✅ 创建成功!\n")
	fmt.Printf("   ID: %s\n", kb.Id)
	fmt.Printf("   名称: %s\n", kb.Name)
	fmt.Printf("   描述: %s\n", kb.Description)
	fmt.Printf("   创建时间: %s\n", time.Unix(kb.CreatedAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Println()

	// ==================== 演示：获取知识库 ====================
	fmt.Println("🔍 演示 GetKnowledgeBase 接口（Query 操作）")
	fmt.Println("================================================")

	getResp, err := client.GetKnowledgeBase(ctx, &pb.GetKnowledgeBaseRequest{
		Id:               kb.Id,
		IncludeDocuments: true,
	})
	if err != nil {
		log.Fatalf("❌ 获取知识库失败: %v", err)
	}

	kbDetail := getResp.KnowledgeBase
	fmt.Printf("✅ 查询成功!\n")
	fmt.Printf("   ID: %s\n", kbDetail.Id)
	fmt.Printf("   名称: %s\n", kbDetail.Name)
	fmt.Printf("   描述: %s\n", kbDetail.Description)
	fmt.Printf("   文档数量: %d\n", kbDetail.DocumentCount)
	fmt.Printf("   更新时间: %s\n", time.Unix(kbDetail.UpdatedAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Println()

	// ==================== DDD 架构说明 ====================
	fmt.Println("📚 DDD + go-zero gRPC 架构说明")
	fmt.Println("================================================")
	fmt.Println(`
请求处理流程：
  gRPC Request 
    → Server (实现 gRPC 接口) 
    → Logic (业务逻辑协调) 
    → Command/Query Handler (应用层) 
    → Domain Service (领域服务) 
    → Repository (仓储) 
    → Database

分层职责：
  1. interfaces/rpc/pb       - Protocol Buffer 定义和生成代码
  2. interfaces/rpc/server   - gRPC 服务实现，创建 Logic 实例
  3. interfaces/rpc/logic    - 业务逻辑协调，调用应用层
  4. application/command     - 命令处理器（写操作）
  5. application/query       - 查询处理器（读操作）
  6. domain/service          - 领域服务（核心业务逻辑）
  7. domain/entity           - 领域实体（业务模型）
  8. infrastructure/persist  - 仓储实现（数据持久化）
`)

	fmt.Println("🎉 gRPC 接口演示完成!")
}

