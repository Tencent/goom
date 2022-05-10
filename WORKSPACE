workspace(
    name = "mocker",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "tencent2",
    sha256 = "1fa5633f9e95014e5179186f65044291786328f7469831366bc374f5fbcbcaca",
    url = "http://mirrors.tencent.com/repository/generic/bazel/legacy/tencent2.tgz",
)

load("@tencent2//third_party:deps.bzl", "dependencies")
dependencies()

load("@tencent2//third_party:install.bzl", "install")

install()

load("//:deps.bzl", "third_party_go_dependencies")

# gazelle:repository_macro deps.bzl%third_party_go_dependencies
third_party_go_dependencies()
