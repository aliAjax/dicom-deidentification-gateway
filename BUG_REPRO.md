# Bug 复现记录

## Bug 是什么

导出清单预览后再追加一张影像，之前那份清单的顺序和内容也会跟着动，落盘摘要偶尔还对不上。麻烦帮我修好这次数据串场。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/exporter/domain -run '^TestCanonicalPathsOwnsBackingArray$' -count=1`
`go test ./internal/exporter/application -run '^TestBuildManifestLeavesInputUntouched$' -count=1`
`go test ./internal/exporter/infrastructure -run '^TestBuilderHashMatchesCanonicalManifest$' -count=1`
`go test ./internal/exporter/application -run '^TestExportServiceKeepsPathSnapshot$' -count=1`

## 错误信息

- `TestCanonicalPathsOwnsBackingArray`：`manifest_test.go:10: input mutated: [changed b.dcm]` —— 规范化函数没有拥有/复制底层数组，直接改写了调用方传入的切片。
- `TestBuildManifestLeavesInputUntouched`：`manifest_test.go:13: manifest="\n" input=[b.dcm a.dcm]` —— 构建清单时把输入路径变成了 `\n`（用 `strings.Join` 之类把切片替换成了字符串），且未做副本，输入被破坏。
- `TestBuilderHashMatchesCanonicalManifest`：`builder_test.go:20: manifest="b.dcm\na.dcm"` —— infrastructure 的 builder 产出的清单内容与 domain 规范化结果不一致，哈希无法匹配。
- `TestExportServiceKeepsPathSnapshot`：`manifest_test.go:41: export aliases caller: ... Files:[]string{"one.dcm", "two.dcm"} ... Files:[]string{"mutated", "changed"}` —— 导出服务返回的 Export 与调用方持有同一底层数组（别名），调用方后续修改波及到了 export 对象，未做快照。
