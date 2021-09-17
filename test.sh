#! /bin/bash

cd ${WORKSPACE}/${projectPath}
bazel clean
# 单元测试命令
bazel coverage ... --test_arg=-test.v
# 收集覆盖率产物 cover.out
python ${WORKSPACE}/ops/ci/coverage/convert_bazel_dat_to_coverfile.py ${WORKSPACE}/${projectPath}/bazel-testlogs
zip coverage cover.out
cp coverage.zip ${testOutputDir}/coverage.zip
# 归置测试报告
zip report -r bazel-testlogs
cp report.zip ${testOutputDir}/report.zip
