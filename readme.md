# CLI Tool

## Installation

```bash
go install github.com/crclz/mg@latest
```

## Context Management

```bash
mg create-context default # 创建新context，会生成文件: mg-context.default.yaml
mg get-context # 获取目前使用的context名称
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

GoTestPrefix: 


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