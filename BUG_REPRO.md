# Bug 复现记录

## Bug 是什么

我们用最小配置起服务后，第一次走鉴权或保存脱敏规则就会崩，空的目标探测也有同样情况。上线窗口快到了，请修复这些零值路径。

goroutine 3 [running]:
testing.tRunner.func1.2({0x102cf1060, 0x102d11a30})
/Users/zhuanzmima0000/.local/go/src/testing/testing.go:1631 +0x1c4
panic({0x102cf1060?, 0x102d11a30?})
/Users/zhuanzmima0000/Desktop/Go_project/new_work_projects/dicom-deidentification-batch-20260821/2026-08-21/dicom-deidentification-gateway__002/env/internal/config/config_test.go:8 +0x98

## 如何触发

在埋有该问题的项目目录中，逐条执行下面的定向测试：

`go test ./internal/config -run '^TestConfigNormalizeInitializesCollections$' -count=1`
`go test ./internal/auth -run '^TestMiddlewareRejectsTypedNilValidator$' -count=1`
`go test ./internal/deidentification/infrastructure -run '^TestFileProfileStoreZeroValueReturnsError$' -count=1`
`go test ./internal/routing/infrastructure -run '^TestHealthCheckerRejectsEmptyTarget$' -count=1`

## 错误信息

- `TestConfigNormalizeInitializesCollections`：`panic: assignment to entry in nil map`。Normalize 未对未初始化的 map 集合做初始化，直接对 nil map 赋值触发 panic（config_test.go:8）。
- `TestMiddlewareRejectsTypedNilValidator`：`panic: typed nil validator called`。中间件未拒绝 typed-nil 的 Validator，反而调用了它，触发了测试桩 `nilValidator.Validate` 的 panic（auth.go:34 / auth_test.go:13）。
- `TestFileProfileStoreZeroValueReturnsError`：`profile_store_test.go:12: zero value store wrote into current directory`。零值 `FileProfileStore` 未返回错误，反而对当前目录执行了写入操作。
- `TestHealthCheckerRejectsEmptyTarget`：`health_test.go:11: empty target accepted`。HealthChecker 未拒绝空 target，直接接受了非法输入。
