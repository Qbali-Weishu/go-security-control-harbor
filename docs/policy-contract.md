# 兼容性策略契约

`as_of = 2026-07-01T00:00:00Z`

此服务由防御工程团队使用，在多控制部署包进入受监管或内部网络段之前对其进行预筛选。该服务必须仅从捆绑的策略文件中产生确定性答案。

## 强制性评估维度

### 1. 平台支持

每个选定的组件必须支持：

- 配置文件的操作系统
- 配置文件的运行时
- 配置文件的内核版本

内核范围严格遵循此规则：

- `min` 是包含的
- `max_exclusive` 是排除的

如果配置文件的内核等于 `max_exclusive`，则必须拒绝该组件。

### 2. 条件性前提

组件可以声明：

- `requires`：当其条件匹配配置文件时，每个匹配条目都是强制性的
- `requires_any`：当其条件匹配配置文件时，必须存在至少一个候选项

评估器必须检查所有匹配的条件，而不仅仅是无条件条目。

### 3. 定向数据路径保护

如果选定的包包含任何 `raw_payload_emitter` 且配置文件区域为 `regulated`，则请求必须在 `data_path` 中第一个收集器之前包含接受的清理器。

接受的清理器取决于事件状态：

- `steady`：`content-sanitizer` 或 `payload-redactor`
- `containment`：`content-sanitizer` 或 `payload-redactor`
- `eradication`：仅 `content-sanitizer`

仅存在是不够的。顺序很重要：

- 接受的清理器必须出现在 `central-collector` 之前
- 如果 `egress_mode = restricted`，`telemetry-relay` 必须出现在 `central-collector` 之后
- 如果 `egress_mode = restricted` 且 `fips_mode` 在 `flow.auditor_required_fips_modes` 中列出，则 `egress-auditor` 必须出现在 `telemetry-relay` 之后

### 4. 冲突和豁免

冲突是双向的。如果组件 `A` 声明与 `B` 冲突，即使 `B` 没有重复声明，选择 `B` 和 `A` 仍然无效。

只有当以下所有条件为真时，才能豁免阻塞器：

- 规则标记为 `waivable`
- 请求包含批准工单
- 批准的 `rule_code` 匹配
- 批准的组件对匹配，不考虑顺序
- 配置文件 ID、区域和事件状态在批准范围内
- 如果批准声明了 `fips_modes`，则配置文件的 FIPS 模式在范围内
- 如果批准声明了 `egress_modes`，则配置文件的出口模式在范围内
- `profile.as_of < approval.expires_at`

如果这些检查中的任何一项失败，则豁免无效，阻塞器保持活动状态。

### 5. 预算聚合、评分和警告建议

总组件开销是所有选定组件在每个预算维度上的总和：

- `cpu_milli`
- `memory_mb`
- `hook_units`

当任何总和维度超过其预算时，配置文件不兼容。

如果包保持兼容但任何利用率大于或等于 `flow.warning_utilization`，响应仍必须为每个超过阈值的维度包含一个建议性 `required_action`：

- `review cpu headroom before rollout`
- `review memory headroom before rollout`
- `review hook headroom before rollout`

评分语义是固定的：

- 被阻塞的包返回 `score = 0`
- 兼容的包返回 `score = 1 - average(cpu_ratio, memory_ratio, hook_ratio)`
- 评分四舍五入到两位小数

### 6. 确定性输出

响应必须是确定性的：

- 阻塞器按 `code` 排序，然后按组件列表排序
- `required_actions` 去重并排序，包括警告建议
- 跟踪总计必须反映用于通过/失败和评分计算的相同预算，即使多个维度同时失败
