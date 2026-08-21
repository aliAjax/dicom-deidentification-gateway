# Bug 复现记录

## Bug 是什么

批量传输一并发就偶尔少结果，碰到取消时还会卡住，日志里出现过 close of closed channel。请修复这段 worker 生命周期，源码就在这里。

goroutine 18 [running]:
panic: send on closed channel
github.com/example/dicom-deidentification-gateway/internal/transfer/application.(*Batch).Run.func1()
internal/transfer/application/batch.go:24 +0x1a4
testing.tRunner.func1.2({0x102cf1060, 0x102d11a30})

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test -race ./internal/transfer/application -run '^TestBatchWorkersRegisterBeforeWait$' -count=1`
`go test -race ./internal/transfer/application -run '^TestBatchCancelStopsProducer$' -count=1`
`go test -race ./internal/transfer/application -run '^TestBatchResultChannelClosesOnce$' -count=1`
`go test -race ./internal/transfer/application -run '^TestBatchRunReturnsEveryAcceptedJob$' -count=1`

## 错误信息

- TestBatchWorkersRegisterBeforeWait：batch_test.go:36 — wait returned before workers stopped（wait 在 worker 停止之前就返回，说明等待逻辑未覆盖所有已注册 worker）
- TestBatchCancelStopsProducer：batch_test.go:56 — producer ignored cancellation（生产者忽略了取消信号，未随取消退出）
- TestBatchResultChannelClosesOnce：batch_test.go:68 — results=0（结果通道过早关闭/无结果，未能按预期收集结果）
- TestBatchRunReturnsEveryAcceptedJob：batch_test.go:21 — got 0 jobs, want 12（Run 未返回每个被接受的 job，期望 12 个实际 0 个）
