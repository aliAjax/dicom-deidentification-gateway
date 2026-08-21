# Bug 复现记录

## Bug 是什么

并发查询时，旧的影像、脱敏规则和路由结果会自己变化，竞态检测也会报共享读写。不要修改代码，帮我把内部引用泄露的位置查出来。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test -race ./internal/dicom/infrastructure -run '^TestMemoryStoreSnapshotIsolation$' -count=1`
`go test -race ./internal/deidentification/infrastructure -run '^TestProfileSnapshotIsolation$' -count=1`
`go test -race ./internal/routing/infrastructure -run '^TestRoutingTargetSnapshotIsolation$' -count=1`
`go test -race ./internal/exporter/infrastructure -run '^TestExportSnapshotIsolation$' -count=1`

## 错误信息

- `TestMemoryStoreSnapshotIsolation`（internal/dicom/infrastructure, memory_test.go:51）：`snapshot=map[PatientID:writer]` — 读取方获取的快照泄露了并发写入方的值 `writer`，快照未被隔离。
- `TestProfileSnapshotIsolation`（internal/deidentification/infrastructure, service_test.go:36）：`profile=domain.Profile{ID:"p", Name:"safe", Rules:[]domain.Rule{domain.Rule{Tag:"reader", Action:"remove", Value:""}}, KeepPrivateTags:false, DateShiftDays:0}` — 快照中带有并发修改方写入的 Rule（Tag:"reader"），快照未被隔离。
- `TestRoutingTargetSnapshotIsolation`（internal/routing/infrastructure, memory_test.go:23）：`target=domain.Target{ID:"t", Name:"archive", ..., Tags:map[string]string{"Modality":"MR"}}` — 快照的 Tags 泄露了并发写入方的修改 `Modality:MR`，快照未被隔离。
- `TestExportSnapshotIsolation`（internal/exporter/infrastructure, memory_test.go:23）：`export=domain.Export{ID:"e", ..., Files:[]string{"reader"}}` — 快照的 Files 切片泄露了并发写入方的修改 `reader`，快照未被隔离。

四项测试均按预期暴露了同一类问题：读取方拿到的快照在并发写入下被污染，未实现快照隔离，符合题面描述的目标行为，BUG 存在。
