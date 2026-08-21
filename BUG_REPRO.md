# Bug 复现记录

## Bug 是什么

影像接收或转发失败后，上层只能看到一段普通文本，原来的存储错误和网络错误都认不出来了。当前目录能直接检查，请帮我修好错误传递。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/platform -run '^TestInternalErrorPreservesSentinel$' -count=1`
`go test ./internal/spool/application -run '^TestQueuePreservesWriteSentinel$' -count=1`
`go test ./internal/transfer/application -run '^TestSendPreservesConnectorSentinel$' -count=1`
`go test ./internal/dicom/application -run '^TestReceivePreservesParserSentinel$' -count=1`

## 错误信息

- `internal/platform` → `errors_test.go:11: sentinel was not preserved` —— InternalError 未保留 sentinel 错误。
- `internal/spool/application` → `queue_test.go:26: err=write spool: spool volume offline` —— Queue 的 Write 把底层 sentinel 包裹成了带前缀的新错误信息，sentinel 被掩盖。
- `internal/transfer/application` → `service_test.go:31: err=send instance: association refused` —— Send 把 connector sentinel 包裹为 `send instance: ...`，未保留原始 sentinel。
- `internal/dicom/application` → `service_test.go:30: err=parse instance: dataset truncated` —— Receive 把 parser sentinel 包裹为 `parse instance: ...`，未保留原始 sentinel。
