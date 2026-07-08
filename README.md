# 安全控制兼容性网关

此服务评估建议的安全控制包是否可以在托管的 Linux 主机配置文件上部署，而不违反平台、数据路径、豁免或资源约束。

## 运行时契约

- `GET /healthz` 返回基本的存活文档。
- `POST /v1/compatibility/assess` 接受包请求并返回确定性的兼容性决策。
- 策略文件默认从 `testdata/policies/` 加载，如果设置了 `COMPAT_POLICY_ROOT` 则从该路径加载。

## 本地运行

```bash
go run ./cmd/compatibilityd
```

服务器运行后，使用环境中可用的 HTTP 客户端提交 `testdata/requests/` 下的任何 JSON 请求。

## 关键工程约束

评估器必须同时跨多个维度进行推理：

- 平台支持：OS、运行时、内核范围
- 前提闭包：直接和条件性需求
- 受监管区域中的定向数据流安全
- 带工单范围豁免的冲突处理
- CPU、内存和钩子预算总和
- 下游审计员的确定性跟踪输出
