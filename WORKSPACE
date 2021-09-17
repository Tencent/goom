workspace(
    name = "mocker",
)

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "tencent2",
    commit = "0c979cd9267063882a90aafcb48e4d1c19d27d1c",
    remote = "git@git.code.oa.com:depot/tencent2.git",
)

load("@tencent2//third_party:deps.bzl", "dependencies")
dependencies()

load("@tencent2//third_party:install.bzl", "install")

install()

load("//:deps.bzl", "third_party_go_dependencies")

# gazelle:repository_macro deps.bzl%third_party_go_dependencies
third_party_go_dependencies()
