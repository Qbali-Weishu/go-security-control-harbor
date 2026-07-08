package main

import (
	"log"
	"net/http"
	"os"

	"compatgate/internal/api"
	"compatgate/internal/policy"
)

func main() {
	// 从环境变量获取策略根目录，默认为 testdata/policies
	root := os.Getenv("COMPAT_POLICY_ROOT")
	if root == "" {
		root = "testdata/policies"
	}
	catalog, err := policy.Load(root)
	if err != nil {
		log.Fatalf("加载策略失败: %v", err)
	}
	server := api.NewServer(catalog)
	addr := ":8080"
	log.Printf("兼容性网关监听在 %s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatalf("服务器退出: %v", err)
	}
}
