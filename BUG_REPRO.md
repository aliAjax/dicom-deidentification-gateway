# Bug 复现记录

## Bug 是什么

spool 批量恢复跑一阵就耗尽资源，某项处理失败后原错误也没了，状态还可能留在处理中。麻烦帮我修好释放和回滚，代码在当前目录。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/spool/application -run '^TestRecoverClosesEachItemImmediately$' -count=1`
`go test ./internal/spool/application -run '^TestRecoverReadFailureMarksFailed$' -count=1`
`go test ./internal/spool/application -run '^TestRecoverKeepsFailedItemPending$' -count=1`
`go test ./internal/spool/application -run '^TestRecoverPreservesMarkError$' -count=1`
`go test ./internal/spool/application -run '^TestRecoverStopsWhenContextCancelled$' -count=1`
`go test ./internal/spool/infrastructure -run '^TestFilesReleaseRemovesSpoolPath$' -count=1`
`go test ./internal/spool/infrastructure -run '^TestMemoryMarkRejectsUnknownItem$' -count=1`

## 错误信息

- TestRecoverClosesEachItemImmediately：recover_test.go:58 `releases before b=[]`，期望在处理到 b 之前已释放前置项，实际为空，说明 Recover 未及时关闭每个已处理项。
- TestRecoverReadFailureMarksFailed：recover_test.go:77 `marks=[a:pending]`，期望读取失败时将 a 标记为 failed，实际仍为 pending，说明读取失败路径未做失败标记。
- TestRecoverKeepsFailedItemPending：recover_test.go:49 `marks=[spool-1:failed]`，期望已失败项保持 pending，实际被改写为 failed，说明 Recover 错误地改写了既有失败项状态。
- TestRecoverPreservesMarkError：recover_test.go:86 `err=<nil>`，期望保留 Mark 返回的错误，实际错误被吞为 nil，说明标记错误未被透传。
- TestRecoverStopsWhenContextCancelled：recover_test.go:96 `err=<nil> called=true`，期望上下文取消时停止且不继续调用下游，实际无错误且已继续调用，说明取消未生效、未提前停止。
- TestFilesReleaseRemovesSpoolPath：files_test.go:15 `path still exists: <nil>`，Release 后 spool 路径仍存在，说明 Files Release 未删除 spool 路径。
- TestMemoryMarkRejectsUnknownItem：memory_test.go:9 `expected missing item error`，对未知项调用 Mark 期望返回错误，实际未返回，说明内存实现未拒绝未知项。
