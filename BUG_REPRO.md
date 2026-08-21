# Bug 复现记录

## Bug 是什么

收到截断的 DICOM 帧时服务会直接 panic，PDU 和 command 两种包都见过。先别改代码，麻烦查清四种长度路径怎么越界的，当前工程能复现。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/dicom/domain -run '^TestDecodePDURejectsTruncatedBody$' -count=1`
`go test ./internal/dicom/adapter -run '^TestDecodeCommandRejectsTruncatedElement$' -count=1`
`go test ./internal/dicom/adapter -run '^TestParseElementRejectsOversizedValue$' -count=1`
`go test ./internal/dicom/infrastructure -run '^TestLimitsRejectsOverflowFrame$' -count=1`

## 错误信息

- `TestDecodePDURejectsTruncatedBody`：`DecodePDU` 在 `pdu.go:18` 发生 panic —— `runtime error: slice bounds out of range [:14] with capacity 7`。预期应拒绝截断的 PDU body（返回错误），实际未做长度校验直接切片导致越界 panic。
- `TestDecodeCommandRejectsTruncatedElement`：`DecodeCommand` 在 `commands.go:59` 发生 panic —— `slice bounds out of range [:134217736] with capacity 9`。预期应拒绝截断元素，实际直接按声明长度切片，越界 panic。
- `TestParseElementRejectsOversizedValue`：`ParseElement` 在 `metadata.go:21` 发生 panic —— `slice bounds out of range [:67108872] with capacity 9`。预期应拒绝超大 value，实际未做上限校验直接切片，越界 panic。
- `TestLimitsRejectsOverflowFrame`：`limits_test.go:8` 报告 `oversized frame allowed`。预期应拒绝超限 frame，实际超限帧被放行通过。
