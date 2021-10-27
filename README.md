# fontend
这是configmanagedemo的后端源码
如果你需要进行二次开发或是打包生产环境代码，可以使用以下命令

## Setup
```
go mod vendor
```

## Build
```
go build -o bin/configmanage-server cmd/web/main.go
```

## Config
将配置模板文件`config.example.toml`改名为`config.toml`，在`config.toml`里设置相关信息。
```
title = "configmanage"

[database]
type="mysql"
host="localhost"
port=3306
username="root"
password=""
name="configmanage"

[server]
appKey="OzMeDSvb6yzeCKHeayB64rpiWQGMNei2JwWXk6l5QlX3SUEJKQIA0qK2qHcqsraz"
ip="0.0.0.0"
port=8080
mode="init"
webhook="http://127.0.0.1:8080/webhook/test"
lruSize=50 # lru缓存量

[redis]
db=0
network="tcp"
address="127.0.0.1:6379"
password=""

```


## Serve
```
bin/configmanage-server config.toml
```

## Api
https://documenter.getpostman.com/view/2406114/UV5dcuF1