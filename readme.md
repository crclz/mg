# CLI Tool

## Installation

```bash
go install github.com/crclz/mg@latest

# or:

# 注意字节内部的goproxy在代理github时不会很及时，所以需要临时采用其他的goproxy.
GOPROXY=https://goproxy.cn go install github.com/crclz/mg@latest
```

## Context Management

```bash
mg create-context default # 创建新context，会生成文件: mg-context.default.yaml

# TODO: mg create-context --preset byted # 创建context，附带doas -p XXX，并且禁止编译优化

mg use-context --query # 获取目前使用的context名称

mg use-context default # 修改目前使用的context
mg use-context other # 修改目前使用的context
```

## Testing

```bash
# automatically discover test: go test -v ./biz/service --run TestXXX_abcd
mg t TestSomeClass_SomeMethod

# run last test
mg t l

# run script
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

# generate a normal singleton service
mg g s --singleton --wire biz/dependency_building biz/service/NetworkService

# generate a scoped service
mg g s --scoped --wire biz/dependency_building biz/service/DataCacheService
```


## 
feature proposal

```bash
mg add-dep NetworkService to DataCacheService

```

config add:
- singleton wire, scoped wire path