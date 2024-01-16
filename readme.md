# CLI Tool

## Installation

```bash
go install github.com/crclz/mg@latest

# or:

# 注意字节内部的goproxy在代理github时不会很及时（天级延迟），所以需要临时采用其他的goproxy.
GOPROXY=https://goproxy.cn go install github.com/crclz/mg@latest
```

## Running Tests (Basic)
```bash
# 自己敲命令: 需指定包名称
go test -v ./biz/utils --run ^TestAbcd$

# mg: 只需指定方法名称
mg t TestAbcd

# TODO: mg t --group TestSomeServiceSomeMethod: running test in a single package
```

## Code Generation

```bash
# generate a simple singleton service in biz/service/network_service.go
mg g s biz/service/network_service

# TODO: generate test (This feature is not implemented yet!)
# TestSomeClass_SomeMethod_when_provide_nil_input_then_return_error
mg g t SomeClass.SomeMethod when provide nil input then return error
mg g t [*]SomeClass[)|.] SomeMethod when provide nil input then return error
```

## Magic

*Not Implemeted!*

`mg magic [--file a/b/c_test.go]`

```go
//! test SomeMethod when some condition return some value
//! test SomeClass.SomeMethod when some condition return some value
//! test SomeClass.SomeMethod:
// when cond1 return value1
// when cond2 return some error
// when cond3 return some shit
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

## Running Tests with Mesh

有时下游会开启~~严格鉴权~~服务鉴定，此时开发调试时需要套一层mesh。

1. 将 [go_test_with_mesh.sample.sh](./internal/application/go_test_with_mesh.sample.sh) 复制到仓库，重命名为 `go_test_with_mesh.sh`
    - 修改sh脚本中的 ServicePsm
    - 根据仓库的需要，做其他定制化的修改
2. 修改 mg-context.*.yaml, 在 `Go` 字段下，添加 `MeshTestCommand: ["bash", "go_test_with_mesh.sh"]`
3. 运行测试时添加 `--mesh` 参数，例如: `mg t --mesh TestXXXX`

## Running Tests (Advanced)

options:
- `--c1`: add --count=1 to argument
- `--dry-run`: only print command, not run test
- `--script $GoScriptName`: add GoScriptName as environment variable, see mgtesting/