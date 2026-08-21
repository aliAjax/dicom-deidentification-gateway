# Bug 复现记录

## Bug 是什么

HTTP 请求中途失败后，下一次重试会读到半截状态：影像像是已收、脱敏标签又没写完，导出也偶尔提前完成。请修复这些公开路径的状态提交。

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/http -run '^TestInstanceReceiveRollsBackOnSpoolFailure$' -count=1`
`go test ./internal/http -run '^TestDeidentifyDoesNotPublishPartialTags$' -count=1`
`go test ./internal/http -run '^TestExportRequestPublishesAtomically$' -count=1`
`go test ./internal/http -run '^TestStudySendDoesNotInventCompletedState$' -count=1`

## 错误信息

- TestInstanceReceiveRollsBackOnSpoolFailure：state_test.go:33 status=201（spool 失败时未回滚，返回 201 而非失败状态）
- TestDeidentifyDoesNotPublishPartialTags：state_test.go:59 partial update=...Instance{..., Status:"deidentified", Tags:{"PatientID":"", "StudyInstanceUID":"1.2.3", "format":"encapsulated"}}（发布了部分标签，Status 已置为 deidentified，但 PatientID 为空、StudyInstanceUID 仍为原值）
- TestExportRequestPublishesAtomically：state_test.go:68 status=202（导出请求未原子发布，返回 202）
- TestStudySendDoesNotInventCompletedState：state_test.go:80 status=202（study send 凭空捏造了 completed 状态，返回 202）
