# mg: golang development at the speed of thought

## 简介

mg是一个轻量级的、简单的命令行工具，它能够改善我们使用golang开发时的体验。它包含了以下功能：

- 测试运行：由于鉴权，当必须使用命令行而非IDE运行测试时，只需复制测试用例名称，无需传入package路径，就可运行测试。
- 代码生成：
  - 简单依赖注入模式下的服务类代码生成
  - 简单依赖注入模式下的测试代码生成

详细的功能介绍，请看后文对应的章节。

## 安装

![GitHub Tag](https://img.shields.io/github/v/tag/crclz/mg?sort=semver)

```bash
go install github.com/crclz/mg@latest

# or:
# 注意字节内部的goproxy在代理github时不会很及时（天级延迟），所以需要临时采用其他的goproxy.
GOPROXY=https://goproxy.cn go install github.com/crclz/mg@latest

# or:
# 任何代理的 @latest 都会有一定时间的延迟，请将上方图片中的 tag 的后面的版本号作为最新版本，即可解决问题。
go install github.com/crclz/mg@v9.9.9 # 请替换 v9.9.9 为有效的版本号

```

## 运行测试

在本地开发的场景下，IDE可方便地运行绝大部分的测试；但是如果在生产网进行开发和测试，那么就会使用到鉴权命令行工具，例如`doas`或者 docker mesh 镜像。这些工具可能无法与IDE兼容，开发者运行测试的方式也变为命令行。

在命令行运行测试的最常用的方法，是这样：

```bash
# 自己敲命令: 需指定包名称
go test -v ./biz/utils --run ^TestAbcd$
```

通过 mg 可以这样运行：

```bash
# mg: 只需指定方法名称
mg t TestAbcd
```


## Context Management

mg-context是自定义的配置文件。

例如，如果想要在 `go test` 之前加一些命令前缀，例如用于鉴权的 `doas -p p.s.m go test XXXX`，可以在mg-context配置文件里面修改GoTestPrefix字段。


```bash
mg create-context default # 创建新context，会生成文件: mg-context.default.yaml

mg use-context --query # 获取目前使用的context名称

mg use-context default # 修改目前使用的context
mg use-context other # 修改目前使用的context
```

示例:

```
Go:
    GoTestPrefix: ["doasp", "-p", "p2.s2.m2"] # go test 前缀命令
    GoBuildNoOptim: true # 禁止编译优化，在测试中使用mockey时常常需要打开此开关
```

## 代码生成

### 服务生成（简单依赖注入风格的）

在了解代码生成前，先了解依赖注入：[dependency-injection](./docs/dependency-injection.md)

生成一个[简单依赖注入风格](./docs/dependency-injection.md#简单依赖注入)的service类，运行命令：

`mg g s biz/service/network_service`


### 测试生成（简单依赖注入风格的）

1. 新建或者打开已有的测试文档。输入：

```go
//!magic test SomeService.SomeMethod return false when command is nil 
```

其中，`//!magic test`是固定语法。`SomeService.SomeMethod`分别代表类名和方法名。`return false when command is nil`是behavior，会被转换为测试函数名称的一部分。

为了方便直接从代码进行复制，这种语法也是支持的：`SomeService) SomeMethod`

2. 命令行执行 `mg magic`，就会生成这种代码：

```go
func TestSomeService_SomeMethod_returnFalseWhenCommandIsNil(t *testing.T) {
	// arrange
	var assert = require.New(t)
	var ctx = context.Background()
	assert.NotNil(ctx)
	// var someService = mg.GetSingletonSomeService()

	// act
	// someService.SomeMethod()

	// assert
}
```

3. 可以在mg-context配置文件中修改arrange的部分代码(TestArrangePart)：
```yaml
Go:
  ...
  Magic:
    TestArrangePart: |-
      var assert = require.New(t)
      var ctx = context.Background()
      assert.NotNil(ctx)
```




## 测试（进阶用法）

### 其他参数

options:
- `--c1`: add --count=1 to argument
- `--dry-run`: only print command, not run test
- `--script $GoScriptName`: add GoScriptName as environment variable, see mgtesting/

### Running Tests with Mesh

有时下游会开启~~严格鉴权~~服务鉴定，此时开发调试时需要套一层mesh。

1. 将 [go_test_with_mesh.sample.sh](./internal/application/go_test_with_mesh.sample.sh) 复制到仓库，重命名为 `go_test_with_mesh.sh`
    - 修改sh脚本中的 ServicePsm
    - 根据仓库的需要，做其他定制化的修改
2. 修改 mg-context.*.yaml, 在 `Go` 字段下，添加 `MeshTestCommand: ["bash", "go_test_with_mesh.sh"]`
3. 运行测试时添加 `--mesh` 参数，例如: `mg t --mesh TestXXXX`

