# CLI Tool

## Installation

```bash
go install github.com/crclz/mg@latest

# or:

# 注意字节内部的goproxy在代理github时不会很及时（天级延迟），所以需要临时采用其他的goproxy.
GOPROXY=https://goproxy.cn go install github.com/crclz/mg@latest
```

## Context Management

```bash
mg create-context default # 创建新context，会生成文件: mg-context.default.yaml

mg use-context --query # 获取目前使用的context名称

mg use-context default # 修改目前使用的context
mg use-context other # 修改目前使用的context
```

## Testing

```bash
# automatically discover test: go test -v ./biz/service --run TestXXX_abcd
mg t TestSomeClass_SomeMethod

# run script
# TODO: script safe assert
# TODO: long running tests
mg t --script TestSomeScript123
```

options:
- `--c1`: add --count=1 to argument

context configs:
- GoTestPrefix: will prepend this prefix to go test command. e.g. `[doas, -p, p.s.m]`


## Generation

```bash
# generate a simple singleton service in biz/service/network_service.go
mg g s biz/service/network_service

# generate test
# TestSomeClass_SomeMethod_whenProvideNilInput_thenReturnError
mg g t SomeClass.SomeMethod when provide nil input then return error
mg g t [*]ExampleService[)] SomeMethod when provide nil input then return error
```
