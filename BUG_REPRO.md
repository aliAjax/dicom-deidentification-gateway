# Bug 复现记录

## Bug 是什么

研究发送重试明明成功了，按 UID 查还是旧状态，有时记录还直接消失。不要修改代码，帮我沿着状态流转查明白，项目可以直接运行。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/study/domain -run '^TestStudyRetrySuccessTransition$' -count=1`
`go test ./internal/study/application -run '^TestAggregatorPreservesActiveState$' -count=1`
`go test ./internal/study/infrastructure -run '^TestStudyMemoryUpdateKeepsUIDIndex$' -count=1`
`go test ./internal/study/application -run '^TestStudyServiceReadsCompletedState$' -count=1`

## 错误信息

- TestStudyRetrySuccessTransition：`status_test.go:8: invalid study transition sending -> completed` —— 状态机不允许 sending → completed 转换。
- TestAggregatorPreservesActiveState：`study_test.go:23` 暴露 study 状态为 `"received"`，聚合器未保留/传播预期的 active 状态（字段全为零值/初始值）。
- TestStudyMemoryUpdateKeepsUIDIndex：`memory_test.go:17: stale uid index remains` —— 更新 study 后 UID 索引未同步刷新，残留旧索引。
- TestStudyServiceReadsCompletedState：`study_test.go:36: status=received` —— service 读取到的状态是 received 而非 completed，未正确读取已完成状态。
