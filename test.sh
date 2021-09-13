#! /bin/bash

# 未 Bazel 化的 Go 项目 test.sh 示例，仅为示例，请根据实际情况调整！
cd ${WORKSPACE}/${projectPath}
# 单元测试命令
go test -v -gcflags=all=-l -covermode=count -coverpkg=./... -coverprofile=coverage_unit.out $(go list ./... | grep -v apitest) | tee report.out
# 归置测试产物
# 归置单测覆盖率报告文件
zip coverage coverage_unit.out
cp coverage.zip ${testOutputDir}/coverage.zip
# 归置测试报告文件
cp report.out ${testOutputDir}/report.out