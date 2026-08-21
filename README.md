# 13-dicom-deidentification-gateway

纯Go医学影像DICOM去标识化与跨机构传输网关。服务提供DICOM实例的HTTP封装接收、元数据校验、规则化去标识化、目标路由、传输状态和研究导出接口，并保留DIMSE TCP会话协议层骨架。

## 启动

```bash
go run ./cmd/server
```

默认监听`http://localhost:8093`，API请求需要`X-API-Key: dev-key`。DIMSE监听地址配置为`:11112`，可通过`DICOM_DIMSE_ADDR`修改。配置支持环境变量覆盖；`configs/config.yaml`和迁移文件用于部署参考。

## HTTP工作流

```bash
curl http://localhost:8093/healthz
curl -H 'X-API-Key: dev-key' --data-binary @examples/sample-instance.txt http://localhost:8093/v1/dicom/instances
curl -H 'X-API-Key: dev-key' -H 'Content-Type: application/json' -d '{"name":"research","date_shift_days":-30,"rules":[{"tag":"PatientID","action":"pseudonymize"},{"tag":"PatientName","action":"remove"},{"tag":"StudyDate","action":"date_shift"}]}' http://localhost:8093/v1/deidentification/profiles
curl -H 'X-API-Key: dev-key' -H 'Content-Type: application/json' -d '{"profile_id":"PROFILE_ID"}' http://localhost:8093/v1/instances/INSTANCE_ID/deidentify
curl -H 'X-API-Key: dev-key' -H 'Content-Type: application/json' -d '{"study_id":"STUDY_ID","instance_ids":["INSTANCE_ID"]}' http://localhost:8093/v1/exports
```

实例接收支持真实DICOM文件（识别`DICM`前导码）和用于开发验证的键值封装样例。重复内容摘要会返回已有实例，超大请求会被拒绝。敏感字段不写入结构化日志。

## 领域结构

`internal/dicom`负责PDU、Association和数据集元数据；`deidentification`负责五类规则和确定性伪名化；`study`负责Patient/Study/Series聚合；`routing`、`transfer`和`spool`负责目标选择、可靠发送和重启恢复；`exporter`负责研究导出清单；`worker`提供租约和退避基础设施。外部存储、连接器、时钟和密码学均通过接口注入。

## 验证

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/server
```

生产环境应使用PostgreSQL迁移、独立对象存储、mTLS/短期令牌和受控DICOM AE白名单替换开发默认值。
