# [0.3.2](https://git.woa.com/goom/mocker/compare/v0.2.8...v0.3.2)(2021-07-18)
### Features
* 支持参数表达式arg.In、arg.Any
* 支持指定多个依次返回的值: Returns(...)
* 支持同时指定多个参数匹配条件: Matches(...)
* 修复若干问题

# [0.2.8](https://git.woa.com/goom/mocker/compare/v0.2.2...v0.2.8)(2021-03-19)
### Features
* 支持go1.16版本

# [0.2.2](https://git.woa.com/goom/mocker/compare/v0.1.6...v0.2.2)(2020-08-19)
### Features
* 支持接口Mock

### Bug Fixes
* 修复指定pkg后未重置问题
* 修复默认andreturn无效问题

# [0.1.6](https://git.woa.com/goom/mocker/compare/v0.1.5...v0.1.6)(2020-07-17)

### Bug Fixes
* 修复reset之后mock失败问题

# [0.1.5](https://git.woa.com/goom/mocker/compare/v0.1.4...v0.1.5)(2020-07-17)

### Bug Fixes
* method when参数个数不匹配问题
* 支持return nil

# [0.1.4](https://git.woa.com/goom/mocker/compare/v0.1.3...v0.1.4)(2020-07-17)

### Bug Fixes
* 修复无返回参数方法的mock (f247d63)
### Features
* 每一次mock支持设置Pkg (0e9407e)

# [0.1.3](https://git.woa.com/goom/mocker/compare/v0.1.2...v0.1.3)(2020-07-06)

### Features
* 支持windows系统
* 支持when条件匹配参数时Return

# [0.1.2](https://git.woa.com/goom/mocker/compare/v0.1.1...v0.1.2) (2020-06-28)


### Bug Fixes

* 修复export-as mock失败的问题 ([50b7c5a](https://git.woa.com/goom/mocker/commits/50b7c5a78e2c33597ebd13fb4a08481cac3d1dab))



# [0.1.1](https://git.woa.com/goom/mocker/compare/v0.1.0...v0.1.1) (2020-06-28)


### Features

* 支持条件匹配特性接口设计 ([2c725e3](https://git.woa.com/goom/mocker/commits/2c725e3df42aeb68c060e620d7a3f7d5a8c927e7))



# [0.1.0](https://git.woa.com/goom/mocker/compare/c869f80c895818959cc5a45ecf6f47466356fedd...v0.1.0) (2020-06-26)


### Features

* Mocker接口实现 ([77de276](https://git.woa.com/goom/mocker/commits/77de276e14aca2395952af654ab1b33949b2cff7))
* 增加测试方法 ([c869f80](https://git.woa.com/goom/mocker/commits/c869f80c895818959cc5a45ecf6f47466356fedd))
* 更新接口 ([2883e35](https://git.woa.com/goom/mocker/commits/2883e356f7f07e0a22f4fa25b748617113723754))
* 更新测试 ([98ae228](https://git.woa.com/goom/mocker/commits/98ae228f28bcb4e3932b32b8b9292751c1b75d0f))
