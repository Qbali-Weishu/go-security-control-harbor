package main

import (
	"log"
	"net/http"
	"os"

	"compatgate/internal/api"
	"compatgate/internal/policy"
)

func main() {
	root := os.Getenv("COMPAT_POLICY_ROOT")
	if root == "" {
		root = "testdata/policies"
	}
	catalog, err := policy.Load(root)
	if err != nil {
		log.Fatalf("load policy: %v", err)
	}
	server := api.NewServer(catalog)
	addr := ":8080"
	log.Printf("compatibility gate listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
