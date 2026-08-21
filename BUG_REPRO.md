# Bug 复现记录

## Bug 是什么

请求已经超时，网关还在读 DIMSE、等路由退避，心跳也继续打下游。不要修改代码，帮我查清取消信号在哪些地方断了，工程就在当前目录。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/dicom/adapter -run '^TestSessionReadPDUHonorsCancellation$' -count=1`
`go test ./internal/platform -run '^TestHTTPDoPropagatesCancellation$' -count=1`
`go test ./internal/routing/application -run '^TestRetryWaitStopsOnCancellation$' -count=1`
`go test ./internal/worker/application -run '^TestHeartbeatStopsCallbackAfterCancel$' -count=1`

## 错误信息

- TestSessionReadPDUHonorsCancellation：dimse_test.go:21 — expected cancellation-aware read, got read pdu header: EOF（读取 PDU 未在取消时及时返回，最终以 EOF 兜底）
- TestHTTPDoPropagatesCancellation：httpclient_test.go:27 — err=external request: Get "http://example.invalid": transport safety timeout（HTTP Do 未响应取消，由 transport safety timeout 兜底而非取消错误）
- TestRetryWaitStopsOnCancellation：retry_test.go:16 — err=<nil> elapsed=1.001109875s（重试等待未在取消时提前停止，完整耗尽约 1s 后才返回且无错误）
- TestHeartbeatStopsCallbackAfterCancel：runner_test.go:47 — callback ran after cancel（取消后心跳回调仍被触发）
