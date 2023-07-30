# cli tool

## changing context
```bash
mg context use default
mg context use other
```

## testing

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


## generation

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